package project

import (
	"github.com/spf13/cobra"

	"github.com/onlishop/onlishop-cli/shop"
)

var projectConfigPath string

var projectRootCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage your Onlishop Project",
}

func Register(rootCmd *cobra.Command) {
	rootCmd.AddCommand(projectRootCmd)
	projectRootCmd.PersistentFlags().StringVar(&projectConfigPath, "project-config", shop.DefaultConfigFileName(), "Path to config")
}
