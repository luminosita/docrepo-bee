package documents

import (
	"github.com/luminosita/bee/common/http"
	"github.com/luminosita/bee/internal/domain/entities"
	"github.com/luminosita/bee/internal/interfaces/use-cases/documents"
)

type GetAllDocumentsHandler struct {
	cd documents.GetAllDocumenter
}

func NewGetAllDocumentsHandler(cd documents.GetAllDocumenter) *GetAllDocumentsHandler {
	return &GetAllDocumentsHandler{
		cd: cd,
	}
}

// GetAllDocuments godoc
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
// @Router       /documents [get]
func (h *GetAllDocumentsHandler) Handle(req *http.HttpRequest) (*http.HttpResponse, error) {
	documentId := req.Params["docId"]

	res, err := h.cd.Execute(&documents.GetAllDocumenterRequest{
		DocumentId: documentId,
	})

	if err != nil {
		return nil, err
	}

	http.Ok(res.Content)

	return nil, nil
}

func (h *GetAllDocumentsHandler) Model(req *http.HttpRequest) *entities.Document {
	return nil
}
