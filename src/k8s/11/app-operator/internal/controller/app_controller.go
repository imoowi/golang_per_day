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

	appsappsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "github.com/example/app-operator/api/v1"
)

// AppReconciler reconciles a App object
type AppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.example.com,resources=apps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.example.com,resources=apps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.example.com,resources=apps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the App object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.24.1/pkg/reconcile
func (r *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// 1. 获取 App 资源实例
	var app appsv1.App
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. 根据 App.Spec 构建目标 Deployment
	desiredDeployment, err := r.deploymentForApp(&app)
	if err != nil {
		return ctrl.Result{}, err
	}

	// 3. 检查底层 Deployment 是否已存在
	var existingDeployment appsappsv1.Deployment
	err = r.Get(ctx, types.NamespacedName{Name: app.Name, Namespace: app.Namespace}, &existingDeployment)

	if errors.IsNotFound(err) {
		// 情况 A：不存在 ──► 创建新 Deployment
		if err := r.Create(ctx, desiredDeployment); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		// 情况 B：存在 ──► 检查镜像或副本数是否有偏差，有则更新
		if existingDeployment.Spec.Template.Spec.Containers[0].Image != app.Spec.Image ||
			*existingDeployment.Spec.Replicas != *app.Spec.Replicas {
			existingDeployment.Spec.Template.Spec.Containers[0].Image = app.Spec.Image
			existingDeployment.Spec.Replicas = app.Spec.Replicas
			if err := r.Update(ctx, &existingDeployment); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// deploymentForApp creates a Deployment object for the App
func (r *AppReconciler) deploymentForApp(app *appsv1.App) (*appsappsv1.Deployment, error) {
	replicas := app.Spec.Replicas
	deployment := &appsappsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			Labels: map[string]string{
				"app": app.Name,
			},
		},
		Spec: appsappsv1.DeploymentSpec{
			Replicas: replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": app.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": app.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  app.Name,
							Image: app.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}
	// 将 App 设置为该 Deployment 的 Owner（宿主）
	if err := ctrl.SetControllerReference(app, deployment, r.Scheme); err != nil {
		return nil, err
	}
	return deployment, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.App{}).             // 监听 App 的增删改
		Owns(&appsappsv1.Deployment{}). // 监听属于 App 的 Deployment 的增删改
		Named("app").
		Complete(r)
}
