package documents

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
