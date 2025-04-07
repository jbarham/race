// NB: Tests in this file must be run with the -race flag
//go:build race

package main

import (
	"net/http/httptest"
	"sync"
	"testing"
)

func testHandler(t *testing.T, h *countHandler) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
		}()
	}
	wg.Wait()
}

func TestSafeHandler(t *testing.T) {
	testHandler(t, newCountHandler(false))
}

func TestUnsafeHandler(t *testing.T) {
	testHandler(t, newCountHandler(true))
}
