package use_cases

import (
	"context"
	documents2 "github.com/luminosita/bee/internal/app/use-cases/documents"
	"github.com/luminosita/bee/internal/infra/db/mockdb/repositories"
	"github.com/luminosita/bee/internal/interfaces/use-cases/documents"
)

func MakeCreateDocument(ctx context.Context) documents.CreateDocumenter {
	docRepo := repositories.MakeMockDocumentsRepository(ctx)
	return documents2.NewCreateDocument(docRepo)
}

func MakeGetAllDocuments(ctx context.Context) documents.GetAllDocumenter {
	docRepo := repositories.MakeMockDocumentsRepository(ctx)
	return documents2.NewGetAllDocuments(docRepo)
}

func MakeGetDocument(ctx context.Context) documents.GetDocumenter {
	docRepo := repositories.MakeMockDocumentsRepository(ctx)
	return documents2.NewGetDocument(docRepo)
}
