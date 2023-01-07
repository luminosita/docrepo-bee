package documents

import (
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/luminosita/docrepo-bee/internal/interfaces/repositories/documents"
	documents2 "github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
)

var PutWireSet = wire.NewSet(NewPutDocument,
	wire.Bind(new(documents2.PutDocumenter), new(*PutDocument)))

type PutDocument struct {
	repo documents.PutDocumentRepositorer
}

func NewPutDocument(r documents.PutDocumentRepositorer) *PutDocument {
	return &PutDocument{
		repo: r,
	}
}

func (d *PutDocument) Execute(
	docData *documents2.PutDocumenterRequest) (*documents2.PutDocumenterResponse, error) {

	if docData == nil || len(docData.Name) == 0 || docData.Size <= 0 || docData.Reader == nil {
		return nil, errors.New(fmt.Sprintf("Bad request: %+v", docData))
	}

	data := &documents.PutDocumentRepositorerRequest{
		Name:   docData.Name,
		Size:   docData.Size,
		Reader: docData.Reader,
	}

	res, err := d.repo.PutDocument(data)
	if err != nil {
		return nil, err
	}

	return &documents2.PutDocumenterResponse{
		DocumentId: res.DocumentId,
	}, nil
}
