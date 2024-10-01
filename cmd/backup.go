package cmd

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/spf13/cobra"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var backupCmd = &cobra.Command{
	Use:     "backup",
	Short:   "Backup",
	Example: os.Args[0] + " backup [saveDir]",
	RunE:    backupRunE,
}

var backupIgnore string

func backupRunE(_ *cobra.Command, args []string) error {
	workDir := echo.Wd()
	compressedFile := bootconfig.SoftwareName + `_` + time.Now().Format(`2006_01_02_15_04_05`)
	if com.IsWindows {
		compressedFile += `.zip`
	} else {
		compressedFile += `.tar.gz`
	}
	var saveDir string
	com.SliceExtract(args, &saveDir)
	if len(saveDir) > 0 {
		compressedFile = filepath.Join(saveDir, compressedFile)
	}
	if len(backupIgnore) > 0 {
		backupIgnore = `|` + backupIgnore
	}
	regexpIgnoreFile := regexp.MustCompile(`^(temp|pid|logs|sessions|dist|html|vendor|application|go\.mod|go\.sum|nohup\.out|upgrade_.+\.log\.html|\..+|.*\.zip|.*\.gz|.*\.log|.*\.go` + backupIgnore + `)$`)
	var regexpFileName *regexp.Regexp

	var err error
	if strings.EqualFold(filepath.Base(compressedFile), `.zip`) {
		_, err = com.Zip(workDir, compressedFile, regexpFileName, regexpIgnoreFile)
	} else {
		err = com.TarGz(workDir, compressedFile, regexpFileName, regexpIgnoreFile)
	}
	if err != nil {
		com.ExitOnFailure(err.Error(), 1)
	}
	return err
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringVar(&backupIgnore, `ignore`, backupIgnore, `忽略的文件正则表达式`)
}
