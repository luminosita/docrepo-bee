package documents

import (
	"github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
	"github.com/luminosita/honeycomb/pkg/http/ctx"
	"github.com/luminosita/honeycomb/pkg/log"
)

const DOCUMENT_FIELD_KEY = "document"

type PutDocumentHandler struct {
	cd documents.PutDocumenter
}

func NewPutDocumentHandler(cd documents.PutDocumenter) *PutDocumentHandler {
	return &PutDocumentHandler{
		cd: cd,
	}
}

// PutDocumenter godoc
// @Summary      Put Document
// @Description  Stores document into repository
// @Tags         documents
// @Accept       json
// @Produce      json
// @Success      200  {object}  string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /documents [post]
func (h *PutDocumentHandler) Handle(ctx *ctx.Ctx) (err error) {
	formFile, err := ctx.FormFile(DOCUMENT_FIELD_KEY)
	if err != nil {
		return err
	}

	name := formFile.Filename

	log.Log().Infof("PutDocumentHandler: %s", name)

	file, err := formFile.Open()
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()

	res, err := h.cd.Execute(&documents.PutDocumenterRequest{
		Name:   name,
		Size:   formFile.Size,
		Reader: file,
	})
	if err != nil {
		return err
	}

	return ctx.SendString(res.DocumentId)
}
