package main

import (
	"etherenum-api/etherenum-service/api/internal/app"
	"etherenum-api/etherenum-service/api/internal/config"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

func main() {
	var cfg config.Config
	port := os.Getenv("PORT")

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		fmt.Printf("failed to read env, %s", err)
	}

	err = app.Run(&cfg, port)
	if err != nil {
		fmt.Printf("error during start app, %s", err)
	}
}
