package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type (
	service struct {
		Group   string
		Name    string
		Version string
	}

	db struct {
		Type     string
		Host     string
		Port     int
		User     string
		Password string
		Database string
		Log      bool
	}

	mgmtEndpoint struct {
		Port            int
		BaseRoutingPath string
	}

	mgmtHealth struct {
		Path string
		Full bool
	}
	management struct {
		Endpoint mgmtEndpoint
		Health   mgmtHealth
	}

	log struct {
		Path     string
		Filename string
		Console  struct {
			Enabled       bool
			DisableColors bool
			Colors        bool
		}
		Level           string
		JSON            bool
		MaxSize         int
		MaxBackups      int
		MaxAge          int
		Compress        bool
		LocalTime       bool
		TimestampFormat string
		FullTimestamp   bool
		ForceFormatting bool
	}

	api struct {
		Endpoint struct {
			Port            int
			BaseRoutingPath string
		}
		Security struct {
			Enabled bool
			Jwt     struct {
				Secret     string
				Expiration struct {
					Enabled bool
					Minutes int32
				}
			}
		}
	}

	// Ldap configuration
	Ldap struct {
		Base   string
		Host   string
		Port   int
		UseSSL bool
		Bind   struct {
			DN       string
			Password string
		}
		UserFilter  string
		GroupFilter string
		Attributes  []string
	}

	// Configuration describe the type for the configuration file
	Configuration struct {
		Service    service
		API        api
		DB         db
		Management management
		Log        log
		Ldap       Ldap
	}
)

var instance *Configuration
var once sync.Once

// GetConfiguration returns the Configuration structure singleton instance.
func GetConfiguration() *Configuration {
	once.Do(func() {
		loadConfiguration()
	})

	return instance
}

func loadConfiguration() {
	viper.SetDefault("logPath", "./log")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if os.Getenv("ENV") != "" {
		viper.SetConfigName("config-" + os.Getenv("ENV"))
	} else {
		viper.SetConfigName("config")
	}

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	if err := viper.Unmarshal(&instance); err != nil {
		panic(fmt.Errorf("fatal error decoding configuration into struct: %v", err))
	}

}

// Get returns a configuration map by key. Used for custom or gear configurations.
func Get(key string) interface{} {
	// just in case!
	conf := GetConfiguration()
	if conf == nil {
		panic("No configuration at all!")
	}
	return viper.Get(key)
}

// IsSet checks to see if the key has been set in any of the data locations.
// IsSet is case-insensitive for a key.
func IsSet(key string) bool {
	// just in case!
	conf := GetConfiguration()
	if conf == nil {
		panic("No configuration at all!")
	}
	return viper.IsSet(key)
}
