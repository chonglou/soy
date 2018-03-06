package env

const (
	// ROOT root dir
	ROOT = "tmp"
)

// SMTP smtp
type SMTP struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
}

// Config configuration model
type Config struct {
	Port           int               `toml:"port"`
	Theme          string            `toml:"theme"`
	Administrators []string          `toml:"administrators"`
	ReCaptcha      ReCaptcha         `toml:"recaptcha"`
	SMTP           SMTP              `toml:"smtp"`
	Site           map[string]string `toml:"site"`
	Header         []Link            `toml:"header"`
	Footer         []Link            `toml:"footer"`
}

// ReCaptcha https://www.google.com/recaptcha/intro/
type ReCaptcha struct {
	SiteKey   string `toml:"site-key"`
	SecretKey string `toml:"secret-key"`
}

// Link link
type Link struct {
	Title    string `toml:"title"`
	URL      string `toml:"url"`
	Children []Link `toml:"children"`
}
