package documents

type RepoRequest = struct {
}
type RepoResponse = any

type CreateDocumentRepositorer interface {
	CreateDocument(req *RepoRequest) (*RepoResponse, error)
}
