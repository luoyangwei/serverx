package config

import "testing"

func TestLoad(t *testing.T) {
	LoadConfig("./websocket.toml")
}
