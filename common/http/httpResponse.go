package http

type HttpResponse struct {
	StatusCode int
	Body       string

	Errors []error
}
