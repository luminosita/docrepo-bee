package documents

import (
	"fmt"
	"github.com/luminosita/common-bee/pkg/log"
	documents2 "github.com/luminosita/docrepo-bee/internal/interface/app/documents"
	"github.com/luminosita/docrepo-bee/internal/interface/data/documents"
)

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
		return nil, log.LogErrorf(fmt.Sprintf("Bad request: %+v", docData))
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
