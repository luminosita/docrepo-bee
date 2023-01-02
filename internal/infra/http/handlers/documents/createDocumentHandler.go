package documents

import (
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
// @Summary      Show something
// @Description  dummy POST method
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /documents [post]
func (h *CreateDocumentHandler) Handle(req *http.HttpRequest) (*http.HttpResponse, error) {
	body := string(req.Body)

	_, err := h.cd.Execute(&documents.CreateDocumenterRequest{
		Content: body,
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *CreateDocumentHandler) Model(req *http.HttpRequest) *entities.Document {
	return nil
}
