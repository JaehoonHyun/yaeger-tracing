package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"app3/pkg/helloworld"
	"app3/pkg/tracer"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//embeding
type server struct {
	helloworld.UnimplementedGreeterServer
	ip string
}

// overloading
func (s *server) OnMatching3(ctx context.Context, in *helloworld.Empty) (*helloworld.HelloReply, error) {
	log.Printf("Come in OnMatching3")

	//Dial to matching
	// s.ip = s.routeMatchingService("app3:20053")
	// s.ip = s.routeMatchingService(endpointApp3)
	s.ip = "localhost"

	return &helloworld.HelloReply{Message: s.ip}, nil

}

var (
	port         int
	endpointApp1 string
	endpointApp2 string
	endpointApp3 string
)

func main() {
	flag.IntVar(&port, "port", 20053, "put the port")
	flag.StringVar(&endpointApp1, "endpointApp1", "localhost:20051", "endpoint of app1")
	flag.StringVar(&endpointApp2, "endpointApp2", "localhost:20052", "endpoint of app2")
	flag.StringVar(&endpointApp3, "endpointApp3", "localhost:20053", "endpoint of app3")

	flag.Parse()

	log.Println("server has started")
	log.Println("port -> ", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	tracer, closer, err := tracer.NewTracer()
	defer closer.Close()
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	opentracing.SetGlobalTracer(tracer)

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			// add opentracing stream interceptor to chain
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// add opentracing unary interceptor to chain
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
		)),
	)
	srv := new(server)

	helloworld.RegisterGreeterServer(s, srv)
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
