package grpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/luminosita/common-bee/pkg/utils"
	pb "github.com/luminosita/docrepo-bee/api/documents/v1"
	"github.com/luminosita/docrepo-bee/internal/interface/app/documents"
	"github.com/luminosita/docrepo-bee/internal/interface/app/documents/mocks"
	server2 "github.com/luminosita/docrepo-bee/internal/server"
	"github.com/luminosita/docrepo-bee/internal/service"
	"github.com/stretchr/testify/assert"
)

const addr = "localhost:9999"

type Mock struct {
	t    *testing.T
	ctrl *gomock.Controller

	pd  *mocks.MockPutDocumenter
	gd  *mocks.MockGetDocumenter
	gdi *mocks.MockGetDocumentInfoer

	c *Client
}

func newMock(t *testing.T) (m *Mock) {
	m = &Mock{}
	m.t = t
	m.ctrl = gomock.NewController(t)

	return
}

func setupServer(m *Mock, middleware ...middleware.Middleware) func() {
	server := service.NewDocumentsServer(m.gdi, m.gd, m.pd)

	var opts = []kgrpc.ServerOption{
		kgrpc.Middleware(
			middleware...,
		),
	}

	opts = append(opts, kgrpc.Address(addr))

	gs := kgrpc.NewServer(opts...)
	pb.RegisterDocumentsServer(gs, server)

	app := kratos.New(
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			gs,
		),
	)

	go func() {
		if err := app.Run(); err != nil {
			log.Printf("grpcServer.Serve: %v", err)
		}
	}()

	return func() {
		_ = app.Stop()
	}
}

func setupClient(m *Mock, middleware ...middleware.Middleware) func() {
	conn, err := kgrpc.DialInsecure(
		context.Background(),
		kgrpc.WithEndpoint(addr),
		kgrpc.WithMiddleware(
			middleware...,
		),
	)
	if err != nil {
		log.Fatalf("grpc.DialContext: %v", err)
	}

	m.c = NewClient(pb.NewDocumentsClient(conn))

	return func() {
		_ = conn.Close()
	}
}

func setupTest(m *Mock, clientMid []middleware.Middleware, serverMid []middleware.Middleware) func() {
	if m == nil {
		panic("Mock not initialized")
	}

	m.pd = mocks.NewMockPutDocumenter(m.ctrl)
	m.gd = mocks.NewMockGetDocumenter(m.ctrl)
	m.gdi = mocks.NewMockGetDocumentInfoer(m.ctrl)

	cleanupServer := setupServer(m, serverMid...)
	cleanupClient := setupClient(m, clientMid...)

	return func() {
		cleanupServer()
		cleanupClient()
	}
}

func setupWithMiddleware(m *Mock, claims jwtv4.MapClaims, testMiddleware middleware.Middleware) func() {
	serverMid := []middleware.Middleware{
		jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
			return []byte(server2.SecretKey), nil
		}, jwt.WithSigningMethod(jwtv4.SigningMethodHS256),
			jwt.WithClaims(func() jwtv4.Claims { return claims })),
		testMiddleware,
	}

	clientMid := []middleware.Middleware{
		jwt.Client(func(token *jwtv4.Token) (interface{}, error) {
			return []byte(server2.SecretKey), nil
		}, jwt.WithSigningMethod(jwtv4.SigningMethodHS256),
			jwt.WithClaims(func() jwtv4.Claims { return claims })),
	}

	return setupTest(m, clientMid, serverMid)
}

func callGetDocumentInfo(t *testing.T, m *Mock) {
	ttt := time.Unix(100, 10)

	repoRequest := &documents.GetDocumentInfoerRequest{DocumentId: "123"}
	repoResponse := &documents.GetDocumentInfoerResponse{Name: "laza", Size: 1234, UploadDate: ttt}

	m.gdi.EXPECT().Execute(gomock.Eq(repoRequest)).Return(repoResponse, nil)

	docInfo, err := m.c.GetDocumentInfo(context.Background(), "123")

	assert.Nil(t, err)
	assert.NotNil(t, docInfo)
	assert.Equal(t, "laza", docInfo.Name)
	assert.Equal(t, int64(1234), docInfo.Size)
	assert.Equal(t, ttt.UTC(), docInfo.UploadDate)
}

