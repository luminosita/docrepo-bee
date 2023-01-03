package documents

import (
	"github.com/luminosita/honeycomb/pkg/http"
	"github.com/luminosita/sample-bee/internal/interfaces/use-cases/documents"
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
// @Summary      Get All Documents
// @Description  dummy GET method
// @Tags         documents
// @Accept       json
// @Produce      json
// @Success      200  {object}  string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /documents [get]
func (h *GetAllDocumentsHandler) Handle(_ *http.HttpRequest) (*http.HttpResponse, error) {
	res, err := h.cd.Execute(&documents.GetAllDocumenterRequest{})

	if err != nil {
		return nil, err
	}

	return http.Ok(res.Documents), nil
}
