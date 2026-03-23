package driver

type Config struct {
	BaseURL  string
	RootPath string
	URLPath  string
}

type Option func(*Config)
