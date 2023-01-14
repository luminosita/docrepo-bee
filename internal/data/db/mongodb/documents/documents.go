package documents

import (
	"github.com/google/wire"
	"github.com/luminosita/docrepo-bee/internal/interface/data/documents"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewGetDocumentInfoRepository, wire.Bind(new(documents.GetDocumentInfoRepositorer), new(*GetDocumentInfoRepository)),
	NewPutDocumentRepository, wire.Bind(new(documents.PutDocumentRepositorer), new(*PutDocumentRepository)),
	NewGetDocumentRepository, wire.Bind(new(documents.GetDocumentRepositorer), new(*GetDocumentRepository)))

const Documents = "documents"
