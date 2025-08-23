package config

type Config struct {
	Endpoint int
	Page     uint
}

func (c *Config) Next() {
	c.Page += 1
}

func (c *Config) Previous() {
	if c.Page == 0 {
		c.Page = 0
	} else {
		c.Page -= 1
	}
}
