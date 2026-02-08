package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/vibesql/vibe/internal/query"
)

const (
	DefaultHost     = "127.0.0.1"
	DefaultPort     = 5173
	MaxConnections  = 2
	ReadTimeout     = 10 * time.Second
	WriteTimeout    = 10 * time.Second
	ShutdownTimeout = 30 * time.Second
	IdleTimeout     = 30 * time.Second
	ReadHeaderTimeout = 5 * time.Second
)

// GetBindHost returns the host to bind to.
// Set VIBE_BIND_HOST=0.0.0.0 to allow LAN access.
func GetBindHost() string {
	if host := os.Getenv("VIBE_BIND_HOST"); host != "" {
		return host
	}
	return DefaultHost
}

type Server struct {
	host       string
	port       int
	httpServer *http.Server
	listener   net.Listener
	handler    *Handler
	ready      atomic.Bool
}

func NewServer(executor query.QueryExecutor) *Server {
	handler := NewHandler(executor)

	server := &Server{
		host:    GetBindHost(),
		port:    DefaultPort,
		handler: handler,
	}
	server.ready.Store(false)
	return server
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	s.handler.RegisterRoutes(mux)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind to %s: %w", addr, err)
	}
	s.listener = listener

	limitListener := &limitedListener{
		Listener:       listener,
		maxConnections: MaxConnections,
		semaphore:      make(chan struct{}, MaxConnections),
	}

	s.httpServer = &http.Server{
		Handler:           mux,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
		IdleTimeout:       IdleTimeout,
		ReadHeaderTimeout: ReadHeaderTimeout,
	}

	s.ready.Store(true)
	log.Printf("[INFO] HTTP server listening on %s (max connections: %d)", addr, MaxConnections)

	go func() {
		if err := s.httpServer.Serve(limitListener); err != nil && err != http.ErrServerClosed {
			log.Printf("[ERROR] HTTP server error: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	if s.httpServer == nil {
		return nil
	}

	log.Printf("[INFO] Shutting down HTTP server gracefully...")
	s.ready.Store(false)

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] HTTP server shutdown error: %v", err)
		return err
	}

	log.Printf("[INFO] HTTP server stopped")
	return nil
}

func (s *Server) IsReady() bool {
	return s.ready.Load()
}

func (s *Server) Addr() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

func (s *Server) WaitForShutdown() {
	if !s.IsReady() {
		log.Printf("[WARN] WaitForShutdown called but server not started")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigChan
	log.Printf("[INFO] Received signal: %v", sig)

	if err := s.Stop(); err != nil {
		log.Printf("[ERROR] Failed to stop server: %v", err)
	}
}

type limitedListener struct {
	net.Listener
	maxConnections int
	semaphore      chan struct{}
}

func (l *limitedListener) Accept() (net.Conn, error) {
	l.semaphore <- struct{}{}

	conn, err := l.Listener.Accept()
	if err != nil {
		<-l.semaphore
		return nil, err
	}

	return &limitedConn{
		Conn:      conn,
		semaphore: l.semaphore,
	}, nil
}

type limitedConn struct {
	net.Conn
	semaphore chan struct{}
	released  bool
}

func (c *limitedConn) Close() error {
	if c.released {
		return nil
	}
	<-c.semaphore
	c.released = true
	return c.Conn.Close()
}
