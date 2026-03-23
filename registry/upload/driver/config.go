package driver

type Config struct {
	BaseURL  string
	RootPath string
	URLPath  string
	StorerID string
}

type Option func(*Config)
