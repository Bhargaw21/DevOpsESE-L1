// main.go
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "iot-sensor-api ok\n")
	})

	http.HandleFunc("/sensor", func(w http.ResponseWriter, r *http.Request) {
		// simulate reading some sensor values
		payload := map[string]interface{}{
			"id":        rand.Intn(1000),
			"temperature": fmt.Sprintf("%.2f", 20+rand.Float64()*15),
			"humidity":    fmt.Sprintf("%.2f", 20+rand.Float64()*60),
			"ts":          time.Now().UTC().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	})

	// Busy CPU burn endpoint to raise CPU usage for testing HPA
	http.HandleFunc("/burn", func(w http.ResponseWriter, r *http.Request) {
		// ?ms=200  => perform CPU-bound work for ~200ms
		msParam := r.URL.Query().Get("ms")
		ms := 200
		if msParam != "" {
			if v, err := strconv.Atoi(msParam); err == nil {
				if v > 0 {
					ms = v
				}
			}
		}
		doCPUBurn(time.Duration(ms) * time.Millisecond)
		w.WriteHeader(200)
		fmt.Fprintf(w, "burned %d ms\n", ms)
	})

	addr := ":8080"
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	fmt.Printf("listening on %s\n", addr)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}

func doCPUBurn(d time.Duration) {
	// busy loop for ~d duration
	end := time.Now().Add(d)
	var x float64 = 0.0001
	for time.Now().Before(end) {
		// some floating-point ops to consume CPU
		for i := 0; i < 1000; i++ {
			x += x * 1.0000001
		}
	}
	_ = x
}
