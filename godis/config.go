package godis

type Config struct {
	ao bool
}

func DefaultConfig() *Config {

	return &Config{ao: false}
}

func (c *Config) SetAO(Value bool) {
	c.ao = Value
}

var Conf = DefaultConfig()
