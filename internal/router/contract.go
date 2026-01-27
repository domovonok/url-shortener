package router

import "net/http"

type LinkHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

type TokenBucket interface {
	Allow() bool
	Capacity() int
	Remaining() int
}
