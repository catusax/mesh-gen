package template

// HandlerSRV is the handler template used for new service projects.
var HandlerSRV = `package handler

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"

	pb "{{.Vendor}}{{.Service}}/proto"
)

type {{title .Service}} struct{
    pb.Unimplemented{{title .Service}}Server
}

func (e *{{title .Service}}) Call(ctx context.Context, req *pb.CallRequest)(*pb.CallResponse,error) {
	zap.L().Info("Received {{title .Service}}.Call request: ", zap.Any("request", req))
	return &pb.CallResponse{
        Msg : "Hello " + req.Name,
    }, nil
}

func (e *{{title .Service}}) ClientStream(stream pb.{{title .Service}}_ClientStreamServer) error {
	var count int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			zap.L().Info("Got pings", zap.Int64("total", count))
			return stream.SendMsg(&pb.ClientStreamResponse{Count: count})
		}
		if err != nil {
			return err
		}
		zap.L().Info(fmt.Sprintf("Got ping %v", req.Stroke))
		count++
	}
}

func (e *{{title .Service}}) ServerStream(req *pb.ServerStreamRequest, stream pb.{{title .Service}}_ServerStreamServer) error {
	zap.L().Info("Received {{title .Service}}.ServerStream request: ", zap.Any("request", req))
	for i := 0; i < int(req.Count); i++ {
		zap.L().Info(fmt.Sprintf("Sending %d", i))
		if err := stream.Send(&pb.ServerStreamResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 250)
	}
	return nil
}

func (e *{{title .Service}}) BidiStream(stream pb.{{title .Service}}_BidiStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		zap.L().Info(fmt.Sprintf("Got ping %v", req.Stroke))
		if err := stream.Send(&pb.BidiStreamResponse{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
`
