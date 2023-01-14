package http

import (
	"bufio"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/luminosita/common-bee/pkg/log"
	"github.com/luminosita/docrepo-bee/internal/interface/app/documents"
	"io"
	http2 "net/http"
)

var PutDocumentWireSet = wire.NewSet(NewPutDocumentHandler)

//	wire.Bind(new(handlers.Handler), new(*PutDocumentHandler)))

const DocumentFieldKey = "document"

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
func (h *PutDocumentHandler) Handle(ctx http.Context) (err error) {
	err = ctx.Request().ParseMultipartForm(2048)
	if err != nil {
		return log.LogError(err)
	}

	formFile, ok := ctx.Request().MultipartForm.File[DocumentFieldKey]
	if !ok {
		return log.LogErrorf("unable to locate file with key (%s)", DocumentFieldKey)
	}

	name := formFile[0].Filename

	log.Infof("PutDocumentHandler: %s", name)

	file, err := formFile[0].Open()
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()

	res, err := h.cd.Execute(&documents.PutDocumenterRequest{
		Name: name,
		Size: formFile[0].Size,
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

	return ctx.String(http2.StatusOK, res.DocumentId)
}

/*
func SendString(body string) error {
	return ctx.fCtx.Status(fiber.StatusOK).JSON(&JsonResponse{
		Body: body,
	})
}

func (ctx *Ctx) SendResponse(obj any) error {
	return ctx.fCtx.Status(fiber.StatusOK).JSON(obj)
}

func (ctx *Ctx) FormFile(key string) (*multipart.FileHeader, error) {
	return ctx.fCtx.FormFile(key)
}

func (ctx *Ctx) SendStream(filename string, reader io.Reader, size ...int) error {
	ctx.fCtx.Attachment(filename)

	if len(size) > 0 && size[0] >= 0 {
		return ctx.fCtx.SendStream(reader, size[0])
	} else {
		return ctx.fCtx.SendStream(reader)
	}

}
*/
