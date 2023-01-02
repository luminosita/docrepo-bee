package documents

import (
    "github.com/luminosita/bee/internal/infra/http"
    "github.com/luminosita/bee/internal/infra/http/handlers"
    "github.com/luminosita/bee/internal/infra/http/interfaces"
    "github.com/luminosita/bee/internal/interfaces/use-cases/documents"
)

type CreateDocumentHandler struct {
    cd documents.CreateDocumenter
    cv interfaces.Validation[*http.HttpRequest]
    *handlers.BaseHandler
}

func (h *CreateDocumentHandler) execute(req *http.HttpRequest) (*http.HttpResponse, error) {
    userId := req.UserId
    body := req.Body

    res, err := h.cd.Execute(&documents.Request{
        UserId: userId,
        Body:   body,
    })

    if err != nil {
        return nil, err
    }

    return &http.HttpResponse{
        StatusCode: 200,
        Body:       res,
    }, nil
}
