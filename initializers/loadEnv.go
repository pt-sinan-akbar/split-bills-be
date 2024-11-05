package initializers

import "github.com/spf13/viper"

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
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
