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
	if env.AUTH.RefreshTTLDays > 0 {
		globalmodel.SetRefreshTTL(env.AUTH.RefreshTTLDays)
	}
	if env.AUTH.AccessTTLMinutes > 0 {
		globalmodel.SetAccessTTL(env.AUTH.AccessTTLMinutes)
	}
	if env.AUTH.MaxSessionRotations > 0 {
		globalmodel.SetMaxSessionRotations(env.AUTH.MaxSessionRotations)
	}
}
