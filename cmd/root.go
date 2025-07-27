package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "runcomfy",
	Short: "Analyze ComfyUI workflows and manage dependencies",
	Long: `runcomfy is a CLI tool for analyzing ComfyUI workflow files 
and managing missing custom nodes and models on RunPod instances.

It scans your local ComfyUI installation to identify missing dependencies
from workflow files and helps you download and install them.`,
	Version: "1.0.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.runcomfy.yaml)")
	rootCmd.PersistentFlags().StringP("comfyui-path", "p", "/workspace/ComfyUI", "path to ComfyUI installation")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "output format (table, json)")

	viper.BindPFlag("comfyui-path", rootCmd.PersistentFlags().Lookup("comfyui-path"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".runcomfy")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}