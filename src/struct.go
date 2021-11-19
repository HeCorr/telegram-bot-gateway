package main

type Bot struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Path     string `yaml:"path"`
}

type Bots struct {
	Bots []Bot `yaml:"bots"`
}

// Prefixes Endpoint with a / if necessary
func (b *Bot) NormalizeEndpoint() {
	if b.Endpoint[0] != '/' {
		b.Endpoint = "/" + b.Endpoint
	}
}
