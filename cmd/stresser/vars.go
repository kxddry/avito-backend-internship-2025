package main

import (
	"flag"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

//nolint:gochecknoglobals
var (
	baseURL       = flag.String("url", "http://localhost:8080", "Base URL of the service")
	duration      = flag.Duration("duration", 60*time.Second, "Test duration")
	workers       = flag.Int("workers", runtime.NumCPU()*2, "Number of worker goroutines")
	requestsTotal atomic.Uint64
	requestsOK    atomic.Uint64
	requestsFail  atomic.Uint64
)

var (
	teamNames  = []string{"backend", "frontend", "mobile", "devops", "qa"}
	userIDs    = make([]string, 0, 200)
	prIDs      = make([]string, 0, 1000)
	prIDsMutex sync.RWMutex
)
