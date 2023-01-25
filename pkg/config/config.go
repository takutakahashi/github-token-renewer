package config

type Config struct {
	URL            string
	PrivateKeyPath string
	AppID          int64
	Output         Output
}

type Output struct {
	KubernetesSecret *OutputKubernetesSecret
}

type OutputKubernetesSecret struct {
	SecretName string
	Key        string
}
