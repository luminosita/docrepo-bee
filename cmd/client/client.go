package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	pb "github.com/luminosita/docrepo-bee/api/documents/v1"
	server2 "github.com/luminosita/docrepo-bee/internal/server"
	grpc2 "github.com/luminosita/docrepo-bee/pkg/client/grpc"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"runtime"
)

var addr string
var version int

func setupClient() *grpc2.Client {
	claims := jwtv4.MapClaims{
		"sub": "laza",
		"id":  "123",
	}

	clientMid := []middleware.Middleware{
		jwt.Client(func(token *jwtv4.Token) (interface{}, error) {
			return []byte(server2.SecretKey), nil
		}, jwt.WithSigningMethod(jwtv4.SigningMethodHS256),
			jwt.WithClaims(func() jwtv4.Claims { return claims })),
	}

	conn, err := kgrpc.DialInsecure(
		context.Background(),
		kgrpc.WithEndpoint(addr),
		kgrpc.WithMiddleware(
			clientMid...,
		),
	)
	if err != nil {
		log.Fatalf("grpc.DialContext: %v", err)
	}

	return grpc2.NewClient(pb.NewDocumentsClient(conn))
}

func uploadDocument(ctx context.Context, c *grpc2.Client, filepath string) string {
	file, err := os.Open(filepath)
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

	log.Printf("File successfully uploaded: %s (%s)\n", filepath, docId)

	return docId
}

func downloadDocument(ctx context.Context, c *grpc2.Client, docId string,
	actDocInfo *grpc2.DocumentInfo, filepath string) {
	docInfo, r, err := c.GetDocument(ctx, docId)
	if err != nil {
		log.Fatalf("GetDocument: %v", err)
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("File open: %v", err)
	}

	actual, err := io.ReadAll(r)
	expected, err := io.ReadAll(file)

	if bytes.Compare(actual, expected) != 0 {
		log.Fatalf("Downloaded file differs from reference file (%s)", filepath)
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

	if actDocInfo.Name != docInfo.Name || actDocInfo.Size != docInfo.Size {
		log.Fatalf("Downloaded documentInfo differs from "+
			"reference documentInfo.\n Expected: %+v\n Actual: %+v", docInfo, actDocInfo)
	}

	log.Printf("Downloaded documentInfo and reference "+
		"documentInfo are identical: (%+v)", actDocInfo)
}

func commandRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "client",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(2)
		},
	}
	rootCmd.AddCommand(commandServe())
	rootCmd.AddCommand(commandVersion())
	return rootCmd
}

func commandServe() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "all filepath",
		Short:   "Launch all tasks",
		Example: "all filepath",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			filepath := args[0]

			return allTasks(filepath)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&addr, "address", "localhost:9000", "gRPC server address")
	flags.IntVar(&version, "version", 1, "Version of client")

	return cmd
}

func commandVersion() *cobra.Command {
	version := "DEV"

	return &cobra.Command{
		Use:   "version",
		Short: "Print the version and exit",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf(
				"Bee Version: %s\nGo Version: %s\nGo OS/ARCH: %s %s\n",
				version,
				runtime.Version(),
				runtime.GOOS,
				runtime.GOARCH,
			)
		},
	}
}

func allTasks(filepath string) error {
	ctx := context.Background()

	c := setupClient()

	docId := uploadDocument(ctx, c, filepath)

	actDocInfo, err := c.GetDocumentInfo(ctx, docId)
	if err != nil {
		log.Fatalf("GetDocumentInfo: %v", err)
	}

	downloadDocument(ctx, c, docId, actDocInfo, filepath)

	return nil
}

func main() {
	if err := commandRoot().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
