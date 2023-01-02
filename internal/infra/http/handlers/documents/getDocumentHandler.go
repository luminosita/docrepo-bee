package documents

import (
	"github.com/luminosita/bee/common/http"
	"github.com/luminosita/bee/internal/domain/entities"
	"github.com/luminosita/bee/internal/interfaces/use-cases/documents"
)

type GetDocumentHandler struct {
	cd documents.GetDocumenter
}

func NewGetDocumentHandler(cd documents.GetDocumenter) *GetDocumentHandler {
	return &GetDocumentHandler{
		cd: cd,
	}
}

// GetDocument godoc
// @Summary      Show something
// @Description  dummy GET method
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Document ID"
// @Success      200  {object}  string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /documents/:id [get]
func (h *GetDocumentHandler) Handle(req *http.HttpRequest) (*http.HttpResponse, error) {
	documentId := req.Params["docId"]

	res, err := h.cd.Execute(&documents.GetDocumenterRequest{
		DocumentId: documentId,
	})

	if err != nil {
		return nil, err
	}

	http.Ok(res.Content)

	return nil, nil
}

func (h *GetDocumentHandler) Model(req *http.HttpRequest) *entities.Document {
	return nil
}
