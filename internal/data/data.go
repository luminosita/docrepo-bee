package data

import (
	"github.com/google/wire"
	"github.com/luminosita/docrepo-bee/internal/data/db/mongodb"
	"github.com/luminosita/docrepo-bee/internal/data/db/mongodb/documents"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(mongodb.NewData, documents.ProviderSet)
