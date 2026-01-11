package initializers

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"reflect"
)

type Config struct {
	DBHost             string `mapstructure:"POSTGRES_HOST"`
	DBUserName         string `mapstructure:"POSTGRES_USER"`
	DBUserPassword     string `mapstructure:"POSTGRES_PASSWORD"`
	DBName             string `mapstructure:"POSTGRES_DB"`
	DBPort             string `mapstructure:"POSTGRES_PORT"`
	ServerPort         string `mapstructure:"PORT"`
	SupabaseStorageURL string `mapstructure:"SUPABASE_STORAGE_URL"`
	SupabaseSecretKey  string `mapstructure:"SUPABASE_SECRET_KEY"`
	SupabaseBucket     string `mapstructure:"SUPABASE_BUCKET"`
	SupbaseFolder      string `mapstructure:"SUPABASE_FOLDER"`
	GinMode            string `mapstructure:"GIN_MODE"`
	GeminiApiKey       string `mapstructure:"GEMINI_API_KEY"`
}

var ConfigSetting = &Config{}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()
	viper.SetDefault("GIN_MODE", "debug")

	if err = viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("? No env file found. Using environment variables.")
			config = parseConfigFromEnv(config)
		}
	}

	err = viper.Unmarshal(&config)
	ConfigSetting = &config
	return
}

func parseConfigFromEnv(config Config) Config {
	r := reflect.TypeOf(config)
	for r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	for i := 0; i < r.NumField(); i++ {
		env := r.Field(i).Tag.Get("mapstructure")
		if err := viper.BindEnv(env); err != nil {
			log.Fatal("? Failed to bind env variable:", err)
		}
	}
	return config
}
