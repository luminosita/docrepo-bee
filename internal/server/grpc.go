package server

import (
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	v1 "github.com/luminosita/docrepo-bee/api/documents/v1"
	"github.com/luminosita/docrepo-bee/internal/conf"
	"github.com/luminosita/docrepo-bee/internal/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, svc *service.Documents) *grpc.Server {
	//authz := opaauthz.NewOpaAuthorizer(opaauthz.OpaURL(
	//	"http://localhost:8181/v1/data/apis/invocation_allowed"))

	claims := jwtv4.MapClaims{}

	var opts = []grpc.ServerOption{
		grpc.Middleware(
			jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(SecretKey), nil
			}, jwt.WithSigningMethod(jwtv4.SigningMethodHS256),
				jwt.WithClaims(func() jwtv4.Claims { return claims })),
			//			authz.OpaMiddleware,
			validate.Validator(),
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterDocumentsServer(srv, svc)
	return srv
}
