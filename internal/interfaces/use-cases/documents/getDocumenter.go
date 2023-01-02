package documents

import "github.com/luminosita/bee/common/interfaces"

type Request = struct {
	UserId string
	Body   string
}
type Response = any

type CreateDocumenter interface {
	interfaces.UseCaser[*Request, *Response]
}
