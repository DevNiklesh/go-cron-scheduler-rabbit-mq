package config

import "os"

type Config struct {
	RabbitUrl   string
	Exchange    string
	QueueWorker string
	KeyWorker   string
}

func New() *Config {
	return &Config{
		RabbitUrl:   getEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672"),
		Exchange:    getEnv("EXCHANGE", "main_exchange"),
		QueueWorker: getEnv("QUEUE_WORKER", "worker_queue"),
		KeyWorker:   getEnv("KEY_WORKER", "worker_key"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
