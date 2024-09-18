package config

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
)

type Config struct {
	StoragePath string
	Address     string
	ExternalAuth
}

type ExternalAuth struct {
	GoogleRedirectURL  string
	GoogleClientID     string
	GoogleClientSecret string
	GithubRedirectURL  string
	GithubClientID     string
	GithubClientSecret string
}

func Loader() *Config {
	addr := flag.String("addr", ":8081", "HTTP network address")
	dsn := flag.String("dsn", "./storage/storage.db", "Sql database storage")
	flag.Parse()

	conf := Config{
		StoragePath: *dsn,
		Address:     *addr,
		ExternalAuth: ExternalAuth{
			GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),

			GithubRedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
			GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		},
	}

	return &conf
}

func init() {
	file, err := os.Open(".env")
	if err != nil {
		log.Fatalf("open env file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				log.Print("ignoring .env line (invalid format):", line)
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if err := os.Setenv(key, value); err != nil {
				log.Fatalf("set env variable: %v", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("env parse scanner error: %v", err)
	}
}
