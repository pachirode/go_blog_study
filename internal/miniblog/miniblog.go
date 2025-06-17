package miniblog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/Miniblog/internal/pkg/log"
	"github.com/marmotedu/Miniblog/pkg/version/verflag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	setting, _ := json.Marshal(viper.AllSettings())
	log.Debugw("Config settings: ", string(setting))

	gin.SetMode(viper.GetString("web.runmode"))
	g := gin.New()

	g.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "Page not found"})
	})

	g.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"code": 200})
	})

	httpServe := http.Server{Addr: viper.GetString("web.addr"), Handler: g}
	log.Infow("Starting listen requests on http address", "addr", viper.GetString("web.addr"))

	if err := httpServe.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw(err.Error())
	}
	return nil
}
