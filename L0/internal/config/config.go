package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Config struct {
	DBPassword string
	HTTPPort   string
}

func LoadConfig() *Config {
	env := loadEnv()

	return &Config{
		DBPassword: env["DB_PASSWORD"],
		HTTPPort:   env["HTTP_PORT"],
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
