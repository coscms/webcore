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

package charset

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	runewidth "github.com/mattn/go-runewidth"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

var ErrUnsupportedCharset = errors.New(`charset is unsupported`)

func Convert(fromEnc string, toEnc string, b []byte) ([]byte, error) {
	toCS := Encoding(toEnc)
	if nil == toCS {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCharset, toEnc)
	}
	if toCS == encoding.Nop {
		return TransformBytes(fromEnc, b)
	}
	fromCS := Encoding(fromEnc)
	if nil == fromCS {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCharset, fromEnc)
	}
	if fromCS == encoding.Nop {
		r := bytes.NewReader(b)
		enc := transform.NewReader(r, toCS.NewEncoder())
		return io.ReadAll(enc)
	}
	r := bytes.NewReader(b)
	dec := transform.NewReader(r, fromCS.NewDecoder())
	enc := transform.NewReader(dec, toCS.NewEncoder())
	return io.ReadAll(enc)
}

func Validate(enc string) bool {
	return Encoding(enc) != nil
}

func Truncate(str string, width int) string {
	w := 0
	b := []byte(str)
	var buf bytes.Buffer
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		rw := runewidth.RuneWidth(r)
		if w+rw > width {
			break
		}
		buf.WriteRune(r)
		w += rw
		b = b[size:]
	}
	return buf.String()
}

func With(str string) int {
	return runewidth.StringWidth(str)
}

func RuneWith(str string) int {
	return utf8.RuneCountInString(str)
}
