package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"app1/pkg/helloworld"
	"app1/pkg/tracer"

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
func (s *server) OnMatching(ctx context.Context, in *helloworld.Empty) (*helloworld.HelloReply, error) {
	log.Printf("Come in OnMatching")

	//Dial to matching
	// s.ip = s.routeMatchingService("app2:20052")
	s.ip = s.routeMatchingService(ctx, endpointApp2)
	// s.ip = "localhost"

	//TODO: tracing ctx를 가져와서 span을 추가할 수 있나?

	return &helloworld.HelloReply{Message: s.ip}, nil

}

func (s *server) routeMatchingService(ctx context.Context, grpcEndpoint string) string {

	// conn, err := grpc.Dial(grpcEndpoint, )

	tracer := opentracing.GlobalTracer()
	span, ctx := opentracing.StartSpanFromContext(ctx, "app1-routeMatchingService")
	defer span.Finish()

	span.LogKV("app1", "log")
	span.SetTag("app1", "tag")

	// initialize client with tracing interceptor using grpc client side chaining
	conn, err := grpc.Dial(
		grpcEndpoint,
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(tracer)),
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer)),
		)),
	)

	if err != nil {
		log.Fatalf("did not connected : %v", err)
	}
	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	msg, err := client.OnMatching2(ctx, &helloworld.Empty{})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Greeting: %s", msg.GetMessage())

	return msg.GetMessage()
}

var (
	port         int
	endpointApp2 string
	endpointApp3 string
)

func main() {
	flag.IntVar(&port, "port", 20051, "put the port")
	flag.StringVar(&endpointApp2, "endpointApp2", "localhost:20052", "endpoint of app2")
	flag.StringVar(&endpointApp3, "endpointApp3", "localhost:20053", "endpoint of app3")

	flag.Parse()

	log.Println("server has started")
	log.Println("port -> ", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	//tracing
	//initialize tracer
	//한 트랜잭션을 위한 tracer를 생성해낸다.
	tracer, closer, err := tracer.NewTracer()
	defer closer.Close()
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	opentracing.SetGlobalTracer(tracer)

	//with tracing
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
