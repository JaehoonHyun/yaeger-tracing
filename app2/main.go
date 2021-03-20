package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"app2/pkg/helloworld"
	"app2/pkg/tracer"

	"github.com/opentracing/opentracing-go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
)

//embeding
type server struct {
	helloworld.UnimplementedGreeterServer
	ip string
}

// overloading
func (s *server) OnMatching2(ctx context.Context, in *helloworld.Empty) (*helloworld.HelloReply, error) {
	log.Printf("Come in OnMatching2")

	//Dial to matching
	// s.ip = s.routeMatchingService("app3:20053")
	s.ip = s.routeMatchingService(ctx, endpointApp3)
	// s.ip = "localhost"

	return &helloworld.HelloReply{Message: s.ip}, nil

}

func (s *server) routeMatchingService(ctx context.Context, grpcEndpoint string) string {

	// conn, err := grpc.Dial(grpcEndpoint, grpc.WithInsecure())

	span, ctx := opentracing.StartSpanFromContext(ctx, "app2-routeMatchingService")
	defer span.Finish()

	span.LogKV("app2", "로그다잇")
	span.SetTag("tag", "tag다잇")

	tracer := opentracing.GlobalTracer()

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

	//ctx를 app3로 전달한다. 그럼. span도 전달될듯?
	msg, err := client.OnMatching3(ctx, &helloworld.Empty{})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Greeting: %s", msg.GetMessage())

	return msg.GetMessage()
}

var (
	port         int
	endpointApp1 string
	endpointApp2 string
	endpointApp3 string
)

func main() {
	flag.IntVar(&port, "port", 20052, "put the port")
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
