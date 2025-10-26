package upload

import (
	"io"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/com"
)

// DetectMIMETypeByReader 通过io.Reader检测文件的MIME类型
func DetectMIMETypeByReader(r io.Reader) (string, error) {
	mime, err := mimetype.DetectReader(r)
	if err != nil {
		return ``, err
	}
	return mime.String(), nil
}

// DetectMIMETypeByFile 通过文件路径检测文件的MIME类型
func DetectMIMETypeByFile(file string) (string, error) {
	mime, err := mimetype.DetectFile(file)
	if err != nil {
		return ``, err
	}
	return mime.String(), nil
}

// DetectMIMETypeByBytes 通过字节数据检测文件的MIME类型
func DetectMIMETypeByBytes(data []byte) string {
	mime := mimetype.Detect(data)
	return mime.String()
}

// IsMIMEType 判断MIME类型是否属于指定的文件类型
func IsMIMEType(mimeType string, fileType string) bool {
	return InMIMEType(mimeType, uploadClient.FileTypeMimeKeywords[fileType])
}

// InMIMEType 判断MIME类型是否在允许的MIME类型列表中
func InMIMEType(mimeType string, allowedMIMEs []string) bool {
	if com.InSlice(mimeType, allowedMIMEs) {
		return true
	}
	for _, v := range strings.SplitN(mimeType, `/`, 2) {
		if com.InSlice(v, allowedMIMEs) {
			return true
		}
	}
	return false
}

// IsImageMIMEType 判断MIME类型是否为图片类型
func IsImageMIMEType(mimeType string) bool {
	return IsMIMEType(mimeType, `image`)
}

// IsVideoMIMEType 判断MIME类型是否为视频类型
func IsVideoMIMEType(mimeType string) bool {
	return IsMIMEType(mimeType, `video`)
}

// IsAudioMIMEType 判断MIME类型是否为音频类型
func IsAudioMIMEType(mimeType string) bool {
	return IsMIMEType(mimeType, `audio`)
}

// IsOfficeDocMIMEType 判断MIME类型是否为Office文档doc类型
func IsOfficeDocMIMEType(mimeType string) bool {
	return IsMIMEType(mimeType, `doc`)
}

// IsOfficeXlsMIMEType 判断MIME类型是否为Office文档xls类型
func IsOfficeXlsMIMEType(mimeType string) bool {
	return IsMIMEType(mimeType, `xls`)
}

// IsOfficePptMIMEType 判断MIME类型是否为Office文档ppt类型
func IsOfficePptMIMEType(mimeType string) bool {
	return IsMIMEType(mimeType, `ppt`)
}

// IsArchiveMIMEType 判断MIME类型是否为压缩包类型
func IsArchiveMIMEType(mimeType string) bool {
	return IsMIMEType(mimeType, `archive`)
}

// IsPdfMIMEType 判断MIME类型是否为PDF类型
func IsPdfMIMEType(mimeType string) bool {
	return IsMIMEType(mimeType, `pdf`)
}

// IsTorrentMIMEType 判断MIME类型是否为BT种子类型
func IsTorrentMIMEType(mimeType string) bool {
	return IsMIMEType(mimeType, `bt`)
}

// IsPhotoshopMIMEType 判断MIME类型是否为Photoshop类型
func IsPhotoshopMIMEType(mimeType string) bool {
	return IsMIMEType(mimeType, `photoshop`)
}
