package documents

import (
	"github.com/google/wire"
	"github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
	"github.com/luminosita/honeycomb/pkg/http/ctx"
	"github.com/luminosita/honeycomb/pkg/http/handlers"
	"github.com/luminosita/honeycomb/pkg/log"
)

var GetDocumentInfoWireSet = wire.NewSet(NewGetDocumentInfoHandler,
	wire.Bind(new(handlers.Handler), new(*GetDocumentInfoHandler)))

type GetDocumentInfoHandler struct {
	cd documents.GetDocumentInfoer
}

func NewGetDocumentInfoHandler(cd documents.GetDocumentInfoer) *GetDocumentInfoHandler {
	return &GetDocumentInfoHandler{
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
func (h *GetDocumentInfoHandler) Handle(ctx *ctx.Ctx) error {
	documentId := ctx.Params["id"]
	log.Log().Infof("GetDocumentInfoHandler %s", documentId)

	res, err := h.cd.Execute(&documents.GetDocumentInfoerRequest{
		DocumentId: documentId,
	})
	if err != nil {
		return err
	}

	return ctx.SendResponse(res)
}
