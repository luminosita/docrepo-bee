package app

import (
	"github.com/google/wire"
	"github.com/luminosita/docrepo-bee/internal/app/documents"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(documents.ProviderSet)
