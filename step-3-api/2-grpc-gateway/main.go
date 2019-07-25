package main

import (
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/contracts"
	"net"
	"net/http"
	"os"
	"sync"
)

var mtx = sync.Mutex{}
var cache = map[int]int{}

func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	res, ok := cache[n]
	if ok {
		return res
	}

	res = fibonacci(n-1) + fibonacci(n-2)

	mtx.Lock()
	cache[n] = res
	mtx.Unlock()

	return cache[n]
}

func runClient() error {
	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	c := contracts.NewMathClient(conn)
	ctx := context.Background()
	req := &contracts.SumRequest{
		A: 5,
		B: 10,
	}
	res, err := c.Sum(ctx, req)
	if err != nil {
		return fmt.Errorf("did not calculate sum: %v", err)
	}
	fmt.Printf("Sum of 5 and 10 is: %d\n", res.GetResult())
	return nil
}

type mathServer struct{}

func (mathServer) Sum(ctx context.Context, req *contracts.SumRequest) (*contracts.Result, error) {
	res := contracts.Result{}
	res.Result = req.GetA() + req.GetB()
	return &res, nil
}

func (mathServer) Fibonacci(ctx context.Context, req *contracts.FactorialRequest) (*contracts.Result, error) {
	res := contracts.Result{}
	res.Result = int32(fibonacci(int(req.GetN())))
	return &res, nil
}

func runGrpcServer() error {
	println("serving grpc at :3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	contracts.RegisterMathServer(s, &mathServer{})

	if err := s.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve grpc: %v", err)
	}

	return nil
}

func runHttpServer() error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := contracts.RegisterMathHandlerFromEndpoint(context.Background(), mux, ":3000", opts)
	if err != nil {
		return fmt.Errorf("failed to register endpoint: %v", err)
	}

	r := http.NewServeMux()
	r.Handle("/", mux)

	println("serving http at :3001")
	if err := http.ListenAndServe(":3001", r); err != nil {
		return fmt.Errorf("failed to serve http: %v", err)
	}

	return nil
}

func main() {
	err := func() error {
		args := os.Args[1:]

		if len(args) == 0 {
			return fmt.Errorf("missing cli argument: can be 'client' or 'server'")
		}

		switch args[0] {
		case "client":
			return runClient()
		case "grpc-server":
			return runGrpcServer()
		case "http-server":
			return runHttpServer()
		case "server":
			errChan := make(chan error)
			go func() {
				errChan <- runHttpServer()
			}()
			go func() {
				errChan <- runGrpcServer()
			}()
			return <-errChan
		}

		return fmt.Errorf("unknown cli argument: can be 'client' or 'server'")
	}()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
