package handlers

import (
	"context"
	"github.com/google/wire"
	grpc2 "github.com/luminosita/common-bee/pkg/grpc"
	pb "github.com/luminosita/docrepo-bee/api/gen/v1"
	documents2 "github.com/luminosita/docrepo-bee/internal/interfaces/use-cases/documents"
	"github.com/luminosita/honeycomb/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var PbDocumentsWireSet = wire.NewSet(NewDocumentsServer,
	wire.Bind(new(pb.DocumentsServer), new(*DocumentsServer)))

type DocumentsServer struct {
	documents2.GetDocumentInfoer
	documents2.GetDocumenter
	documents2.PutDocumenter
}

func NewDocumentsServer(gdi documents2.GetDocumentInfoer,
	gd documents2.GetDocumenter, pd documents2.PutDocumenter) *DocumentsServer {
	return &DocumentsServer{
		GetDocumentInfoer: gdi,
		GetDocumenter:     gd,
		PutDocumenter:     pd,
	}
}

func (s *DocumentsServer) GetDocumentInfo(ctx context.Context,
	req *pb.GetDocumentInfoRequest) (*pb.GetDocumentInfoReply, error) {

	res, err := s.GetDocumentInfoer.Execute(&documents2.GetDocumentInfoerRequest{
		DocumentId: req.DocumentId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetDocumentInfoReply{
		Info: &pb.DocumentInfo{
			Name: res.Name,
			Size: res.Size,
		},
		UploadDate: timestamppb.New(res.UploadDate),
	}, nil
}

func (s *DocumentsServer) GetDocument(req *pb.GetDocumentRequest,
	srv pb.Documents_GetDocumentServer) error {
	res, err := s.GetDocumenter.Execute(&documents2.GetDocumenterRequest{
		DocumentId: req.DocumentId,
	})
	if err != nil {
		return err
	}

	info := &pb.GetDocumentReply_Info{
		Info: &pb.DocumentInfo{
			Name: res.Name,
			Size: res.Size,
		},
	}

	serverErr := srv.Send(&pb.GetDocumentReply{Data: info})
	if serverErr != nil {
		return log.LogError(status.Errorf(codes.Internal,
			"unable to send document info: %v", serverErr))
	}

	grpc2.CopyToServerStream(res.Reader, srv, func(buffer []byte) any {
		return &pb.PutDocumentRequest{
			Data: &pb.PutDocumentRequest_ChunkData{
				ChunkData: buffer,
			},
		}
	})

	return nil
}

func (s *DocumentsServer) PutDocument(srv pb.Documents_PutDocumentServer) error {
	req, err := srv.Recv()
	if err != nil {
		return log.LogError(status.Errorf(codes.Unknown, "cannot receive document info"))
	}

	info := req.GetInfo()

	res, err := s.PutDocumenter.Execute(&documents2.PutDocumenterRequest{
		Name: info.Name,
		Size: info.Size,
	})

	reply := new(pb.GetDocumentReply)

	grpc2.CopyFromServerStream(res.Writer, srv, reply, func(reply any) []byte {
		return reply.(*pb.GetDocumentReply).GetData().(*pb.GetDocumentReply_ChunkData).ChunkData
	})

	err = res.Writer.Close()
	if err != nil {
		return log.LogError(status.Errorf(codes.Internal, "unable to close PutDocument writer"))
	}

	return srv.SendAndClose(&pb.PutDocumentReply{
		DocumentId: res.DocumentId,
	})
}
