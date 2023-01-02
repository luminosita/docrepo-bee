package bee

import (
	"context"
	"github.com/luminosita/bee/common/server"
	"github.com/luminosita/bee/internal/bee/factories/handlers"
)

type BeeServer struct {
}

func NewBeeServer() *BeeServer {
	return &BeeServer{}
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
