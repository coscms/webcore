package formbuilder_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/admpub/log"
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/param"
)

var (
	_ echo.BeforeValidate = &TestRequest{}
	_ echo.AfterValidate  = &TestRequest{}
)

type TestRequest struct {
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

func (t *TestRequest) BeforeValidate(ctx echo.Context) error {
	return nil
}

func (t *TestRequest) AfterValidate(ctx echo.Context) error {
	return nil
}

func (t *TestRequest) Short_() string {
	return `test`
}

func TestFormbuilder(t *testing.T) {
	defer log.Close()
	com.MkdirAll(`testdata/template`, os.ModePerm)
	d := render.New(`standard`, `testdata/template`)
	d.Init()
	defaults.Use(render.Middleware(d))
	dbschema.DBI.Fields[`test`] = map[string]*factory.FieldInfo{
		`name`: &factory.FieldInfo{
			Name:         `name`,
			DataType:     `varchar`,
			GoType:       `string`,
			GoName:       `Name`,
			Multilingual: true,
		},
	}
	langs := echo.NewKVData()
	langs.Add(`en`, `English`)
	langs.Add(`zh-CN`, `简体中文`)
	bean := &TestRequest{}
	ctx := defaults.NewMockContext()
	ctx.SetRenderer(d)
	form := formbuilder.New(ctx, bean,
		formbuilder.LanguagesGetter(func(ctx echo.Context) *echo.KVData {
			return langs
		}, `en`),
		formbuilder.ConfigFile(`test`))
	form.OnPost(func() error {
		var err error
		fmt.Printf("%#v\n", bean)
		return err
	})
	form.Generate()
	//fmt.Printf("%#v\n", ctx.Get(`forms`))
	assert.Equal(t, form.Forms, ctx.Get(`forms`))
	htmlResult := string(form.Render())
	fmt.Println(htmlResult)
	var spaceClearRegex = regexp.MustCompile(`(>)[\s]+(&|<)`)
	htmlResult = spaceClearRegex.ReplaceAllString(htmlResult, `$1$2`)
	expected := `<form generator="forms" class="form-horizontal" id="Forms" role="form" method="POST" action="" required-redstar="true">
	<div class="form-group">
<label class="col-sm-2 control-label">Name</label>
<div class="col-sm-8">
	<input type="text" name="Name" class="form-control">
</div>
</div>
	<div class="form-group">
<label class="col-sm-2 control-label">Age</label>
<div class="col-sm-8">
	<input type="text" name="Age" class="form-control" value="0">
</div>
</div>
	<div class="form-group form-submit-group">
	<div class="col-md-offset-2 col-md-10">
<button type="submit" class='btn btn-lg btn-primary'>
<i class="fa fa-check"></i>
    Submit
</button>
			&nbsp; &nbsp; &nbsp;
<button type="reset" class='btn btn-lg'>
<i class="fa fa-undo"></i>
    Reset
</button>
			&nbsp; &nbsp; &nbsp;
	</div>
</div>
</form>`
	expected = spaceClearRegex.ReplaceAllString(expected, `$1$2`)
	//assert.Equal(t, expected, htmlResult)

	// -----
	expectedReq := &TestRequest{
		Name: `test`,
		Age:  123,
	}
	ctx.Request().SetMethod(`POST`)
	ctx.Request().Form().Set(`Language[en][name]`, expectedReq.Name)
	ctx.Request().Form().Set(`age`, param.AsString(expectedReq.Age))
	err := form.RecvSubmission()
	if form.Exited() {
		err = form.Error()
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedReq, bean)
}
