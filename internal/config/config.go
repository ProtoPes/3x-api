package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

// TODO: Choose between yaml or env

type WgConfig struct {
	Port       string `env:"AWG_SERVER_PORT,required"`
	SubNetIP   string `env:"AWG_SUBNET_IP" envDefault:"10.8.1.0"`
	SubNetMask string `env:"AWG_SUBNET_MASK" envDefault:"24"`
	DNS        string `env:"AWG_DNS" envDefault:"9.9.9.9 149.112.112.112"`
	JcMin      int    `env:"JUNK_PACKET_COUNT_MIN" envDefault:"3"`
	JcMax      int    `env:"JUNK_PACKET_COUNT_MAX" envDefault:"10"`
	JPMin      int    `env:"JUNK_PACKET_MIN_SIZE" envDefault:"50"`
	JPMax      int    `env:"JUNK_PACKET_MAX_SIZE" envDefault:"1000"`
}

type Config struct {
	Env     string `env:"ENV" envDefault:"local"`
	HostURL string `env:"HOST_URL,required"`
	WgConfig
	// HTTPServer
}

// type HTTPServer struct {
// 	Address     string
// 	StoragePath string
// 	Timeout     time.Duration
// 	IdleTimeout time.Duration
// }

func MustLoad() *Config {
	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Can not load config: %s", err)
	}
	return cfg
}
