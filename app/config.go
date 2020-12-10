package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/conf"
)

type AppConfig struct {
	Web struct {
		ServiceServerTLS bool          `conf:"default:false"`
		Address          string        `conf:"default:0.0.0.0:3000"`
		Debug            string        `conf:"default:0.0.0.0:4000"`
		ReadTimeout      time.Duration `conf:"default:5s"`
		WriteTimeout     time.Duration `conf:"default:5s"`
		ShutdownTimeout  time.Duration `conf:"default:5s"`
	}
	DB struct {
		User       string `conf:"default:postgres"`
		Password   string `conf:"default:postgres,noprint"`
		Host       string `conf:"default:db"`
		Name       string `conf:"default:postgres"`
		DisableTLS bool   `conf:"default:true"`
	}
	Trace struct {
		URL         string  `conf:"default:http://zipkin:9411/api/v2/spans"`
		Service     string  `conf:"default:illuminatingdeposits-rest"`
		Probability float64 `conf:"default:1"`
	}
}

func ParsedConfig(cfg AppConfig) (AppConfig, error) {
	if err := conf.Parse(os.Args[1:], "DEPOSITS", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("DEPOSITS", &cfg)
			if err != nil {
				return AppConfig{}, errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return AppConfig{}, nil
		}
		return AppConfig{}, errors.Wrap(err, "parsing config")
	}
	log.Printf("AppConfig is %v", cfg)
	return cfg, nil
}
