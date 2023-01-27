package output

import (
	"context"

	"github.com/takutakahashi/github-token-renewer/pkg/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (k Kubernetes) Output(token string) error {
	ctx := context.Background()
	c := k.c.CoreV1().Secrets(k.cfg.SecretNamespace)
	secret, err := c.Get(ctx, k.cfg.SecretName, v1.GetOptions{})
	if err != nil {
		return err
	}
	secret.StringData[k.cfg.Key] = token
	if _, err := c.Update(ctx, secret.DeepCopy(), v1.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}
