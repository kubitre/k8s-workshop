package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func main() {
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
}
