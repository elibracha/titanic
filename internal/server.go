package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"titanic-api/internal/config"
	"titanic-api/internal/healthcheck"
	"titanic-api/internal/passenger"
	"titanic-api/internal/web"
)

type Server interface {
	Start()
}

type server struct {
	conf *config.Config
}

func (s *server) Start() {
	router, err := s.router()
	if err != nil {
		log.Fatalf("setup error: %s", err)
	}

	p := fmt.Sprintf(":%d", s.conf.GetPort())
	srv := &http.Server{
		Addr:    p,
		Handler: router,
	}

	// start server
	go func() {
		log.Printf("Starting server on %s...", p)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server setup error: %s", err)
		}
	}()

	// listen for signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)

	// block until a signal received
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shut down gracefully
	log.Println("shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %s", err)
	}

	log.Println("server gracefully stopped.")
}

func (s *server) router() (*chi.Mux, error) {
	service, err := s.initService()
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()

	// setup middlewares
	router.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(time.Second*60),
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET"},
			MaxAge:         300,
		}),
	)

	// setup ui routes
	router.Route("/ui", func(r chi.Router) {
		r.Mount("/", web.NewHandler(service).RegisterHandler())
	})

	// setup static docs route
	fs := http.FileServer(http.Dir("docs"))
	router.Mount("/", http.StripPrefix("/", fs))

	// setup docs routes
	router.Get("/api/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/openapi.json"),
	))

	// setup api routes
	router.Route("/api/v1", func(r chi.Router) {
		// setup passenger routes
		r.Mount("/passenger", passenger.NewHandler(service).RegisterHandler())
		// setup health check routes
		r.Mount("/health", healthcheck.NewHandler().RegisterHandler())
	})

	return router, nil
}

func (s *server) initService() (passenger.Service, error) {
	var store passenger.Store
	storeType := s.conf.GetStoreType()
	switch storeType {
	case passenger.StoreTypeCSV:
		store = passenger.NewStoreCSV(s.conf.GetStorePath())
	case passenger.StoreTypeSQLite:
		store = passenger.NewStoreSQLite(passenger.NewConnector(s.conf.GetStorePath()))
	default:
		return nil, fmt.Errorf("store type provided not supported")
	}

	return passenger.NewService(store), nil

}

func NewServer() Server {
	return &server{conf: config.NewConfig()}
}
