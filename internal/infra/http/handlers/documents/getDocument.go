package documents

import (
	"github.com/google/wire"
	"github.com/luminosita/common-bee/pkg/log"
	"github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
	"github.com/luminosita/honeycomb/pkg/http/ctx"
	"github.com/luminosita/honeycomb/pkg/http/handlers"
)

var GetDocumentWireSet = wire.NewSet(NewGetDocumentHandler,
	wire.Bind(new(handlers.Handler), new(*GetDocumentHandler)))

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
func (h *GetDocumentHandler) Handle(ctx *ctx.Ctx) error {
	documentId := ctx.Params["id"]
	log.Log().Infof("GetDocumentHandler %s", documentId)

	res, err := h.cd.Execute(&documents.GetDocumenterRequest{
		DocumentId: documentId,
	})
	if err != nil {
		return err
	}

	return ctx.SendStream(res.Name, res.Reader, int(res.Size))
}
