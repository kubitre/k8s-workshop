# Step 3: Golang

1. Prepare golang environment:
    - Locally (https://golang.org/doc/install);
    - Using docker: 
        - `mkdir -p ~/go/.cache`;
        - `export GOPATH=~/go`;
        - `alias gontainer="docker run -it --net=host -u $(id -u):$(id -g) -e XDG_CACHE_HOME=/tmp/cache -v $GOPATH/.cache:/tmp/cache -v ~/go:/go -w /go/$(realpath --relative-to=$GOPATH $PWD) golang:1.12.7"`;
        - `alias go="gontainer go`;
        - `go version`;
2. Lets create a new hello-world project:
    - `mkdir -p ~/go/src/math && cd ~/go/src/math`;
    - Create a `main.go` with content: 
    ```go
    package main
    import "fmt"
    func main() {
       fmt.Printf("hello, world\n")
    } 
    ```
    - Run `go fmt .`
    - Run `go run main.go`
3. Declare our "math" functions:
    ```go
    package main
    
    import "fmt"
    
    func sum(a int, b int) int {
    	return a + b
    }
    
    func fibonacci(n int) int {
    	if n < 2 {
    		return n
    	}
    	return fibonacci(n-1) + fibonacci(n-2)
    }
    
    func main() {
    	fmt.Printf("sum 1 + 2: %d\n", sum(1, 2))
    	fmt.Printf("fibonacci 20: %d\n", fibonacci(20))
    }

    ```
4. Then we can implement tty to add some concurrency:
    ```go
    reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter count to get its fibonacci or q to quit\n")
	for {
		s, _ := reader.ReadString('\n')
		s = strings.Trim(s, "\n")
		if s == "" {
			continue
		}
		if s == "q" {
			return
		}
		count, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			fmt.Print("Entered count is not an integer\n")
			continue
		}
		go func(n int) {
			startedAt := time.Now()
			result := fibonacci(n)
			fmt.Printf("Fibonacci for %d is %d (took %d ms)\n", n, result, time.Now().Sub(startedAt) * time.Millisecond)
		}(int(count))
	}
    ```
5. But calculating fibonacci every time is too expensive for resources. Lets add some cache:
   -  ```go
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
      ```
   - Note that we used sync.Mutex to avoid simultaneous access to the map.
6. Finally, lets implement http server:
    - ```go
      http.HandleFunc("/fibonacci", func(w http.ResponseWriter, r *http.Request) {
    		w.Header().Set("Content-Type", "application/json")
    		res := map[string]interface{}{}
    
    		err := func() error {
    			query := r.URL.Query()
    
    			rawN, ok := query["n"]
    			if !ok || len(rawN[0]) < 1 {
    				return fmt.Errorf("missing param: %s", "n")
    			}
    
    			n, err := strconv.ParseInt(rawN[0], 10, 0)
    			if err != nil {
    				return fmt.Errorf("param %s is not an integer", "n")
    			}
    
    			res["result"] = fibonacci(int(n))
    
    			return nil
    		}()
    		if err != nil {
    			w.WriteHeader(400)
    			res["error"] = err.Error()
    		}
    		jsonRes, _ := json.Marshal(res)
    		fmt.Fprintf(w, "%s", jsonRes)
    	})
    
    	if err := http.ListenAndServe(":3000", nil); err != nil {
    		println(err)
    	}
      ```
    - Note that we used anonymous function call to implement convenient error handling: 
    ```go
    err := func() error {
        // ...
    }()
    ```
// @todo: use `-v $PWD:$PWD -w $PWD` in go docker alias.