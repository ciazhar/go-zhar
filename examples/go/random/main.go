package main

import (
	"net/http"
	"sync"
)

type Handler struct {
	// ...
}

func (h *Handler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {

}

func main() {
	mu := new(sync.Mutex)

}
