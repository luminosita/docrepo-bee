package documents

import usecases "github.com/luminosita/bee/internal/interfaces/use-cases"

type Request = struct {
	UserId string
	Body   string
}
type Response = any

type CreateDocumenter interface {
	usecases.UseCaser[*Request, *Response]
}
