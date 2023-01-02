package interfaces

import "github.com/luminosita/bee/internal/infra/http"

type Handler interface {
	Process(req *http.HttpRequest) (*http.HttpResponse, error)
}
