package rpc

import (
    "fmt"
    "github.com/go-ozzo/ozzo-config"
    "testing"
)

func TestConfig(t *testing.T) {
    c := config.New()
    c.Load("./config/main.json")
    Config(c)

    fmt.Println(rpc_config_group, rpc_client)
}
