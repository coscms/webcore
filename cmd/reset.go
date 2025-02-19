package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/admpub/oauth2/v4/errors"
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/config"
	"github.com/spf13/cobra"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset config",
	RunE:  resetRunE,
	Example: rootCmd.Use + ` reset captcha
` + rootCmd.Use + ` reset password <username> <new_password>
` + rootCmd.Use + ` reset u2f <username>`,
}

func resetRunE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Usage()
	}
	conf, err := config.InitConfig()
	if err != nil {
		return err
	}
	conf.AsDefault()
	switch args[0] {
	case `captcha`: // 重置验证码配置
		ctx := defaults.NewMockContext()
		m := dbschema.NewNgingConfig(ctx)
		var affected int64
		affected, err = m.UpdatexField(nil, `value`, `default`, db.And(
			db.Cond{`key`: `type`},
			db.Cond{`group`: `captcha`},
		))
		if err != nil {
			return err
		}
		if affected == 0 {
			err = errors.New("no record updated")
		}
	case `password`: // 重置密码配置
		ctx := defaults.NewMockContext()
		if len(args) < 2 {
			return errors.New("user name required")
		}
		if len(args) < 3 {
			return errors.New("user new password required")
		}
		flagFile := filepath.Join(echo.Wd(), `cli_can_modify_password`)
		if !com.FileExists(flagFile) {
			return errors.New("can not modify password. can not find file: " + flagFile)
		}
		var b []byte
		b, err = os.ReadFile(flagFile)
		if err != nil {
			return err
		}
		cliCanModifyPassword, _ := strconv.ParseBool(string(b))
		if !cliCanModifyPassword {
			return errors.New("can not modify password. The content of the file cli_can_modify_password is not true.")
		}
		username := strings.TrimSpace(args[1])
		password := strings.TrimSpace(args[2])
		m := dbschema.NewNgingUser(ctx)
		err = m.Get(func(r db.Result) db.Result {
			return r.Select(`id`, `salt`)
		}, `username`, username)
		if err != nil {
			return err
		}
		var affected int64
		password = com.MakePassword(password, m.Salt)
		affected, err = m.UpdatexField(nil, `password`, password, `id`, m.Id)
		if err != nil {
			return err
		}
		if affected == 0 {
			err = errors.New("no record updated")
		}
	case `u2f`: // 取消用户的U2F绑定
		ctx := defaults.NewMockContext()
		if len(args) < 2 {
			return errors.New("user name required")
		}
		flagFile := filepath.Join(echo.Wd(), `cli_can_reset_u2f`)
		if !com.FileExists(flagFile) {
			return errors.New("can not reset u2f. can not find file: " + flagFile)
		}
		var b []byte
		b, err = os.ReadFile(flagFile)
		if err != nil {
			return err
		}
		cliCanResetU2F, _ := strconv.ParseBool(string(b))
		if !cliCanResetU2F {
			return errors.New("can not reset u2f. The content of the file cli_can_reset_u2f is not true.")
		}
		username := strings.TrimSpace(args[1])
		userM := dbschema.NewNgingUser(ctx)
		err = userM.Get(func(r db.Result) db.Result {
			return r.Select(`id`)
		}, `username`, username)
		if err != nil {
			return err
		}
		m := dbschema.NewNgingUserU2f(ctx)
		var affected int64
		affected, err = m.Deletex(nil, db.And(
			db.Cond{`uid`: userM.Id},
			db.Cond{`step`: db.In([]uint{0, 2})},
		))
		if err != nil {
			return err
		}
		if affected == 0 {
			err = errors.New("no record deleted")
		}
	default:
		err = fmt.Errorf("unsupported config: %s", args[0])
	}
	if err == nil {
		fmt.Println("reset success:", args[0])
	}
	return err
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
