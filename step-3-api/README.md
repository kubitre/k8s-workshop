# Step 3: API

1. Lets specify contracts for our api; Interface will be the same:
    ```
    int function sum(int a, int b) {...}
    
    int function factorial(int n) {...}
    ```  
    - Using same path: `~/go/src/math`
    - Add new directory: `mkdir ~/go/src/math/contracts`
    - Create a .proto file for math service: `math.proto`
2. Define & generate service and data structures:
    - ```proto
      syntax = "proto3";
    
      package math;
    
      option go_package = "contracts";
    
      service math {
    
          rpc Sum (SumRequest) returns (Result) {}
    
          rpc Fibonacci(FibonacciRequest) returns (Result) {}
    
      }
    
      message SumRequest {
          int32 a = 1;
          int32 b = 2;
      }
    
      message FibonacciRequest {
          int32 n = 1;
      }
    
      message Result {
          int32 result = 1;
      }

      ```
    - Run `docker run --rm -u $(id -u):$(id -g) -v $PWD/contracts:/contracts -w /contracts thethingsindustries/protoc "--go_out=plugins=grpc:. -I. ./*.proto"`
    - Review generated `math.pb.go`;
    
3. Then we can use it in our generated http server:
    - For grpc client:
      ```go
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
      ```
    - For grpc server:
      ```go
      type mathServer struct {}
      
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
      ```
      
4. Lets add grpc gateway to our server to be accessible via http1:
    - To `math.proto`:
      ```proto
      rpc Sum (SumRequest) returns (Result) {
          option (google.api.http).get = "/v1alpha/math/sum";
      }
      
      rpc Fibonacci (FactorialRequest) returns (Result) {
          option (google.api.http).get = "/v1alpha/math/fibonacci";
      }
      ```
    - Run `docker run --rm -u $(id -u):$(id -g) -v $PWD:/contracts -w /contracts thethingsindustries/protoc --swagger_out=logtostderr=true:. --grpc-gateway_out=logtostderr=true:. -I. ./*.proto`
    - So we got some code for gateway, lets include it to our server:
      ```go
      	mux := runtime.NewServeMux()
      	opts := []grpc.DialOption{grpc.WithInsecure()}
      	err := contracts.RegisterMathHandlerFromEndpoint(context.Background(), mux, ":3000", opts)
      	if err != nil {
      		return fmt.Errorf("failed to register endpoint: %v", err)
      	}
      
      	r := http.NewServeMux()
      	r.Handle("/", mux)
      
      	if err := http.ListenAndServe(":3001", r); err != nil {
      		return fmt.Errorf("failed to serve http: %v", err)
      	}
      
      	return nil
      ``` 
    - Then we can run both (its important) grpc and http servers and try `curl 127.0.0.1:3000/v1alpha/math/fibonacci?n=10`
    - So we can keep contracts (i.e. in separate repo) and generate clients and servers for every language to be sure that our interfaces are compatible.
    