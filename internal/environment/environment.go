package environment

type Envs struct {
	ServerPort string `env:"SERVER_PORT"`
	Region     string `env:"REGION"`
}

var Env Envs
