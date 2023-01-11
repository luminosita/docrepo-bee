package bee

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	proto "github.com/luminosita/docrepo-bee/api/gen/v1"
	"github.com/luminosita/honeycomb/pkg/http"
	"github.com/luminosita/honeycomb/pkg/server"
	"google.golang.org/grpc"
)

type Config struct {
	Sc server.Config `mapstructure:"server"`
}

func (c *Config) ServerConfig() *server.Config {
	return &c.Sc
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
type BeeServer struct {
	c *Config
}

func NewBeeServer(c *Config) *BeeServer {
	c.Sc = server.Config{}

	return &BeeServer{
		c: c,
	}
}

func (bs *BeeServer) Config() server.ServerConfigurer {
	return bs.c
}

func (bs *BeeServer) OverrideConfigItems() map[string]string {
	return map[string]string{"config.server.baseUrl": "BaseUrl"}
}

func (*BeeServer) Routes(ctx context.Context) []*http.Route {
	routes := make([]*http.Route, 0)

	routes = append(routes, &http.Route{Type: http.STATIC, Path: "/assets"})

	routes = append(routes, &http.Route{
		Type: http.GET, Path: "/documents/:id", Handler: MakeGetDocumentHandler(ctx)})
	routes = append(routes, &http.Route{
		Type: http.GET, Path: "/documents/:id/info", Handler: MakeGetDocumentInfoHandler(ctx)})
	routes = append(routes, &http.Route{
		Type: http.POST, Path: "/documents", Handler: MakePutDocumentHandler(ctx)})

	return routes
}

func (*BeeServer) GrpcRegFunc(server *grpc.Server) {
	proto.RegisterDocumentsServer(server, MakeGrpcDocumentServer(context.Background()))
}

func (*BeeServer) GwRegFunc(ctx context.Context, mux *runtime.ServeMux,
	endpoint string, opts []grpc.DialOption) error {
	return proto.RegisterDocumentsHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
