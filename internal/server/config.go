package server

// Config ...
type Config struct {
	AppName  string `env:"APP" envDefault:"core"`
	Port     int    `env:"PORT" envDefault:"8080"`
	DDHost   string `env:"DD_HOST"`
	DDPort   string `env:"DD_PORT"`
	HostName string `env:"HOST_NAME"`
	AuthUrl  string `env:"AUTH_URL" envDefault:"http://ms-user-management/api/auth/verify-token"`
}
