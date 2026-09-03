/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	aiv1 "github.com/example/ai-operator/api/v1"
)

// AIServiceReconciler reconciles a AIService object.
type AIServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ai.ai.example.com,resources=aiservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ai.ai.example.com,resources=aiservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ai.ai.example.com,resources=aiservices/finalizers,verbs=update

func (r *AIServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 1. 获取 AIService 实例
	var aiService aiv1.AIService
	if err := r.Get(ctx, req.NamespacedName, &aiService); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconciling AIService", "name", aiService.Name, "model", aiService.Spec.Model)

	// 2. 部署或更新关联的 Deployment
	deploy := &appsv1.Deployment{}
	deployName := req.NamespacedName

	err := r.Get(ctx, deployName, deploy)
	if errors.IsNotFound(err) {
		// 创建 Deployment
		desiredDeploy, err := r.buildDeployment(&aiService)
		if err != nil {
			logger.Error(err, "Failed to build Deployment")
			return ctrl.Result{}, err
		}
		logger.Info("Creating Deployment", "deployment", desiredDeploy.Name)
		if err := r.Create(ctx, desiredDeploy); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		// 存在则校验并保持同步（幂等更新副本数与镜像）
		desiredReplicas := int32(1)
		if aiService.Spec.Replicas != nil {
			desiredReplicas = *aiService.Spec.Replicas
		}

		updated := false
		if *deploy.Spec.Replicas != desiredReplicas {
			deploy.Spec.Replicas = &desiredReplicas
			updated = true
		}
		if deploy.Spec.Template.Spec.Containers[0].Image != aiService.Spec.Image {
			deploy.Spec.Template.Spec.Containers[0].Image = aiService.Spec.Image
			updated = true
		}

		if updated {
			logger.Info("Updating Deployment spec", "deployment", deploy.Name)
			if err := r.Update(ctx, deploy); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// 3. 部署或更新关联的 Service
	svc := &corev1.Service{}
	err = r.Get(ctx, deployName, svc)
	if errors.IsNotFound(err) {
		desiredSvc, err := r.buildService(&aiService)
		if err != nil {
			logger.Error(err, "Failed to build Service")
			return ctrl.Result{}, err
		}
		logger.Info("Creating Service", "service", desiredSvc.Name)
		if err := r.Create(ctx, desiredSvc); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// 4. 回写 Status
	readyReplicas := deploy.Status.ReadyReplicas
	phase := "Pending"
	message := "Waiting for pods to be ready"

	if readyReplicas > 0 && readyReplicas == *deploy.Spec.Replicas {
		phase = "Running"
		message = "All model server pods are ready"
	} else if readyReplicas > 0 {
		phase = "Progressing"
		message = fmt.Sprintf("%d/%d pods are ready", readyReplicas, *deploy.Spec.Replicas)
	}

	if aiService.Status.Phase != phase || aiService.Status.ReadyReplicas != readyReplicas {
		aiService.Status.Phase = phase
		aiService.Status.ReadyReplicas = readyReplicas
		aiService.Status.Message = message

		logger.Info("Updating AIService status", "phase", phase, "readyReplicas", readyReplicas)
		if err := r.Status().Update(ctx, &aiService); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *AIServiceReconciler) buildDeployment(aiService *aiv1.AIService) (*appsv1.Deployment, error) {
	replicas := int32(1)
	if aiService.Spec.Replicas != nil {
		replicas = *aiService.Spec.Replicas
	}

	labels := map[string]string{
		"app":        aiService.Name,
		"controller": "aiservice-operator",
	}

	// 将业务 Spec 中的 GPU 数量映射为 nvidia.com/gpu 资源限制
	resourceLimits := corev1.ResourceList{
		corev1.ResourceMemory: resource.MustParse(aiService.Spec.Resources.Memory),
	}
	if aiService.Spec.Resources.GPU > 0 {
		resourceLimits[corev1.ResourceName("nvidia.com/gpu")] = resource.MustParse(strconv.Itoa(int(aiService.Spec.Resources.GPU)))
	}

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      aiService.Name,
			Namespace: aiService.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "model-server",
							Image: aiService.Spec.Image,
							Args: []string{
								"--model", aiService.Spec.Model,
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8000,
									Name:          "http",
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits:   resourceLimits,
								Requests: resourceLimits,
							},
						},
					},
				},
			},
		},
	}

	// 绑定 OwnerReference，方便垃圾回收
	if err := controllerutil.SetControllerReference(aiService, deploy, r.Scheme); err != nil {
		return nil, err
	}

	return deploy, nil
}
func (r *AIServiceReconciler) buildService(aiService *aiv1.AIService) (*corev1.Service, error) {
	labels := map[string]string{
		"app":        aiService.Name,
		"controller": "aiservice-operator",
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      aiService.Name,
			Namespace: aiService.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       8000,
					TargetPort: intstr.FromInt(8000),
				},
			},
		},
	}

	// 绑定 OwnerReference
	if err := controllerutil.SetControllerReference(aiService, svc, r.Scheme); err != nil {
		return nil, err
	}

	return svc, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AIServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiv1.AIService{}).
		Named("aiservice").
		Complete(r)
}
