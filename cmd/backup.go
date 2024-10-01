package cmd

import (
	"path/filepath"
	"regexp"
	"time"

	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/spf13/cobra"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup",
	RunE:  backupRunE,
}

func backupRunE(_ *cobra.Command, args []string) error {
	// configFile := config.FromCLI().Conf
	// lockFile := config.InstalledLockFile()
	workDir := echo.Wd()
	compressedFile := bootconfig.SoftwareName + `_` + time.Now().Format(`2006_01_02_15_04_05`) + `.zip`
	var saveDir string
	com.SliceExtract(args, &saveDir)
	regexpIgnoreFile := regexp.MustCompile(`^(temp|pid|logs|sessions|dist|html|upgrade_.+\.log\.html|\..+|.*\.zip|.*\.gz)$`)
	var regexpFileName *regexp.Regexp
	if len(saveDir) > 0 {
		compressedFile = filepath.Join(saveDir, compressedFile)
	}
	_, err := com.Zip(workDir, compressedFile, regexpFileName, regexpIgnoreFile)
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	return err
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
