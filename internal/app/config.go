package app

type Config struct {
	Host        string
	Port        string
	Database    string
	User        string
	Password    string
	ForceRemove bool
}

func NewConfig() Config {
	return Config{
		Host:        "localhost",
		Port:        "5432",
		Database:    "postgres",
		User:        "postgres",
		Password:    "postgres",
		ForceRemove: false,
	}
}
