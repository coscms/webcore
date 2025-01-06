package config

import "github.com/coscms/webcore/library/codec"

func init() {
	codec.RegisterStringCryptor([3]byte{'A', 'E', 'S'}, NewAES())
}

func NewAES() *aes {
	return &aes{}
}

type aes struct {
}

func (r *aes) getConfig() (*Config, error) {
	cfg := FromFile()
	if cfg != nil {
		return cfg, nil
	}
	return InitConfig()
}

// EncryptString 加密
func (r *aes) EncryptString(input string) (string, error) {
	cfg, err := r.getConfig()
	if err != nil {
		return ``, err
	}
	encrypted := cfg.Encode256(input)
	return encrypted, nil
}

// DecryptString 解密
func (r *aes) DecryptString(input string) (string, error) {
	cfg, err := r.getConfig()
	if err != nil {
		return ``, err
	}
	decrypted := cfg.Decode256(input)
	return decrypted, nil
}
