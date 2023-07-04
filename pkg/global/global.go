package global

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-proxypool/pkg/storage"
)

var (
	Storage storage.Storage
	Logger  *log.Logger
	Config  *viper.Viper
)

func Initialize() {
	Logger = log.New()
	Logger.SetFormatter(&log.TextFormatter{})

	Config = viper.New()
	Config.SetEnvPrefix("PROXYPOOL")
	Config.SetConfigName("config")
	Config.SetConfigType("yaml")
	Config.AddConfigPath("/etc/proxypool")
	Config.AddConfigPath("$HOME/.proxypool")
	Config.AddConfigPath(".")

	Config.AutomaticEnv()
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := Config.ReadInConfig()
	if err != nil {
		Logger.Fatalf("failed to read config: %v", err)
	}

	// Set logger level.
	lvl := Config.GetString("log_level")
	level, err := log.ParseLevel(lvl)
	if err != nil {
		Logger.SetLevel(log.WarnLevel)
	}
	Logger.SetLevel(level)

	s := Config.GetString("storage")
	switch s {
	case "memory":
		Storage = storage.NewMemoryStorage()
		break
	case "redis":
		Storage = storage.NewRedisStorage(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", Config.GetString("redis.host"), Config.GetInt("redis.port")),
			DB:   Config.GetInt("redis.db"),
		})
		break
	default:
		Logger.Fatalf("unknown storage: %s", s)
	}
}
