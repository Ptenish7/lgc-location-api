package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Build information -ldflags .
const (
	version    string = "dev"
	commitHash string = "-"
)

var cfg *Config

// GetConfigInstance returns service config
func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

// Database - contains all parameters database connection.
type Database struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Migrations string `yaml:"migrations"`
	Name       string `yaml:"name"`
	SslMode    string `yaml:"sslmode"`
	Driver     string `yaml:"driver"`
	MaxRetry   uint64 `yaml:"maxRetry"`
}

// Project - contains all parameters project information.
type Project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string
	CommitHash  string
}

// Retranslator - contains retranslator config
type Retranslator struct {
	ChannelSize     uint64        `yaml:"channelSize"`
	ConsumerCount   uint64        `yaml:"consumerCount"`
	ConsumerSize    uint64        `yaml:"consumerSize"`
	ConsumerTimeout time.Duration `yaml:"consumerTimeout"`
	ProducerCount   uint64        `yaml:"producerCount"`
	WorkerCount     uint64        `yaml:"workerCount"`
	BatchSize       uint64        `yaml:"batchSize"`
}

// Kafka - contains all parameters kafka information.
type Kafka struct {
	Capacity uint64   `yaml:"capacity"`
	Topic    string   `yaml:"topic"`
	GroupID  string   `yaml:"groupId"`
	Brokers  []string `yaml:"brokers"`
	MaxRetry uint64   `yaml:"maxRetry"`
}

// Telemetry config for service.
type Telemetry struct {
	GraylogPath string `yaml:"graylogPath"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project      Project      `yaml:"project"`
	Retranslator Retranslator `yaml:"retranslator"`
	Database     Database     `yaml:"database"`
	Kafka        Kafka        `yaml:"kafka"`
	Telemetry    Telemetry    `yaml:"telemetry"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return nil
}
