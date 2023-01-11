package documents

import (
	"bufio"
	"github.com/google/wire"
	"github.com/luminosita/common-bee/pkg/log"
	"github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
	"github.com/luminosita/honeycomb/pkg/http/ctx"
	"github.com/luminosita/honeycomb/pkg/http/handlers"
	"io"
)

var PutDocumentWireSet = wire.NewSet(NewPutDocumentHandler,
	wire.Bind(new(handlers.Handler), new(*PutDocumentHandler)))

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
		Name: name,
		Size: formFile.Size,
	})
	if err != nil {
		return err
	}

	defer func() {
		err = res.Writer.Close()
	}()

	reader := bufio.NewReader(file)

	buf := make([]byte, 3*1024) //the chunk size

	_, err = io.CopyBuffer(res.Writer, reader, buf)
	if err != nil {
		return err
	}

	return ctx.SendString(res.DocumentId)
}