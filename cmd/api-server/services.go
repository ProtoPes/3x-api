package main

import awgconfig "3x-api/cmd/awg-config"
import "3x-api/internal/config"

// Send a signal to setup configuration for vpn services
func Bootstrap(c *config.Config) error {
	return awgconfig.Boot(c)
}
