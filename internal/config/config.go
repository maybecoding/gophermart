package config

import (
	"flag"
	"fmt"
	"os"
)

type (
	Config struct {
		HTTP
		PG
		AccrualSystem
		Log
		JWT
	}
	HTTP struct {
		Address string
	}
	PG struct {
		URI string
	}
	AccrualSystem struct {
		Address string
	}
	Log struct {
		Level string
	}
	JWT struct {
		Secret       string
		ExpiresHours int
	}
)

func New() (cfg *Config, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("config.New error: %w", err)
		}
	}()

	fHttpAddress := flag.String("a", "localhost:8080", "HTTP server address")
	httpAddress := os.Getenv("RUN_ADDRESS")

	fPgURI := flag.String("d", "postgres://api:pwd@localhost:5432/mart?sslmode=disable", "Postgres database URI")
	pgURI := os.Getenv("DATABASE_URI")

	fAccrualSystemAddress := flag.String("r", "localhost:8090", "Accrual system address")
	accrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")

	fLogLevel := flag.String("l", "debug", "Log level")
	logLevel := os.Getenv("LOG_LEVEL")

	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.Usage()
		return nil, fmt.Errorf("error due flag.Parse (extra arguments are passed) :%w", err)
	}

	if *fHttpAddress != "" {
		httpAddress = *fHttpAddress
	}
	if *fPgURI != "" {
		pgURI = *fPgURI
	}
	if *fAccrualSystemAddress != "" {
		accrualSystemAddress = *fAccrualSystemAddress
	}
	if *fLogLevel != "" {
		logLevel = *fLogLevel
	}
	return &Config{
		HTTP:          HTTP{httpAddress},
		PG:            PG{pgURI},
		AccrualSystem: AccrualSystem{accrualSystemAddress},
		Log:           Log{logLevel},
		JWT: JWT{
			Secret:       "super complex secret nobody can read it",
			ExpiresHours: 10,
		},
	}, nil
}
