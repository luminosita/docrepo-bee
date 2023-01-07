package handlers

import (
	"context"
	use_cases "github.com/luminosita/docrepo-bee/internal/bee/factories/use-cases"
	"github.com/luminosita/docrepo-bee/internal/infra/http/handlers/documents"
	"github.com/luminosita/honeycomb/pkg/http/handlers"
)

func MakeGetDocumentHandler(ctx context.Context) handlers.Handler {
	useCase := use_cases.MakeGetDocument(ctx)
	return documents.NewGetDocumentHandler(useCase).Handle
}

func MakePutDocumentHandler(ctx context.Context) handlers.Handler {
	useCase := use_cases.MakePutDocument(ctx)
	return documents.NewPutDocumentHandler(useCase).Handle
}
