package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func newKubernetesClient() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, err
		}
	}
	return kubernetes.NewForConfig(config)
}

type ListPodsInput struct {
	Namespace string `json:"namespace"`
}

type GetPodInput struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func registerK8sTools(ctx context.Context, s *mcp.Server) error {
	client, err := newKubernetesClient()
	if err != nil {
		return err // 集群不可用时整组降级，不阻塞其他 Tool
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_pods",
		Description: "List pods in a Kubernetes namespace with their phase.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in ListPodsInput) (*mcp.CallToolResult, any, error) {
		pods, err := client.CoreV1().Pods(in.Namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("list pods: %w", err)
		}
		out := make([]string, 0, len(pods.Items))
		for _, p := range pods.Items {
			out = append(out, fmt.Sprintf("%s [%s]", p.Name, p.Status.Phase))
		}
		return nil, out, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pod",
		Description: "Get a single pod by name and namespace.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in GetPodInput) (*mcp.CallToolResult, any, error) {
		pod, err := client.CoreV1().Pods(in.Namespace).Get(ctx, in.Name, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("get pod: %w", err)
		}
		return nil, pod, nil
	})
	return nil
}
