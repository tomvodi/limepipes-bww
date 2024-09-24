package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/common"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/grpcplugin"
	"github.com/tomvodi/limepipes-plugin-bww/internal/bww"
	"github.com/tomvodi/limepipes-plugin-bww/internal/bww/parser"
	"github.com/tomvodi/limepipes-plugin-bww/internal/bww/symbolmapper"
	"github.com/tomvodi/limepipes-plugin-bww/internal/bwwfile"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common/helper"
	"github.com/tomvodi/limepipes-plugin-bww/internal/pluginimplementation"
	"google.golang.org/grpc"
)

// defaultGRPCServer returns a new gRPC server with the given options.
// Acts as a factory method for gRPC servers.
func defaultGRPCServer(opts []grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(opts...)
}

func main() {
	tok := bwwfile.NewTokenizer()
	tokConv := bwwfile.NewTokenConverter()
	sp := bwwfile.NewStructureParser(
		tok,
		tokConv,
	)
	symmap := symbolmapper.New()
	fsconv := bww.NewConverter(symmap)
	impl := pluginimplementation.NewPluginImplementation(
		parser.New(sp, fsconv),
		helper.NewTuneFixer(),
	)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: common.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			fileformat.Format_BWW.String(): grpcplugin.NewGrpcPlugin(
				impl,
			),
		},

		GRPCServer: defaultGRPCServer,
	})
}
