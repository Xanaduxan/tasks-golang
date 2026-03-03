package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func LoadEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("cannot open .env file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		err := os.Setenv(key, value)
		if err != nil {
			return
		}
	}
}
