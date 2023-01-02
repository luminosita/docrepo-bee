package server

import (
	"context"
	"github.com/luminosita/bee/common/http/handlers"
)

type Method int

const (
	GET   Method = iota // Head = 0
	POST                // Shoulder = 1
	PUT                 // Knee = 2
	PATCH               // Toe = 3
)

func (m Method) String() string {
	return []string{"GET", "PUT", "HEAD", "PATCH"}[m]
}

type Route struct {
	Method  Method
	Path    string
	Handler handlers.Handler
}

type Server interface {
	Routes(ctx context.Context) []*Route
}
