package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	Token string
	Link  string
}

func GetConfig() (*Config, error) {
	file, err := os.Open("./config/app.env")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	env, err := godotenv.Parse(file)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:  env["PORT"],
		Token: env["BOT_TOKEN"],
		Link:  env["WEBHOOK_LINK"],
	}, nil
}
