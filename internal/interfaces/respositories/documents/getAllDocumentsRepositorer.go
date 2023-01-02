package documents

type GetDocumentRepositorerRequest = struct {
	DocumentID string
}
type GetDocumentRepositorerResponse = struct {
	Content string
}

type GetDocumentRepositorer interface {
	GetDocument(req *GetDocumentRepositorerRequest) (*GetDocumentRepositorerResponse, error)
}
