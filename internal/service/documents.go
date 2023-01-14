package service

import (
	"context"
	"github.com/google/wire"
	grpc2 "github.com/luminosita/common-bee/pkg/grpc"
	"github.com/luminosita/common-bee/pkg/log"
	v1 "github.com/luminosita/docrepo-bee/api/documents/v1"
	documents2 "github.com/luminosita/docrepo-bee/internal/interface/app/documents"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var PbDocumentsWireSet = wire.NewSet(NewDocumentsServer,
	wire.Bind(new(v1.DocumentsServer), new(*Documents)))

type Documents struct {
	documents2.GetDocumentInfoer
	documents2.GetDocumenter
	documents2.PutDocumenter
}

func NewDocumentsServer(gdi documents2.GetDocumentInfoer,
	gd documents2.GetDocumenter, pd documents2.PutDocumenter) *Documents {
	return &Documents{
		GetDocumentInfoer: gdi,
		GetDocumenter:     gd,
		PutDocumenter:     pd,
	}
}

func (s *Documents) GetDocumentInfo(_ context.Context,
	req *v1.GetDocumentInfoRequest) (*v1.GetDocumentInfoReply, error) {

	res, err := s.GetDocumentInfoer.Execute(&documents2.GetDocumentInfoerRequest{
		DocumentId: req.DocumentId,
	})
	if err != nil {
		return nil, log.LogError(status.Errorf(codes.Unknown, "cannot execute service method: %s", err))
	}

	return &v1.GetDocumentInfoReply{
		Info: &v1.DocumentInfo{
			Name: res.Name,
			Size: res.Size,
		},
		UploadDate: timestamppb.New(res.UploadDate),
	}, nil
}

func (s *Documents) GetDocument(req *v1.GetDocumentRequest,
	srv v1.Documents_GetDocumentServer) error {
	res, err := s.GetDocumenter.Execute(&documents2.GetDocumenterRequest{
		DocumentId: req.DocumentId,
	})
	if err != nil {
		return log.LogError(status.Errorf(codes.Unknown, "cannot execute service method: %s", err))
	}

	info := &v1.GetDocumentReply_Info{
		Info: &v1.DocumentInfo{
			Name: res.Name,
			Size: res.Size,
		},
	}

	serverErr := srv.Send(&v1.GetDocumentReply{Data: info})
	if serverErr != nil {
		return log.LogError(status.Errorf(codes.Internal,
			"unable to send document info: %v", serverErr))
	}

	err = grpc2.CopyToMessageStream(res.Reader, srv, func(buffer []byte) any {
		return &v1.PutDocumentRequest{
			Data: &v1.PutDocumentRequest_ChunkData{
				ChunkData: buffer,
			},
		}
	})
	if err != nil {
		return log.LogError(status.Errorf(codes.Unknown, "cannot copy from message stream"))
	}

	return nil
}

func (s *Documents) PutDocument(srv v1.Documents_PutDocumentServer) error {
	req, err := srv.Recv()
	if err != nil {
		return log.LogError(status.Errorf(codes.Unknown, "cannot receive document info"))
	}

	info := req.GetInfo()

	res, err := s.PutDocumenter.Execute(&documents2.PutDocumenterRequest{
		Name: info.Name,
		Size: info.Size,
	})
	if err != nil {
		return log.LogError(status.Errorf(codes.Unknown, "cannot execute service method: %s", err))
	}

	reply := new(v1.GetDocumentReply)

	err = grpc2.CopyFromMessageStream(res.Writer, srv, reply, func(reply any) []byte {
		return reply.(*v1.GetDocumentReply).GetData().(*v1.GetDocumentReply_ChunkData).ChunkData
	})
	if err != nil {
		return log.LogError(status.Errorf(codes.Unknown, "cannot copy from message stream"))
	}

	err = res.Writer.Close()
	if err != nil {
		return log.LogError(status.Errorf(codes.Internal, "unable to close PutDocument writer"))
	}

	return srv.SendAndClose(&v1.PutDocumentReply{
		DocumentId: res.DocumentId,
	})
}
