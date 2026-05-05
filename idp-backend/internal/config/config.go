package config

import (
	"os"
)

type Config struct {
	DBURL         string
	KafkaBrokers  string
	KubeConfig    string
	BaseDomain    string
	ServerPort    string
	RegistryURL   string
}

func LoadConfig() *Config {
	return &Config{
		DBURL:        getEnv("DB_URL", "postgres://postgres:postgres@192.168.2.3:5432/idp?sslmode=disable"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "192.168.2.2:9092"),
		KubeConfig:   getEnv("KUBECONFIG", os.Getenv("USERPROFILE")+"/.kube/config"), // Default to local kubeconfig
		BaseDomain:   getEnv("BASE_DOMAIN", "idp.dev"),
		ServerPort:   getEnv("PORT", "8080"),
		RegistryURL:  getEnv("DOCKER_REGISTRY", "localhost:5000"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
