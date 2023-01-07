//go:build wireinject
// +build wireinject

package bee

import (
	"context"
	"github.com/google/wire"
	"github.com/luminosita/docrepo-bee/internal/app/use-cases/documents"
	"github.com/luminosita/docrepo-bee/internal/infra/db/mongodb/repositories"
	documents2 "github.com/luminosita/docrepo-bee/internal/infra/http/handlers/documents"
)

func MakeGetDocumentHandler(ctx context.Context) *documents2.GetDocumentHandler {
	wire.Build(documents2.GetWireSet, documents.GetWireSet, repositories.GetWireSet)

	return nil
}

func MakePutDocumentHandler(ctx context.Context) *documents2.PutDocumentHandler {
	wire.Build(documents2.PutWireSet, documents.PutWireSet, repositories.PutWireSet)

	return nil
}
