package main

import (
	"context"
	"os"

	golog "log"

	"github.com/confluentinc/cli/command/connect"
	log "github.com/confluentinc/cli/log"
	"github.com/confluentinc/cli/shared"
	proto "github.com/confluentinc/cli/shared/connect"
	plugin "github.com/hashicorp/go-plugin"
)

func main() {
	var impl *Connect
	{
		logger := log.New()
		logger.Log("msg", "hello")
		defer logger.Log("msg", "goodbye")

		f, err := os.Create("/tmp/confluent-connect-plugin.log")
		check(err)
		logger.Logger.Out = f

		impl = &Connect{Logger: logger}
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"connect": &connect.Plugin{Impl: impl},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

type Connect struct {
	Logger *log.Logger
}

func (c *Connect) List(ctx context.Context) ([]*proto.Connector, error) {
	c.Logger.Log("msg", "connect.List()")
	return []*proto.Connector{
		{
			Id: "connector-1",
		},
	}, nil
}

func check(err error) {
	if err != nil {
		golog.Fatal(err)
	}
}
