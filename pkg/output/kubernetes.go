package output

import (
	"github.com/takutakahashi/github-token-renewer/pkg/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Kubernetes struct {
	cfg config.OutputKubernetesSecret
	c   *kubernetes.Clientset
}

func NewKubernetes(output config.OutputKubernetesSecret) (*Kubernetes, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Kubernetes{cfg: output, c: clientset}, nil
}
