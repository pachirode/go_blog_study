package miniblog

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/marmotedu/Miniblog/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	recommendedHomeDir = ".miniblog"
	defaultConfigName  = "miniblog.yaml"
)

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, recommendedHomeDir))

		viper.AddConfigPath(".")

		viper.SetConfigType("yaml")

		viper.SetConfigName(defaultConfigName)
	}

	viper.AutomaticEnv()

	viper.SetEnvPrefix("MINIBLOG")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("Read confing file failed: ", err)
	}

	log.Infow("Using config file: ", viper.ConfigFileUsed())
}

func logOptions() *log.Options {
	return &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}
