package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDBURI  string
	JWTSecret   string
	OpenAIKey   string
	EmailAPIKey string
	OpenAIUrl   string
	EmailIdFrom string
	EmailIdTo   string
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")    // path to look for the config file in
	viper.SetConfigName("PROD") // name of config file (without extension)
	viper.SetConfigType("env")  // REQUIRED if the config file does not have the extension in the name
	// viper.AutomaticEnv()       // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	return &Config{
		MongoDBURI:  viper.GetString("MONGO_URI"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		OpenAIKey:   viper.GetString("OPEN_AI_KEY"),
		EmailAPIKey: viper.GetString("EMAIL_API_KEY"),
		OpenAIUrl:   viper.GetString("OPEN_AI_URL"),
		EmailIdFrom: viper.GetString("EMAIL_ID_FROM"),
		EmailIdTo:   viper.GetString("EMAIL_ID_TO"),
	}, nil
}
