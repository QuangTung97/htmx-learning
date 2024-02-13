package config

import (
	"fmt"
	"strings"

	"github.com/QuangTung97/svloc"
	"github.com/spf13/viper"
)

type Config struct {
	Env   string      `mapstructure:"env"`
	Auth  Auth        `mapstructure:"auth"`
	MySQL MySQLConfig `mapstructure:"mysql"`
}

type Auth struct {
	CSRFHMACSecret string `mapstructure:"csrf_hmac_secret"`

	GoogleClientID     string `mapstructure:"google_client_id"`
	GoogleClientSecret string `mapstructure:"google_client_secret"`
}

func Load() Config {
	vip := viper.New()

	vip.SetConfigName("config")
	vip.SetConfigType("yml")
	vip.AddConfigPath(".")
	return loadConfig(vip)
}

var IsProdLoc = svloc.Register[bool](func(unv *svloc.Universe) bool {
	return Loc.Get(unv).Env == "production"
})

var Loc = svloc.RegisterEmpty[Config]()

func loadConfig(vip *viper.Viper) Config {
	vip.SetEnvPrefix("docker")
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vip.AutomaticEnv()

	err := vip.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// workaround https://github.com/spf13/viper/issues/188#issuecomment-399518663
	// to allow read from environment variables when Unmarshal
	for _, key := range vip.AllKeys() {
		val := vip.Get(key)
		vip.Set(key, val)
	}

	fmt.Println("Config file used:", vip.ConfigFileUsed())

	var cfg Config
	err = vip.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
