//go:build wireinject
// +build wireinject

package bee

import (
	"context"
	"github.com/google/wire"
	"github.com/luminosita/docrepo-bee/internal/app/use-cases/documents"
	"github.com/luminosita/docrepo-bee/internal/infra/db/mongodb/repositories"
	"github.com/luminosita/docrepo-bee/internal/infra/grpc/handlers"
	documents2 "github.com/luminosita/docrepo-bee/internal/infra/http/handlers/documents"
)

func MakeGetDocumentHandler(ctx context.Context) *documents2.GetDocumentHandler {
	wire.Build(documents2.GetDocumentWireSet, documents.GetDocumentWireSet,
		repositories.GetDocumentWireSet)

	return nil
}

func MakeGetDocumentInfoHandler(ctx context.Context) *documents2.GetDocumentInfoHandler {
	wire.Build(documents2.GetDocumentInfoWireSet, documents.GetDocumentInfoWireSet,
		repositories.GetDocumentInfoWireSet)

	return nil
}

func MakePutDocumentHandler(ctx context.Context) *documents2.PutDocumentHandler {
	wire.Build(documents2.PutDocumentWireSet, documents.PutDocumentWireSet,
		repositories.PutDocumentWireSet)

	return nil
}

func MakeGrpcDocumentServer(ctx context.Context) *handlers.DocumentsServer {
	wire.Build(handlers.PbDocumentsWireSet, documents.GetDocumentInfoWireSet,
		repositories.GetDocumentInfoWireSet, documents.PutDocumentWireSet,
		repositories.PutDocumentWireSet, documents.GetDocumentWireSet,
		repositories.GetDocumentWireSet)

	return nil
}
