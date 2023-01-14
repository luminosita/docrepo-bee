package http

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/luminosita/common-bee/pkg/log"
	documents "github.com/luminosita/docrepo-bee/internal/interface/app/documents"
	http2 "net/http"
)

var GetDocumentWireSet = wire.NewSet(NewPutDocumentHandler)

//	wire.Bind(new(handlers.Handler), new(*PutDocumentHandler)))

type GetDocumentHandler struct {
	cd documents.GetDocumenter
}

func NewGetDocumentHandler(cd documents.GetDocumenter) *GetDocumentHandler {
	return &GetDocumentHandler{
		cd: cd,
	}
}

// GetDocument godoc
// @Summary      Retrieves document for repository
// @Description  Get Document By ID
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Document ID"
// @Success      200  {object}  string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /documents/{id} [get]
func (h *GetDocumentHandler) Handle(ctx http.Context) (err error) {
	documentId := ctx.Vars().Get("id")
	log.Infof("GetDocumentHandler %s", documentId)

	res, err := h.cd.Execute(&documents.GetDocumenterRequest{
		DocumentId: documentId,
	})
	if err != nil {
		return err
	}

	ctx.Response().Header().Set("Content-Disposition", "attachment; filename="+res.Name)

	return ctx.Stream(http2.StatusOK, "binary/octet-stream", res.Reader)
}
