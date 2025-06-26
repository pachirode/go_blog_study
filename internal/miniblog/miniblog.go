package miniblog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marmotedu/Miniblog/internal/pkg/core"
	"github.com/marmotedu/Miniblog/internal/pkg/errno"
	"github.com/marmotedu/Miniblog/internal/pkg/log"
	mv "github.com/marmotedu/Miniblog/internal/pkg/middleware"
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

	gin.SetMode(viper.GetString("web.runmode"))
	g := gin.New()

	mws := []gin.HandlerFunc{gin.Recovery(), mv.NoCache, mv.Cors, mv.RequestID(), mv.Secure}
	g.Use(mws...)

	g.NoRoute(func(ctx *gin.Context) {
		core.WriteResponse(ctx, errno.ErrPageNotFound, nil)
	})

	g.GET("/test", func(ctx *gin.Context) {
		log.C(ctx).Infow("/test function called")
		core.WriteResponse(ctx, nil, map[string]string{"status": "ok"})
	})

	httpServe := http.Server{Addr: viper.GetString("web.addr"), Handler: g}
	log.Infow("Starting listen requests on http address", "addr", viper.GetString("web.addr"))

	go func() {
		if err := httpServe.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infow("Starting stop server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServe.Shutdown(ctx); err != nil {
		log.Fatalw("Shutdown server err", err)
	}

	log.Infow("Serer stop succeed")

	return nil
}
