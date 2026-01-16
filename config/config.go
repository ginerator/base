package config

type DbConfig struct {
	Name         string `env:"RDS_DBNAME"`
	Host         string `env:"RDS_HOST"`
	Port         string `env:"RDS_PORT"`
	Username     string `env:"RDS_USERNAME"`
	Password     string `env:"RDS_PASSWORD"`
	MaxOpenConns string `env:"MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns string `env:"MAX_IDLE_CONNS" default:"5"`
}

type AppConfig struct {
	Name string `env:"APP_NAME" default:"ta.item-service"`
	Port string `env:"APP_PORT" default:"3000"`
	Env  string `env:"APP_ENV" default:"development"`
	Host string `env:"APP_HOST" default:"localhost"`
}
type AuthConfig struct {
	Auth0Url string `env:"AUTH0_URL"`
}
