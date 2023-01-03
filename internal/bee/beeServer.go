package bee

import (
	"context"
	"github.com/luminosita/bee/common/server"
	"github.com/luminosita/bee/internal/bee/factories/handlers"
)

type Config struct {
	Sc server.Config `mapstructure:"server"`
}

func (c *Config) ServerConfig() *server.Config {
	return &c.Sc
}

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

func (*BeeServer) Routes(ctx context.Context) []*server.Route {
	routes := make([]*server.Route, 0)

	//router.get('/posts', expressRouteAdapter(makeGetLatestPostsController()))
	//router.get('/posts/:id', expressRouteAdapter(makeGetPostByIdController()))
	//router.post('/posts', authMiddleware, expressRouteAdapter(makeCreatePostController()))
	//router.patch('/posts/:id', authMiddleware, expressRouteAdapter(makeUpdatePostController()))
	//router.delete('/posts/:id', authMiddleware, expressRouteAdapter(makeDeletePostController()))

	routes = append(routes, &server.Route{
		Method: server.GET, Path: "/documents/:id", Handler: handlers.MakeGetDocumentHandler(ctx)})
	routes = append(routes, &server.Route{
		Method: server.GET, Path: "/documents", Handler: handlers.MakeGetAllDocumentsHandler(ctx)})
	routes = append(routes, &server.Route{
		Method: server.POST, Path: "/documents", Handler: handlers.MakeCreateDocumentHandler(ctx)})

	return routes
}
