package config

import "github.com/spf13/viper"

func ViperBind() {
	viper.BindEnv("ENV")

	// Binding Swagger Auth
	viper.BindEnv("SWAGGER_USERNAME")
	viper.BindEnv("SWAGGER_PASSWORD")

	// Binding Database
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("SSL_MODE")
	viper.BindEnv("TZ")

	// Binding GCP Cred
	viper.BindEnv("GOOGLE_APPLICATION_CREDENTIALS_BASE64")
}
