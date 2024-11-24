package settings

import "github.com/coscms/webcore/dbschema"

type Config struct {
	Group   string
	Items   map[string]*dbschema.NgingConfig
	Forms   []*SettingForm
	Encoder Encoder
	Decoder Decoder
}

func NewConfig(group string) *Config {
	return &Config{
		Group: group,
		Items: map[string]*dbschema.NgingConfig{},
	}
}

func (c *Config) AddItem(item *dbschema.NgingConfig) *Config {
	item.Group = c.Group
	c.Items[item.Key] = item
	return c
}

func (c *Config) AddForm(cfgs ...*SettingForm) *Config {
	c.Forms = append(c.Forms, cfgs...)
	return c
}

func (c *Config) SetTransferType(dest interface{}) *Config {
	rType := GetReflectType(dest)
	c.Encoder = MakeEncoder(rType)
	c.Decoder = MakeDecoder(rType)
	return c
}

func (c *Config) SetEncoder(encoder Encoder) *Config {
	c.Encoder = encoder
	return c
}

func (c *Config) SetDecoder(decoder Decoder) *Config {
	c.Decoder = decoder
	return c
}

func (c *Config) Apply() {
	if c.Encoder != nil {
		RegisterEncoder(c.Group, c.Encoder)
	}
	if c.Decoder != nil {
		RegisterDecoder(c.Group, c.Decoder)
	}
	AddDefaultConfig(c.Group, c.Items)
	Register(c.Forms...)
}
