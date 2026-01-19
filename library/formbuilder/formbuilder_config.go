package formbuilder

import (
	"fmt"
	"os"
	"path"

	"github.com/coscms/forms"
	"github.com/coscms/forms/common"
	formsconfig "github.com/coscms/forms/config"
	"github.com/webx-top/echo/middleware/render/driver"
	"gopkg.in/yaml.v3"
)

// ParseConfigFile 解析配置文件 xxx.form.json
func (f *FormBuilder) ParseConfigFile(jsonformat ...bool) (*formsconfig.Config, error) {
	renderer, ok := f.ctx.Renderer().(driver.Driver)
	if !ok {
		return nil, fmt.Errorf(`FormBuilder: Expected renderer is "driver.Driver", but got "%T"`, f.ctx.Renderer())
	}
	var isJSON bool
	var configFile string
	if len(jsonformat) > 0 {
		isJSON = jsonformat[0]
	} else {
		switch path.Ext(f.configFile) {
		case `.json`:
			isJSON = true
			configFile = f.configFile
		case `.yaml`, `.yml`:
			isJSON = false
			configFile = f.configFile
		default:
			isJSON = false
		}
	}
	if len(configFile) == 0 {
		configFile = f.configFile + `.form`
		if isJSON {
			configFile += `.json`
		} else {
			configFile += `.yaml`
		}
	}
	configFile = renderer.TmplPath(f.ctx, configFile)
	if len(configFile) == 0 {
		return nil, ErrJSONConfigFileNameInvalid
	}
	var cfg *formsconfig.Config
	b, err := renderer.RawContent(configFile)
	if err != nil || len(b) == 0 {
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf(`read file %s: %w`, configFile, err)
		}
		if renderer.Manager() == nil {
			return nil, fmt.Errorf(`renderer.Manager() is nil: %s`, configFile)
		}
		cfg = f.ToConfig()
		if f.snippet {
			cfg.Template = `allfields`
			cfg.WithButtons = false
		}
		if isJSON {
			b, err = f.ToJSONBlob(cfg)
			if err != nil {
				return nil, fmt.Errorf(`[form.ToJSONBlob] %s: %w`, configFile, err)
			}
		} else {
			b, err = yaml.Marshal(cfg)
			if err != nil {
				return nil, fmt.Errorf(`[form:yaml.Marshal] %s: %w`, configFile, err)
			}
		}
		err = renderer.Manager().SetTemplate(configFile, b)
		if err != nil {
			return nil, fmt.Errorf(`%s: %w`, configFile, err)
		}
		f.ctx.Logger().Infof(f.ctx.T(`生成表单配置文件“%v”成功。`), configFile)
	} else {
		if isJSON {
			cfg, err = forms.Unmarshal(b, configFile)
			if err != nil {
				return nil, fmt.Errorf(`[forms.Unmarshal] %s: %w`, configFile, err)
			}
		} else {
			cfg, err = common.GetOrSetCachedConfig(configFile, func() (*formsconfig.Config, error) {
				cfg := &formsconfig.Config{}
				err := yaml.Unmarshal(b, cfg)
				return cfg, err
			})
			if err != nil {
				return nil, fmt.Errorf(`[form:yaml.Unmarshal] %s: %w`, configFile, err)
			}
		}
	}
	if cfg != nil {
		return cfg.Clone(), err
	}
	cfg = f.NewConfig()
	return cfg, err
}

// SetConfig sets the form builder configuration and returns the FormBuilder instance for method chaining.
func (f *FormBuilder) SetConfig(cfg *formsconfig.Config) *FormBuilder {
	f.config = cfg
	return f
}

// InitConfig initializes the FormBuilder configuration by either parsing the config file
// or cloning the existing config. It also handles language set conversion if languages are specified.
// Returns an error if config parsing fails.
func (f *FormBuilder) InitConfig() error {
	var cfg *formsconfig.Config
	var err error
	if f.config == nil {
		cfg, err = f.ParseConfigFile()
		if err != nil {
			return err
		}
	} else {
		cfg = f.config.Clone()
	}

	if f.configPrepare != nil {
		err = f.configPrepare(cfg)
		if err != nil {
			return err
		}
	}

	if f.snippet {
		f.setSnippetConfig(cfg)
	}

	if f.Languages() != nil {
		f.toLangset(cfg)
	}

	f.Init(cfg)
	return err
}