func callGetDocument(t *testing.T, m *Mock) {
	s := "test test test"

	reader := io.NopCloser(strings.NewReader(s))

	repoRequest := &documents.GetDocumenterRequest{DocumentId: "123"}
	repoResponse := &documents.GetDocumenterResponse{Name: "laza", Size: 1234, Reader: reader}

	m.gd.EXPECT().Execute(gomock.Eq(repoRequest)).Return(repoResponse, nil)

	docInfo, r, err := m.c.GetDocument(context.Background(), "123")
	defer func() { _ = r.Close() }()

	assert.Nil(t, err)
	assert.NotNil(t, docInfo)
	assert.Equal(t, "laza", docInfo.Name)
	assert.Equal(t, int64(1234), docInfo.Size)
	assert.NotNil(t, r)

	resBytes, _ := io.ReadAll(r)
	assert.Equal(t, []byte(s), resBytes)
}

func TestGetDocumentInfo(t *testing.T) {
	m := newMock(t)

	serverMid := []middleware.Middleware{
		func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
				return handler(ctx, req)
			}
		},
	}

	defer setupTest(m, nil, serverMid)()

	callGetDocumentInfo(t, m)
}

func TestGetDocumentInfoBad(t *testing.T) {
	m := newMock(t)

	serverMid := []middleware.Middleware{validate.Validator()}

	defer setupTest(m, nil, serverMid)()

	_, err := m.c.GetDocumentInfo(context.Background(), "")

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "value length must be at least")
}

func TestGetDocument(t *testing.T) {
	m := newMock(t)
	defer setupTest(m, nil, nil)()

	callGetDocument(t, m)
}

func TestPutDocument(t *testing.T) {
	m := newMock(t)
	defer setupTest(m, nil, nil)()

	s := "some test data"

	buffer := &bytes.Buffer{}
	writer := utils.NewNopWriteCloser(buffer)
	docInfo := &DocumentInfo{Name: "laza", Size: 1234}

	repoRequest := &documents.PutDocumenterRequest{Name: "laza", Size: 1234}
	repoResponse := &documents.PutDocumenterResponse{
		DocumentId: "123", Writer: writer}

	m.pd.EXPECT().Execute(gomock.Eq(repoRequest)).Return(repoResponse, nil)

	reader := io.NopCloser(strings.NewReader(s))

	docId, err := m.c.PutDocument(context.Background(), docInfo, reader)

	assert.Nil(t, err)
	assert.NotNil(t, "123", docId)
	assert.Equal(t, []byte(s), buffer.Bytes())
}

func TestJwt(t *testing.T) {
	m := newMock(t)

	tm := false

	testMiddleware := func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			c, ok := jwt.FromContext(ctx)
			assert.True(t, ok)
			assert.Equal(t, "laza", c.(jwtv4.MapClaims)["sub"])
			assert.Equal(t, "123", c.(jwtv4.MapClaims)["id"])

			tm = true

			return handler(ctx, req)
		}
	}

	claims := jwtv4.MapClaims{
		"sub": "laza",
		"id":  "123",
	}

	defer setupWithMiddleware(m, claims, testMiddleware)()

	callGetDocumentInfo(t, m)

	assert.True(t, tm) //check if test method got called
}

func TestJwtForStream(t *testing.T) {
	m := newMock(t)

	tm := false

	testMiddleware := func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			c, ok := jwt.FromContext(ctx)
			assert.True(t, ok)
			assert.Equal(t, "laza", c.(jwtv4.MapClaims)["sub"])
			assert.Equal(t, "123", c.(jwtv4.MapClaims)["id"])

			tm = true

			return handler(ctx, req)
		}
	}

	claims := jwtv4.MapClaims{
		"sub": "laza",
		"id":  "123",
	}

	defer setupWithMiddleware(m, claims, testMiddleware)()

	callGetDocument(t, m)

	assert.True(t, tm) //check if test method got called
}

func TestJwtBad(t *testing.T) {
	m := newMock(t)

	tm := false

	testMiddleware := func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			c, ok := jwt.FromContext(ctx)
			assert.True(t, ok)
			assert.Empty(t, c.(jwtv4.MapClaims))

			tm = true

			return nil, errors.New("jwt bad test")
		}
	}

	claims := jwtv4.MapClaims{}

	defer setupWithMiddleware(m, claims, testMiddleware)()

	_, err := m.c.GetDocumentInfo(context.Background(), "123456")

	assert.NotNil(t, err)
	assert.True(t, tm) //check if test method got called
	assert.ErrorContains(t, err, "jwt bad test")
}
