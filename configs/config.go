package config

// Config config value
type Config struct {
	Db Database
}

// Database type
type Database struct {
	Host string
	Port int
	User string
	Pwd  string
}
