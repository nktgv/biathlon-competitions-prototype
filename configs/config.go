package configs

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Laps          int    `json:"laps"`
	LapLength     int    `json:"lapLen"`
	PenaltyLength int    `json:"penaltyLen"`
	FiringLines   int    `json:"firingLines"`
	Start         string `json:"start"`
	StartDelta    string `json:"startDelta"`
}

func LoadConfig(configPath string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
