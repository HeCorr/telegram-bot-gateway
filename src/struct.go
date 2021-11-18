package main

type Bot struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Path     string `yaml:"path"`
}

type Bots struct {
	Bots []Bot `yaml:"bots"`
}
