package use_cases

import (
	"context"
	documents2 "github.com/luminosita/bee/internal/app/use-cases/documents"
	"github.com/luminosita/bee/internal/infra/db/mongodb/repositories"
	"github.com/luminosita/bee/internal/interfaces/use-cases/documents"
)

func MakeCreateDocument(ctx context.Context) documents.CreateDocumenter {
	docRepo := repositories.NewCreateDocumentRepository(ctx)
	return documents2.NewCreateDocument(docRepo)
}

func MakeGetAllDocuments(ctx context.Context) documents.GetAllDocumenter {
	docRepo := repositories.NewGetAllDocumentsRepository(ctx)
	return documents2.NewGetAllDocuments(docRepo)
}

func MakeGetDocument(ctx context.Context) documents.GetDocumenter {
	docRepo := repositories.NewGetDocumentRepository(ctx)
	return documents2.NewGetDocument(docRepo)
}
