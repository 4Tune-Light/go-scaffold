package greeter

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ServiceName = "greeter.GreeterService"
	GreetMethod = "/greeter.GreeterService/Greet"
)

type GreetRequest struct {
	Name string
}

type GreetResponse struct {
	Message string
}

type GRPCServer struct {
	svc Service
}

func NewGRPCServer(svc Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) Greet(ctx context.Context, req *GreetRequest) (*GreetResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	resp, err := s.svc.Greet(ctx, name)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("greet failed: %v", err))
	}

	return &GreetResponse{Message: resp.Message}, nil
}

func RegisterGRPC(s *grpc.Server, svc Service) {
	g := NewGRPCServer(svc)
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: ServiceName,
		HandlerType: (*GRPCServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Greet",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					req := &GreetRequest{}
					if err := dec(req); err != nil {
						return nil, err
					}
					if interceptor == nil {
						return g.Greet(ctx, req)
					}
					info := &grpc.UnaryServerInfo{
						Server:     srv,
						FullMethod: GreetMethod,
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return g.Greet(ctx, req.(*GreetRequest))
					}
					return interceptor(ctx, req, info, handler)
				},
			},
		},
		Streams: []grpc.StreamDesc{},
	}, g)
}
