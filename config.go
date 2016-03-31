package main

import "github.com/go-ini/ini"

type ServerConfig struct {
	BindPort int    `ini:"bind_port"`
	LocalDC  string `ini:"local_dc"`
}

func ParseConfigFile(filename string) (*ServerConfig, error) {
	/*
	   Read in the configuration values.
	*/
	cfg, err := ini.Load(filename)

	if err != nil {
		return nil, err
	}

	config := &ServerConfig{
		BindPort: 12345,
	}
	err = cfg.MapTo(config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
