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

package license

import (
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
)

var publicKey atomic.Value

func PublicKey() string {
	s, _ := publicKey.Load().(string)
	return s
}

func GetOrLoadPublicKey() string {
	pubKey := PublicKey()
	if len(pubKey) == 0 {
		pubKey = LoadPublicKey()
	}
	return pubKey
}

func LoadPublicKey() string {
	pubKeyFile := filepath.Join(echo.Wd(), `data`, `nging.pem.pub`)
	b, err := os.ReadFile(pubKeyFile)
	if err != nil {
		log.Error(`Failed to reading public key file [ ` + pubKeyFile + ` ]: ` + err.Error())
		return ``
	}
	pubKey := string(b)
	publicKey.Store(pubKey)
	return pubKey
}

func SetPublicKey(pubKey string) {
	publicKey.Store(pubKey)
}
