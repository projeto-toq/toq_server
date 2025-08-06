package config

import (
	"fmt"
	"os"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"gopkg.in/yaml.v3"
)

func (c *config) LoadEnv() {
	env := globalmodel.Environment{}

	data, err := os.ReadFile("configs/env.yaml")
	if err != nil {
		fmt.Printf("error reading the env file: %v\n", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &env)
	if err != nil {
		fmt.Printf("error unmarshalling the env: %v\n", err)
		os.Exit(1)
	}

	c.env = env
	globalmodel.SetJWTSecret(env.JWT.Secret)
}
