/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"gitee.com/admpub/certmagic"
	figure "github.com/admpub/go-figure"
	"github.com/admpub/go-ps"
	"github.com/admpub/log"
	"github.com/admpub/service"
	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/subdomains"

	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/config/startup"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/license"
	"github.com/coscms/webcore/library/msgbox"
	"github.com/coscms/webcore/library/selfupdate"
)

// Nging 启动入口

// rootCmd represents the base command when called without any subcommands
var rootCmd = NewRoot()
var dumpCli bool

func NewRoot() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	return &cobra.Command{
		Use:          filepath.Base(os.Args[0]),
		Short:        ``,
		Long:         ``,
		SilenceUsage: true,
		RunE:         rootRunE,
	}
}

func callStartup() error {
	if !service.Interactive() { //非交互模式(即后台服务模式)
		return nil
	}
	pproc, err := ps.FindProcess(os.Getppid())
	if err != nil {
		return err
	}
	if pproc == nil {
		return nil
	}
	executor := filepath.Base(pproc.Executable())
	expected := []string{filepath.Base(os.Args[0]), `startup`, `go`}
	var ext string
	if com.IsWindows {
		ext = `.exe`
		expected = append(expected, `startup`+ext, `go`+ext)
	}
	if com.InSlice(executor, expected) {
		return nil
	}
	executable, err := osext.Executable()
	if err != nil {
		return err
	}
	workDir := filepath.Dir(executable)
	executable = filepath.Join(workDir, `startup`+ext)
	if !com.IsFile(executable) {
		return nil
	}
	procArgs := []string{executable}
	if len(os.Args) > 1 {
		procArgs = append(procArgs, os.Args[1:]...)
	}
	_, err = os.StartProcess(executable, procArgs, &os.ProcAttr{
		Dir:   workDir,
		Env:   os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Sys:   selfupdate.NewSysProcAttr(),
	})
	if err != nil {
		return err
	}
	signal.Reset()
	os.Exit(0)
	return err
}

func rootRunE(cmd *cobra.Command, args []string) error {
	if err := callStartup(); err != nil {
		return err
	}
	if !config.Version.Licensed {
		message := `Invalid license!
授权无效!
		
If you have already purchased a license, please place the ` + license.FileName() + ` to:
如果您已经购买了授权，请将授权证书` + license.FileName() + `放到：
		
%s`
		msgbox.Error("WARNING",
			message,
			license.FilePath())

		fmt.Println(``)
		fmt.Println(`To purchase a license, please go to our official website:`)
		fmt.Println(`购买授权请前往官方网站：`)
		fmt.Println(``)
		fmt.Println(license.ProductDetailURL())
		if bootconfig.MustLicensed {
			return nil
		}
	}

	//独立模块
	if config.FromCLI().OnlyRunServer() {
		bootconfig.SetServerType(config.FromCLI().Type)
		httpserver.Clear()
		return nil
	}

	bootconfig.SetServerType(`web`)
	//Manager
	config.FromCLI().RunStartup()

	if config.IsInstalled() && bootconfig.AutoUpgradeDBStruct && !config.FromFile().Sys.DisableAutoUpgradeDB {
		if err := bootconfig.Upgrade(); err != nil && os.ErrNotExist != err {
			log.Error(`upgrade.sql: `, err)
		}
	}

	// LOGO
	fmt.Println(strings.TrimSuffix(figure.NewFigure(bootconfig.SoftwareName, `big`, false).String(), "\n"), config.Version.VString()+"\n")

	bootconfig.WebStart()
	startup.FireBefore(`web`)
	if config.IsInstalled() {
		startup.FireAfter(`web.installed`)
	}
	defer func() {
		startup.FireAfter(`web`)
	}()

	c := &engine.Config{
		ReusePort:          true,
		TLSAuto:            config.FromFile().Sys.SSLAuto,
		TLSEmail:           config.FromFile().Sys.SSLEmail,
		TLSHosts:           config.FromFile().Sys.SSLHosts,
		TLSCacheDir:        config.FromFile().Sys.SSLCacheDir,
		TLSCertFile:        config.FromFile().Sys.SSLCertFile,
		TLSKeyFile:         config.FromFile().Sys.SSLKeyFile,
		MaxRequestBodySize: config.FromFile().GetMaxRequestBodySize(),
	}
	c.Address = fmt.Sprintf(`%s:%v`, config.FromCLI().Address, config.FromCLI().Port)
	hasCert := (len(c.TLSCertFile) > 0 && len(c.TLSKeyFile) > 0)
	//c.TLSAuto = true
	if c.TLSAuto || hasCert {
		if config.FromCLI().Port == 80 {
			if c.TLSAuto {
				echo.PanicIf(initCertMagic(c))
				//c.SupportAutoTLS(nil, config.FromFile().Sys.SSLHosts...)
			} else {
				c.Address = fmt.Sprintf(`%s:443`, config.FromCLI().Address)
				e2 := echo.New()
				e2.Use(middleware.HTTPSRedirect(), middleware.LogWithConfig(middleware.LogConfig{Skipper: config.FromFile().Sys.HTTPLogSkipper}), middleware.Recover())
				go e2.Run(standard.New(fmt.Sprintf(`%s:80`, config.FromCLI().Address)))
			}
			subdomains.Default.Protocol = `https`
		}
	}
	if len(bootconfig.Welcome) > 0 {
		now := time.Now()
		msgbox.Success(`Welcome`,
			bootconfig.Welcome,
			config.Version.VString(),
			now.Format("Monday, 02 Jan 2006"))
	}
	subdomains.Default.SetDebug(config.FromFile().Debug())
	echo.Fire(`nging.httpserver.run.before`)

	serverEngine := standard.NewWithConfig(c)
	go handleSignal(serverEngine)
	subdomains.Default.Run(serverEngine)
	return nil
}

func initCertMagic(c *engine.Config) error {
	fileStorage := &certmagic.FileStorage{
		Path: filepath.Join(echo.Wd(), `data`, `cache`, `certmagic`),
	}
	if err := com.MkdirAll(fileStorage.Path, os.ModePerm); err != nil {
		return err
	}
	if bootconfig.Develop { // use the staging endpoint while we're developing
		certmagic.Default.CA = certmagic.LetsEncryptStagingCA
	} else {
		certmagic.Default.CA = certmagic.LetsEncryptProductionCA
	}
	certmagic.Default.Email = c.TLSEmail
	certmagic.Default.Agreed = true
	certmagic.Default.Storage = fileStorage
	ln, err := certmagic.Listen(c.TLSHosts)
	if err == nil {
		c.SetListener(ln)
	}
	return err
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if dumpCli {
		writeReceived()
	}
	config.FromCLI().InitFlag(rootCmd.PersistentFlags())
	Init()
	if len(rootCmd.Use) == 0 {
		rootCmd.Use = os.Args[0]
	}

	if err := rootCmd.Execute(); err != nil {
		com.ExitOnFailure(err.Error() + "\n")
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolVar(&dumpCli, "dumpcli", false, "--dumpcli false")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
