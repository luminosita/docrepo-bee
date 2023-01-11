package grpc

import (
	"context"
	"errors"
	grpc2 "github.com/luminosita/common-bee/pkg/grpc"
	pb "github.com/luminosita/docrepo-bee/api/gen/v1"
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

	grpc2.CopyToClientStream(r, stream, func(buffer []byte) any {
		return &pb.PutDocumentRequest{
			Data: &pb.PutDocumentRequest_ChunkData{
				ChunkData: buffer,
			},
		}
	})

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
		return nil, nil, errors.New("Invalid or malformatted document info sent from server")
	}

	docInfo := &DocumentInfo{
		Name: rInfo.Name,
		Size: rInfo.Size,
	}

	r, w := io.Pipe()

	reply := new(pb.GetDocumentReply)

	go grpc2.CopyFromClientStream(w, res, reply, func(reply any) []byte {
		return reply.(*pb.GetDocumentReply).GetData().(*pb.GetDocumentReply_ChunkData).ChunkData
	})

	return docInfo, r, nil
}
