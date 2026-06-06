package controller

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/golangmonster/workspace-booking-service/internal/config"
	interceptor "github.com/golangmonster/workspace-booking-service/internal/controller/middleware/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type ImplementationAdapter interface {
	RegisterHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
	RegisterServer(grpcServer *grpc.Server)
}

type controller struct {
	cfg             *config.Config
	implementations []ImplementationAdapter
	grpcServer      *grpc.Server
	httpServer      *http.Server
}

func New(cfg *config.Config, implementations ...ImplementationAdapter) *controller {
	return &controller{
		cfg:             cfg,
		implementations: implementations,
	}
}

func (c *controller) Run(ctx context.Context) {
	c.ServeGRPC()
	c.ServeHTTP(ctx)
}

func (c *controller) Stop(ctx context.Context) {
	log.Info("shutting down servers...")

	if c.grpcServer != nil {
		stopped := make(chan struct{})
		go func() {
			c.grpcServer.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
		case <-ctx.Done():
			c.grpcServer.Stop()
			log.Warn("grpc server force stopped due to context deadline")
		}
	}

	if c.httpServer != nil {
		if err := c.httpServer.Shutdown(ctx); err != nil {
			log.Errorf("http server shutdown error: %v", err)
		} else {
			log.Info("http server stopped gracefully")
		}
	}
}

func (c *controller) ServeGRPC() {
	lis, err := net.Listen("tcp", c.cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("%v: failed to listen grpc port: %s", c.cfg.GRPCAddr, err)
	}

	c.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.WithValidation(),
		),
	)

	for _, imp := range c.implementations {
		imp.RegisterServer(c.grpcServer)
	}

	go func() {
		reflection.Register(c.grpcServer)
		log.Infof("grpc addr: %s", lis.Addr())
		if err = c.grpcServer.Serve(lis); err != nil {
			log.Fatal(err, "failed to serve grpc")
		}
	}()
}

func (c *controller) ServeHTTP(ctx context.Context) {
	runtimeMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{UseProtoNames: true, EmitUnpopulated: true},
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	var err error
	for _, imp := range c.implementations {
		err = multierr.Append(err, imp.RegisterHandlerFromEndpoint(ctx, runtimeMux, c.cfg.GRPCAddr, opts))
	}

	if err != nil {
		log.Fatal(err, "failed to register gateway")
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", runtimeMux)
	prefix := "/docs/"

	fs := http.FileServer(http.Dir("./swagger/"))
	httpMux.Handle(prefix, http.StripPrefix(prefix, fs))

	c.httpServer = &http.Server{
		Addr:    c.cfg.HTTPAddr,
		Handler: cors.AllowAll().Handler(httpMux),
	}

	go func() {
		log.Infof("http addr: %s", c.cfg.HTTPAddr)
		log.Infof("swagger addr: %s/docs", c.cfg.HTTPAddr)
		if err = c.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err, "failed to serve http")
		}
	}()
}
