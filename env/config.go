package env

import (
	"crypto/rand"
	"path"
)

const (
	// ROOT root dir
	ROOT = "tmp"
	// TITLE title key
	TITLE = "title"
)

// RandomBytes random bytes
func RandomBytes(l int) ([]byte, error) {
	buf := make([]byte, l)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// Config config filename
func Config() string {
	return path.Join(ROOT, "config.toml")
}

// SMTP smtp
type SMTP struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
}

// Env configuration model
type Env struct {
	Port           int               `toml:"port"`
	Theme          string            `toml:"theme"`
	Secrets        string            `toml:"secrets"`
	Administrators []string          `toml:"administrators"`
	SMTP           SMTP              `toml:"smtp"`
	Google         Google            `toml:"google"`
	Site           map[string]string `toml:"site"`
	Header         []Link            `toml:"header"`
	Footer         []Link            `toml:"footer"`
}

// Google google
type Google struct {
	VerifyID  string    `toml:"verify-id"`
	ReCaptcha ReCaptcha `toml:"recaptcha"`
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
