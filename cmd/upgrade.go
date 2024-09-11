package cmd

import (
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/config"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade Database Struct",
	RunE:  upgradeRunE,
}

func upgradeRunE(cmd *cobra.Command, args []string) error {
	conf, err := config.InitConfig()
	config.MustOK(err)
	conf.AsDefault()
	return bootconfig.Upgrade()
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
