package documents

import (
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/luminosita/docrepo-bee/internal/interfaces/repositories/documents"
	documents2 "github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
)

var GetDocumentWireSet = wire.NewSet(NewGetDocument,
	wire.Bind(new(documents2.GetDocumenter), new(*GetDocument)))

type GetDocument struct {
	repo documents.GetDocumentRepositorer
}

func NewGetDocument(r documents.GetDocumentRepositorer) *GetDocument {
	return &GetDocument{
		repo: r,
	}
}

func (d *GetDocument) Execute(
	docData *documents2.GetDocumenterRequest) (*documents2.GetDocumenterResponse, error) {

	if docData == nil || len(docData.DocumentId) == 0 {
		//TODO: Externalize
		return nil, errors.New(fmt.Sprintf("Bad request: %+v", docData))
	}

	data := &documents.GetDocumentRepositorerRequest{
		DocumentId: docData.DocumentId,
	}

	res, err := d.repo.GetDocument(data)

	if err != nil {
		return nil, err
	}

	return &documents2.GetDocumenterResponse{
		Name:   res.Name,
		Size:   res.Size,
		Reader: res.Reader,
	}, nil
}
