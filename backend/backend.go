package backend

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/api/test", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
