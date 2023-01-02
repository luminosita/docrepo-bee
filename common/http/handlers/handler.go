package handlers

import "github.com/luminosita/bee/common/http"

type Handler interface {
	Handle(req *http.HttpRequest) (*http.HttpResponse, error)
}
