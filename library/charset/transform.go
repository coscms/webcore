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
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/webx-top/com"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var aliases = map[string]string{
	`utf8`:      `utf-8`,
	`hz-gb2312`: `gb2312`,
}

var encodings = map[string]encoding.Encoding{
	`gb18030`: simplifiedchinese.GB18030,
	`gb2312`:  simplifiedchinese.HZGB2312,
	`gbk`:     simplifiedchinese.GBK,
	`utf-8`:   encoding.Nop,
}

func Supported() []string {
	r := make([]string, 0, len(encodings))
	for k := range encodings {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}

func Register(charset string, encoding encoding.Encoding, alias ...string) {
	charset = strings.ToLower(charset)
	encodings[charset] = encoding
	for _, a := range alias {
		a = strings.ToLower(a)
		aliases[a] = charset
	}
}

func Encoding(charset string) encoding.Encoding {
	charset = strings.ToLower(charset)
	if cs, ok := aliases[charset]; ok {
		charset = cs
	}
	if enc, ok := encodings[charset]; ok {
		return enc
	}
	if enc, err := htmlindex.Get(charset); err == nil {
		return enc
	}
	return nil
}

func NewTransformWriter(charset string, dst io.WriteCloser) (io.WriteCloser, error) {
	cs := Encoding(charset)
	if nil == cs {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCharset, charset)
	}
	if cs == encoding.Nop {
		return dst, nil
	}
	return transform.NewWriter(dst, cs.NewDecoder()), nil
}

func NewTransformReader(charset string, src io.Reader) (io.Reader, error) {
	cs := Encoding(charset)
	if nil == cs {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCharset, charset)
	}
	if cs == encoding.Nop {
		return src, nil
	}
	return transform.NewReader(src, cs.NewDecoder()), nil
}

func Transform(charset string, content string) (string, error) {
	r := strings.NewReader(content)
	tr, err := NewTransformReader(charset, r)
	if err != nil {
		return content, err
	}
	b, err := io.ReadAll(tr)
	if err != nil {
		return content, err
	}
	return com.Bytes2str(b), nil
}

func TransformBytes(charset string, content []byte) ([]byte, error) {
	r := bytes.NewReader(content)
	tr, err := NewTransformReader(charset, r)
	if err != nil {
		return content, err
	}
	b, err := io.ReadAll(tr)
	if err != nil {
		return content, err
	}
	return b, nil
}

func NewTransformFunc(charset string) (func(string) (string, error), error) {
	cs := Encoding(charset)
	if nil == cs {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCharset, charset)
	}
	if cs == encoding.Nop {
		return func(v string) (string, error) { return v, nil }, nil
	}
	t := cs.NewDecoder()
	return func(content string) (string, error) {
		r := strings.NewReader(content)
		tr := transform.NewReader(r, t)
		b, err := io.ReadAll(tr)
		if err != nil {
			return content, err
		}
		return com.Bytes2str(b), nil
	}, nil
}

func NewTransformBytesFunc(charset string) (func([]byte) ([]byte, error), error) {
	cs := Encoding(charset)
	if nil == cs {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCharset, charset)
	}
	if cs == encoding.Nop {
		return func(v []byte) ([]byte, error) { return v, nil }, nil
	}
	t := cs.NewDecoder()
	return func(content []byte) ([]byte, error) {
		r := bytes.NewReader(content)
		tr := transform.NewReader(r, t)
		b, err := io.ReadAll(tr)
		if err != nil {
			return content, err
		}
		return b, nil
	}, nil
}
