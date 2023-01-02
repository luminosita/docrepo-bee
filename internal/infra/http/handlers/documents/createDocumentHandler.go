package documents

import (
	"encoding/json"
	"github.com/luminosita/bee/common/http"
	"github.com/luminosita/bee/internal/domain/entities"
	"github.com/luminosita/bee/internal/interfaces/use-cases/documents"
)

type CreateDocumentHandler struct {
	cd documents.CreateDocumenter
}

func NewCreateDocumentHandler(cd documents.CreateDocumenter) *CreateDocumentHandler {
	return &CreateDocumentHandler{
		cd: cd,
	}
}

// CreateDocument godoc
// @Summary      Create Document
// @Description  dummy POST method
// @Tags         documents
// @Accept       json
// @Produce      json
// @Success      200  {object}  string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /documents [post]
func (h *CreateDocumentHandler) Handle(req *http.HttpRequest) (*http.HttpResponse, error) {
	doc := &entities.Document{}

	json.Unmarshal(req.Body, doc)

	res, err := h.cd.Execute(&documents.CreateDocumenterRequest{
		Document: doc,
	})

	if err != nil {
		return nil, err
	}

	return http.Ok(res.DocumentId), nil
}
