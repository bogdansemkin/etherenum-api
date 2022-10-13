package main

import (
	"etherenum-api/etherenum-service/api/internal/app"
	"etherenum-api/etherenum-service/api/internal/config"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

func main () {
	var cfg config.Config
	err :=cleanenv.ReadEnv(&cfg)
	if err != nil {
		fmt.Errorf("failed to read env, %s", err)
	}

	app.Run(&cfg)
}
