package main

import (
	"etherenum-api/etherenum-service/api/internal/app"
	"etherenum-api/etherenum-service/api/internal/config"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	var cfg config.Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		fmt.Printf("failed to read env, %s", err)
	}

	err = app.Run(&cfg)
	if err != nil {
		fmt.Printf("error during start app, %s", err)
	}
}
