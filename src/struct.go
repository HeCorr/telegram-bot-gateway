package main

type Bot struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Path     string `yaml:"path"`
}

type Options struct {
	Key    string `yaml:"key"`
	Cert   string `yaml:"cert"`
	Listen string `yaml:"listen"`
	Strict bool   `yaml:"strict"`
}

type Config struct {
	Bots    []Bot   `yaml:"bots"`
	Options Options `yaml:"config"`
}

// Prefixes Endpoint with a / if necessary
func (b *Bot) NormalizeEndpoint() {
	if b.Endpoint[0] != '/' {
		b.Endpoint = "/" + b.Endpoint
	}
}
