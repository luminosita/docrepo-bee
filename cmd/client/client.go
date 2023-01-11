package main

import (
	"bytes"
	"context"
	"flag"
	pb "github.com/luminosita/docrepo-bee/api/gen/v1"
	grpc2 "github.com/luminosita/docrepo-bee/pkg/client/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"time"
)

var addr string
var filePath string
var version int

func init() {
	flag.StringVar(&addr, "address", "localhost:9080", "filestorage address")
	flag.StringVar(&filePath, "file-path", "", "Path to get a file")
	flag.IntVar(&version, "version", 1, "Version of client")
}

func main() {
	flag.Parse()

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

	c := grpc2.NewClient(pb.NewDocumentsClient(conn))

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("File open: %v", err)
	}

	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("File stat: %v", err)
	}

	docInfo := &grpc2.DocumentInfo{
		Name: stat.Name(),
		Size: stat.Size(),
	}

	docId, err := c.PutDocument(ctx, docInfo, file)
	if err != nil {
		log.Fatalf("PutDocument: %v", err)
	}

	log.Printf("File successfully uploaded: %s (%s)\n", filePath, docId)

	docInfo, r, err := c.GetDocument(ctx, docId)
	if err != nil {
		log.Fatalf("GetDocument: %v", err)
	}

	file, err = os.Open(filePath)
	if err != nil {
		log.Fatalf("File open: %v", err)
	}

	actual, err := io.ReadAll(r)
	expected, err := io.ReadAll(file)

	if bytes.Compare(actual, expected) != 0 {
		log.Fatalf("Downloaded file differs from reference file (%s)", filePath)
	}

	log.Println("Downloaded file and reference file are identical.")

	err = file.Close()
	if err != nil {
		log.Fatalf("File close: %v", err)
	}

	err = r.Close()
	if err != nil {
		log.Fatalf("Download stream close: %v", err)
	}

	actDocInfo, err := c.GetDocumentInfo(ctx, docId)
	if err != nil {
		log.Fatalf("GetDocumentInfo: %v", err)
	}

	if actDocInfo.Name != docInfo.Name || actDocInfo.Size != docInfo.Size {
		log.Fatalf("Downloaded documentInfo differs from "+
			"reference documentInfo.\n Expected: %+v\n Actual: %+v", docInfo, actDocInfo)
	}

	log.Printf("Downloaded documentInfo and reference "+
		"documentInfo are identical: (%+v)", actDocInfo)
}
