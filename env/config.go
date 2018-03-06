package env

// SMTP smtp
type SMTP struct {
	User     string
	Password string
	Port     string
	Host     string
}

// Config configuration model
type Config struct {
	Port           int
	Theme          string
	Administrators []string
}
