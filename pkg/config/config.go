package config

type Config struct {
	URL            string
	PrivateKeyPath string
	AppID          int64
	Installations  []Installation
}

type Installation struct {
	ID     int64
	Output Output
}

type Output struct {
	KubernetesSecret *OutputKubernetesSecret
}

type OutputKubernetesSecret struct {
	SecretName      string
	SecretNamespace string
	Key             string
}
