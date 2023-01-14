package documents

import (
	"github.com/google/wire"
	documents2 "github.com/luminosita/docrepo-bee/internal/interface/app/documents"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewGetDocument, wire.Bind(new(documents2.GetDocumenter), new(*GetDocument)),
	NewPutDocument, wire.Bind(new(documents2.PutDocumenter), new(*PutDocument)),
	NewGetDocumentInfo, wire.Bind(new(documents2.GetDocumentInfoer), new(*GetDocumentInfo)))
