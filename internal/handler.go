package internal

import (
	"net/http"
)

type Handler struct {
	Queue Queue
}

func NewHandler(maxQueuesCount, maxMessagesInOneQueueCount, defaultTimeout int) http.Handler {
	h := &Handler{Queue: Queue{
		MaxQueuesCount:             maxQueuesCount,
		MaxMessagesInOneQueueCount: maxMessagesInOneQueueCount,
		DefaultTimeout:             defaultTimeout,
		Store:                      make(map[string]*[]Message),
	}}

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && queue.MatchString(r.URL.Path):
		h.getMessageFromQueue(w, r)
		return
	case r.Method == http.MethodPut && queue.MatchString(r.URL.Path):
		h.putMessageToQueue(w, r)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found handler"))
		return
	}
}
