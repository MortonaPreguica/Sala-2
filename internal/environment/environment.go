package environment

type Envs struct {
	ServerPort string `env:"SERVER_PORT"`
}

var Env Envs
