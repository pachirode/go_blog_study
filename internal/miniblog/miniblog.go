package miniblog

import (
	"encoding/json"
	"fmt"

	"github.com/marmotedu/Miniblog/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func NewMiniBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "miniblog",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
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

	return cmd
}

func run() error {
	setting, _ := json.Marshal(viper.AllSettings())
	log.Infow("Config settings: ", string(setting))
	return nil
}
