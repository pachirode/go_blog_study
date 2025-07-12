package miniblog

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/marmotedu/Miniblog/internal/miniblog/controller/v1/user"
	"github.com/marmotedu/Miniblog/internal/miniblog/store"
	"github.com/marmotedu/Miniblog/internal/pkg/known"
	"github.com/marmotedu/Miniblog/internal/pkg/log"
	mv "github.com/marmotedu/Miniblog/internal/pkg/middleware"
	pb "github.com/marmotedu/Miniblog/pkg/proto/miniblog/v1"
	"github.com/marmotedu/Miniblog/pkg/token"
	"github.com/marmotedu/Miniblog/pkg/version/verflag"
)

var cfgFile string

func NewMiniBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "miniblog",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested()
			log.Init(logOptions())
			defer log.Sync()

			return run()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}
	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the miniblog config file")

	verflag.AddFlags(cmd.PersistentFlags())
	return cmd
}

func run() error {
	if err := initStore(); err != nil {
		return err
	}

	token.Init(viper.GetString("web.jwt-secret"), known.Usernamekey)

	gin.SetMode(viper.GetString("web.runmode"))
	g := gin.New()

	mws := []gin.HandlerFunc{gin.Recovery(), mv.NoCache, mv.Cors, mv.RequestID(), mv.Secure}
	g.Use(mws...)

	if err := installRouters(g); err != nil {
		return nil
	}

	httpServer := startInsecureServer(g)
	httpsServer := startSercureServer(g)
	grpcServer := *startGrpcServer()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infow("Starting stop server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalw("Shutdown http server err", err)
	}
	if err := httpsServer.Shutdown(ctx); err != nil {
		log.Fatalw("Shutdown https server err", err)
	}

	grpcServer.GracefulStop()

	log.Infow("Servers stop succeed")

	return nil
}

func startInsecureServer(g *gin.Engine) *http.Server {
	httpServer := &http.Server{Addr: viper.GetString("web.addr"), Handler: g}
	log.Infow("Starting listen requests on http address", "addr", viper.GetString("web.addr"))

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	return httpServer
}

func startSercureServer(g *gin.Engine) *http.Server {
	httpsServer := &http.Server{Addr: viper.GetString("tls.addr"), Handler: g}
	log.Infow("Starting listen requests on https address", "addr", viper.GetString("web.addr"))

	cert, key := viper.GetString("tls.cert"), viper.GetString("tls.key")
	if cert != "" && key != "" {
		go func() {
			if err := httpsServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalw(err.Error())
			}
		}()
	}

	return httpsServer
}

func startGrpcServer() *grpc.Server {
	lis, err := net.Listen("tcp", viper.GetString("grpc.addr"))
	if err != nil {
		log.Fatalw("Failed to listen", "err", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMiniBlogServer(grpcServer, user.New(store.S, nil))

	log.Infow("Starting listen requests on grpc address", "addr", viper.GetString("grpc.addr"))
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalw(err.Error())
		}
	}()

	return grpcServer
}
