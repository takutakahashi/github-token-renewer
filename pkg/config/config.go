package config

type Config struct {
	URL            string
	PrivateKeyPath string
	Output         Output
}

type Output struct {
	KubernetesSecret *OutputKubernetesSecret
}

type OutputKubernetesSecret struct {
	SecretName string
	Key        string
}
