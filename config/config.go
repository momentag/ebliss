package config

import "github.com/hashicorp/hcl/v2"

type BackendConfig struct {
	Name    string   `hcl:",label"`
	Address string   `hcl:"addr"`
	Token   string   `hcl:"token"`
	Options hcl.Body `hcl:",remain"`
}

type ServiceConfig struct {
	Name       string           `hcl:",label"`
	ListenAddr string           `hcl:"listen"`
	Backends   []*BackendConfig `hcl:"backend,block"`
}

type Config struct {
	Services []*ServiceConfig `hcl:"service,block"`
}
