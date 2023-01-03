package handlers

import (
	"context"
	"github.com/luminosita/honeycomb/pkg/http/handlers"
	use_cases "github.com/luminosita/sample-bee/internal/bee/factories/use-cases"
	"github.com/luminosita/sample-bee/internal/infra/http/handlers/documents"
)

func MakeGetDocumentHandler(ctx context.Context) handlers.Handler {
	useCase := use_cases.MakeGetDocument(ctx)
	return documents.NewGetDocumentHandler(useCase)
}

func MakeGetAllDocumentsHandler(ctx context.Context) handlers.Handler {
	useCase := use_cases.MakeGetAllDocuments(ctx)
	return documents.NewGetAllDocumentsHandler(useCase)
}

func MakeCreateDocumentHandler(ctx context.Context) handlers.Handler {
	useCase := use_cases.MakeCreateDocument(ctx)
	return documents.NewCreateDocumentHandler(useCase)
}
