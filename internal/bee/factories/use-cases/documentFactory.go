package use_cases

import (
	"context"
	documents2 "github.com/luminosita/docrepo-bee/internal/app/use-cases/documents"
	"github.com/luminosita/docrepo-bee/internal/infra/db/mongodb/repositories"
	"github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
)

func MakePutDocument(ctx context.Context) documents.PutDocumenter {
	docRepo := repositories.NewPutDocumentRepository(ctx)
	return documents2.NewPutDocument(docRepo)
}

func MakeGetDocument(ctx context.Context) documents.GetDocumenter {
	docRepo := repositories.NewGetDocumentRepository(ctx)
	return documents2.NewGetDocument(docRepo)
}
