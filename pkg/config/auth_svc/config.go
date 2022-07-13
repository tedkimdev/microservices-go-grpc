package auth_svc

import "github.com/spf13/viper"

type EnvConfig struct {
  Port         string `mapstructure:"PORT"`
  ReadOnlyDBURL        string `mapstructure:"READ_ONLY_DB_URL"`
  ReadWriteDBURL        string `mapstructure:"READ_WRITE_DB_URRL"`
  JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`
}

func LoadEnvConfig() (config EnvConfig, err error) {
  viper.AddConfigPath("./pkg/config/auth_svc")
  viper.SetConfigName("dev")
  viper.SetConfigType("env")

  viper.AutomaticEnv()

  err = viper.ReadInConfig()

  if err != nil {
    return
  }

  err = viper.Unmarshal(&config)

  return
}