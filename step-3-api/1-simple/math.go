package __simple

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/contracts"
	"net"
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

func runServer() error {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	contracts.RegisterMathServer(s, &mathServer{})
	if err := s.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
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
		case "server":
			return runServer()
		}

		return fmt.Errorf("unknown cli argument: can be 'client' or 'server'")
	}()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
