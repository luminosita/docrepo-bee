package server

import (
	"context"
	http2 "net/http"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	http "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/luminosita/docrepo-bee/api/documents/v1"
	"github.com/luminosita/docrepo-bee/internal/conf"
	"github.com/luminosita/docrepo-bee/internal/service"
	"google.golang.org/grpc"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.Documents) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	mux := runtime.NewServeMux()
	// setting up a dail up for gRPC service by specifying endpoint/target url
	_ = v1.RegisterDocumentsHandlerFromEndpoint(context.Background(), mux,
		"localhost:9000", []grpc.DialOption{grpc.WithInsecure()})
	// Creating a normal HTTP server
	server := &http2.Server{
		Handler: mux,
	}

	srv.Server = server

	return srv
}
