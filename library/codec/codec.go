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

package codec

import (
	"errors"
	"fmt"
	"strings"
)

type Codec interface {
	Encode(raw string, keys ...string) string
	Decode(encrypted string, keys ...string) string
}

type StringCryptor interface {
	EncryptString(input string) (string, error)
	DecryptString(input string) (string, error)
}

var stringCryptors = map[string]StringCryptor{
	`SM2`: DefaultSM2,
	`RSA`: DefaultRSA,
}

type StringCryptorName [3]byte

func (s StringCryptorName) String() string {
	key := strings.Builder{}
	key.WriteByte(s[0])
	key.WriteByte(s[1])
	key.WriteByte(s[2])
	return key.String()
}

func RegisterStringCryptor(name StringCryptorName, cryptor StringCryptor) {
	stringCryptors[name.String()] = cryptor
}

var ErrUnsupported = errors.New(`unsupported codecs`)

func AutoDecrypt(encrypted string) (string, error) {
	if len(encrypted) < 6 {
		return encrypted, nil
	}
	pre := encrypted[0:4]
	if pre[3] != ':' {
		return encrypted, nil
	}
	pre = pre[0:3]
	cd, ok := stringCryptors[pre]
	if !ok {
		return encrypted, fmt.Errorf(`%w: %s`, ErrUnsupported, pre)
	}
	return cd.DecryptString(encrypted[4:])
}

func AutoEncrypt(plainText string, cryptorType string) (string, error) {
	cd, ok := stringCryptors[cryptorType]
	if !ok {
		return plainText, fmt.Errorf(`%w: %s`, ErrUnsupported, cryptorType)
	}
	encrypted, err := cd.EncryptString(plainText)
	if err != nil {
		return plainText, err
	}
	return cryptorType + `:` + encrypted, err
}
