package grpc

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/luminosita/common-bee/pkg/utils"
	pb "github.com/luminosita/docrepo-bee/api/gen/v1"
	"github.com/luminosita/docrepo-bee/internal/infra/grpc/handlers"
	"github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
	"github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"
)

const addr = "localhost:9999"

var (
	server *handlers.DocumentsServer
)

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

func setupServer(server *handlers.DocumentsServer, done chan os.Signal) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterDocumentsServer(grpcServer, server)
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("grpcServer.Serve: %v", err)
		}
	}()

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	grpcServer.GracefulStop()
	log.Printf("Server stopped")
}

func TestMain(m *testing.M) {
	done := make(chan os.Signal, 1)

	server = handlers.NewDocumentsServer(nil, nil, nil)

	go func() {
		setupServer(server, done)
	}()

	time.Sleep(1)

	code := m.Run()

	done <- os.Interrupt

	os.Exit(code)
}

func setupTest(m *Mock) func() {
	if m == nil {
		panic("Mock not initialized")
	}

	m.pd = mocks.NewMockPutDocumenter(m.ctrl)
	m.gd = mocks.NewMockGetDocumenter(m.ctrl)
	m.gdi = mocks.NewMockGetDocumentInfoer(m.ctrl)

	server.GetDocumentInfoer = m.gdi
	server.GetDocumenter = m.gd
	server.PutDocumenter = m.pd

	ctx := context.Background()
	connCtx, _ := context.WithTimeout(ctx, time.Second*5)

	conn, err := grpc.DialContext(
		connCtx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("grpc.DialContext: %v", err)
	}

	m.c = NewClient(pb.NewDocumentsClient(conn))

	return func() {
	}
}

func TestGetDocumentInfo(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	ttt := time.Unix(100, 10)

	repoRequest := &documents.GetDocumentInfoerRequest{DocumentId: "123"}
	repoResponse := &documents.GetDocumentInfoerResponse{Name: "laza", Size: 1234, UploadDate: ttt}

	m.gdi.EXPECT().Execute(gomock.Eq(repoRequest)).Return(repoResponse, nil)

	docInfo, err := m.c.GetDocumentInfo(context.Background(), "123")
	fmt.Println(err)

	assert.Nil(t, err)
	assert.NotNil(t, docInfo)
	assert.Equal(t, "laza", docInfo.Name)
	assert.Equal(t, int64(1234), docInfo.Size)
	assert.Equal(t, ttt.UTC(), docInfo.UploadDate)
}

func TestGetDocument(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

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

func TestPutDocument(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

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
