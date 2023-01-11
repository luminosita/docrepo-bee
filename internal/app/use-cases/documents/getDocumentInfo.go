package documents

import (
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/luminosita/docrepo-bee/internal/interfaces/repositories/documents"
	documents2 "github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
)

var GetDocumentInfoWireSet = wire.NewSet(NewGetDocumentInfo,
	wire.Bind(new(documents2.GetDocumentInfoer), new(*GetDocumentInfo)))

type GetDocumentInfo struct {
	repo documents.GetDocumentInfoRepositorer
}

func NewGetDocumentInfo(r documents.GetDocumentInfoRepositorer) *GetDocumentInfo {
	return &GetDocumentInfo{
		repo: r,
	}
}

func (d *GetDocumentInfo) Execute(
	docData *documents2.GetDocumentInfoerRequest) (*documents2.GetDocumentInfoerResponse, error) {

	if docData == nil || len(docData.DocumentId) == 0 {
		//TODO: Externalize
		return nil, errors.New(fmt.Sprintf("Bad request: %+v", docData))
	}

	data := &documents.GetDocumentInfoRepositorerRequest{
		DocumentId: docData.DocumentId,
	}

	res, err := d.repo.GetDocumentInfo(data)

	if err != nil {
		return nil, err
	}

	return &documents2.GetDocumentInfoerResponse{
		Name:       res.Name,
		Size:       res.Size,
		UploadDate: res.UploadDate,
	}, nil
}
