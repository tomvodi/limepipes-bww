package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/common"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/grpc_plugin"
	"github.com/tomvodi/limepipes-plugin-bww/internal/bww"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common/helper"
	"github.com/tomvodi/limepipes-plugin-bww/internal/plugin_implementation"
	"google.golang.org/grpc"
)

// defaultGRPCServer returns a new gRPC server with the given options.
// Acts as a factory method for gRPC servers.
func defaultGRPCServer(opts []grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(opts...)
}

func main() {
	impl := plugin_implementation.NewPluginImplementation(
		bww.NewBwwParser(),
		helper.NewTuneFixer(),
		bww.NewBwwFileTuneSplitter(),
	)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: common.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"bww": grpc_plugin.NewGrpcPlugin(
				impl,
			),
		},

		GRPCServer: defaultGRPCServer,
	})
}
