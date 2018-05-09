package main

import (
	"context"
	golog "log"
	"os"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"

	schedv1 "github.com/confluentinc/cc-structs/kafka/scheduler/v1"
	"github.com/confluentinc/cli/command/connect"
	chttp "github.com/confluentinc/cli/http"
	log "github.com/confluentinc/cli/log"
	metric "github.com/confluentinc/cli/metric"
	"github.com/confluentinc/cli/shared"
)

func main() {
	var logger *log.Logger
	{
		logger = log.New()
		logger.Log("msg", "hello")
		defer logger.Log("msg", "goodbye")

		f, err := os.OpenFile("/tmp/confluent-connect-plugin.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		check(err)
		logger.SetLevel(logrus.DebugLevel)
		logger.Logger.Out = f
	}

	var metricSink shared.MetricSink
	{
		metricSink = metric.NewSink()
	}

	var config *shared.Config
	{
		config = &shared.Config{
			MetricSink: metricSink,
			Logger:     logger,
		}
		err := config.Load()
		if err != nil && err != shared.ErrNoConfig {
			logger.WithError(err).Errorf("unable to load config")
		}
	}

	var impl *Connect
	{
		client := chttp.NewClientWithJWT(context.Background(), config.AuthToken, config.AuthURL, config.Logger)
		impl = &Connect{Logger: logger, Config: config, Client: client}
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
	Config *shared.Config
	Client *chttp.Client
}

func (c *Connect) List(ctx context.Context) ([]*schedv1.ConnectCluster, error) {
	c.Logger.Log("msg", "connect.List()")
	if c.Config.Auth == nil {
		return nil, shared.ErrUnauthorized
	}
	ret, _, err := c.Client.Connect.List(c.Config.Auth.Account.Id)
	return ret, shared.ConvertAPIError(err)
}

func (c *Connect) Describe(ctx context.Context, cluster *schedv1.ConnectCluster) (*schedv1.ConnectCluster, error) {
	c.Logger.Log("msg", "connect.Describe()")
	if c.Config.Auth == nil {
		return nil, shared.ErrUnauthorized
	}
	cluster.AccountId = c.Config.Auth.Account.Id
	ret, _, err := c.Client.Connect.Describe(cluster)
	return ret, shared.ConvertAPIError(err)
}

func check(err error) {
	if err != nil {
		golog.Fatal(err)
	}
}