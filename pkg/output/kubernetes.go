package output

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/takutakahashi/github-token-renewer/pkg/config"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Kubernetes struct {
	cfg config.OutputKubernetesSecret
	c   *kubernetes.Clientset
}

func NewKubernetes(output config.OutputKubernetesSecret) (*Kubernetes, error) {
	clientset, err := kubernetes.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		return nil, err
	}

	return &Kubernetes{cfg: output, c: clientset}, nil
}

func (k Kubernetes) Output(ctx context.Context, token string) error {
	c := k.c.CoreV1().Secrets(k.cfg.SecretNamespace)
	secret, err := c.Get(ctx, k.cfg.SecretName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		secret = &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      k.cfg.SecretName,
				Namespace: k.cfg.SecretNamespace,
			},
			StringData: map[string]string{},
		}
		secret.StringData[k.cfg.Key] = token
		if _, err := c.Create(ctx, secret, metav1.CreateOptions{}); err != nil {
			return err
		}
		logrus.Infof("new kubernetes secret created. name: %s/%s, key: %s", secret.Namespace, secret.Name, k.cfg.Key)
		return nil

	} else if err != nil {
		return err
	}
	secret.StringData = map[string]string{}
	secret.StringData[k.cfg.Key] = token
	if _, err := c.Update(ctx, secret.DeepCopy(), metav1.UpdateOptions{}); err != nil {
		return err
	}
	logrus.Infof("kubernetes secret updated. name: %s/%s, key: %s", secret.Namespace, secret.Name, k.cfg.Key)
	return nil
}
