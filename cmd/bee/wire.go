//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"
	"github.com/luminosita/common-bee/pkg/log"
	"github.com/luminosita/docrepo-bee/internal/app"
	"github.com/luminosita/docrepo-bee/internal/conf"
	"github.com/luminosita/docrepo-bee/internal/data"
	"github.com/luminosita/docrepo-bee/internal/server"
	"github.com/luminosita/docrepo-bee/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(context.Context, *conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, app.ProviderSet, service.ProviderSet, newApp))
}
