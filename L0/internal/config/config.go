package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Config struct {
	DBPassword  string
	HTTPPort    string
	KafkaBroker string
	HostName    string
	IsKafka     bool
}

func LoadConfig() *Config {
	env := loadEnv()
	isDocker := checkDockerEnvironment()
	isKafka := false

	hostName := env["POSTGRES_HOST"]
	if hostName == "" {
		if isDocker {
			hostName = "postgres"
			isKafka = true
		} else {
			hostName = "localhost"
		}
	}

	return &Config{
		DBPassword:  env["DB_PASSWORD"],
		HTTPPort:    env["HTTP_PORT"],
		KafkaBroker: env["KAFKA_BROKERS"],
		HostName:    hostName,
		IsKafka:     isKafka,
	}
}

func loadEnv() map[string]string {
	env := make(map[string]string)

	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			part := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(part[0])
			value := strings.TrimSpace(part[1])
			env[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return env

}

func checkDockerEnvironment() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	content, err := os.ReadFile("/proc/1/cgroup")
	if err == nil && strings.Contains(string(content), "docker") {
		return true
	}

	return false
}
