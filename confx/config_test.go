package confx

import "testing"

func TestLoad(t *testing.T) {
	LoadConfig("./websocket.toml")
}
