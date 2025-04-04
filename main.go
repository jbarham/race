package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"
)

var unsafe = flag.Bool("unsafe", false, "set to turn off locks")

type RWLocker interface {
	sync.Locker
	RLock()
	RUnlock()
}

type DummyLock struct{}

func (d DummyLock) Lock()    {}
func (d DummyLock) Unlock()  {}
func (d DummyLock) RLock()   {}
func (d DummyLock) RUnlock() {}

// Inspired by https://pkg.go.dev/net/http#example-Handle
// but here we use a map which must be locked for both writes and reads (including iteration).
type countHandler struct {
	rwlock RWLocker
	m      map[string]int
}

func newCountHandler(unsafe bool) *countHandler {
	var rwlock RWLocker
	if unsafe {
		rwlock = DummyLock{}
	} else {
		rwlock = &sync.RWMutex{}
	}
	return &countHandler{
		m:      make(map[string]int),
		rwlock: rwlock,
	}
}

func (s *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Lock for write
	s.rwlock.Lock()
	s.m["count"] = s.m["count"] + 1
	s.rwlock.Unlock()
	// Lock for read
	s.rwlock.RLock()
	json.NewEncoder(w).Encode(s.m)
	s.rwlock.RUnlock()
}

func main() {
	flag.Parse()
	http.Handle("/", newCountHandler(*unsafe))
	log.Fatal(http.ListenAndServe(":5000", nil))
}
