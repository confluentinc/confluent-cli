package connect

import (
	"context"

	"github.com/confluentinc/cli/shared"
	proto "github.com/confluentinc/cli/shared/connect"
	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type Plugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl Connect
}

type Connect interface {
	List(ctx context.Context) ([]*proto.Connector, error)
}

func (p *Plugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewConnectClient(c)}, nil
}

// Check that Plugin satisfies GPRCPlugin interface.
var _ plugin.GRPCPlugin = &Plugin{}

func (p *Plugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterConnectServer(s, &GRPCServer{p.Impl})
	return nil
}

func init() {
	shared.PluginMap["connect"] = &Plugin{}
}
