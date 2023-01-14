package grpc

import (
	"context"
	"fmt"
	grpc2 "github.com/luminosita/common-bee/pkg/grpc"
	"github.com/luminosita/common-bee/pkg/log"
	pb "github.com/luminosita/docrepo-bee/api/documents/v1"
	"io"
	"time"
)

type DocumentInfo struct {
	Name       string
	Size       int64
	UploadDate time.Time
}

type Client struct {
	dc pb.DocumentsClient
}

func NewClient(dc pb.DocumentsClient) *Client {
	return &Client{dc: dc}
}

func (c *Client) GetDocumentInfo(ctx context.Context, id string) (*DocumentInfo, error) {
	res, err := c.dc.GetDocumentInfo(ctx, &pb.GetDocumentInfoRequest{DocumentId: id})
	if err != nil {
		return nil, err
	}

	return &DocumentInfo{
		Name:       res.Info.Name,
		Size:       res.Info.Size,
		UploadDate: res.UploadDate.AsTime(),
	}, nil
}

func (c *Client) PutDocument(ctx context.Context, docInfo *DocumentInfo, r io.ReadCloser) (string, error) {
	stream, err := c.dc.PutDocument(ctx)
	if err != nil {
		return "", err
	}

	req := &pb.PutDocumentRequest{
		Data: &pb.PutDocumentRequest_Info{
			Info: &pb.DocumentInfo{
				Name: docInfo.Name,
				Size: docInfo.Size,
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		return "", err
	}

	err = grpc2.CopyToMessageStream(r, stream, func(buffer []byte) any {
		return &pb.PutDocumentRequest{
			Data: &pb.PutDocumentRequest_ChunkData{
				ChunkData: buffer,
			},
		}
	})
	if err != nil {
		return "", err
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", nil
	}

	return res.DocumentId, nil
}

func (c *Client) GetDocument(ctx context.Context, id string) (*DocumentInfo, io.ReadCloser, error) {
	res, err := c.dc.GetDocument(ctx, &pb.GetDocumentRequest{DocumentId: id})
	if err != nil {
		return nil, nil, err
	}

	recv, err := res.Recv()
	if err != nil {
		return nil, nil, err
	}

	rInfo := recv.GetInfo()
	if rInfo == nil {
		return nil, nil, log.LogErrorf("Invalid or malformatted document info sent from server")
	}

	docInfo := &DocumentInfo{
		Name: rInfo.Name,
		Size: rInfo.Size,
	}

	r, w := io.Pipe()

	reply := new(pb.GetDocumentReply)

	go func() {
		err := grpc2.CopyFromMessageStream(w, res, reply, func(reply any) []byte {
			return reply.(*pb.GetDocumentReply).GetData().(*pb.GetDocumentReply_ChunkData).ChunkData
		})
		if err != nil {
			fmt.Printf("unexpected error occuried: %+v", err)
			_ = w.CloseWithError(err)
		}
		_ = w.Close()
		err = res.CloseSend()
		if err != nil {
			fmt.Printf("unexpected error occuried: %+v", err)
		}
	}()

	return docInfo, r, nil
}
