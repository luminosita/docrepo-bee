package handlers

import (
	"context"
	"github.com/luminosita/bee/internal/infra/http/handlers"
	"github.com/luminosita/bee/internal/infra/http/handlers/documents"
	validatoradapters "github.com/luminosita/bee/internal/infra/http/validator-adapters"
	documents2 "github.com/luminosita/bee/internal/server/factories/use-cases/documents"
)

func MakeCreateDocumentHandler(ctx context.Context) *handlers.BaseHandler {
	useCase := documents2.MakeCreateDocument(ctx)
	validation := validatoradapters.NewValidatorAdapter()
	return documents.NewCreateDocumentHandler(useCase, validation)
}
