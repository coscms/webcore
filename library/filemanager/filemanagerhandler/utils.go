package filemanagerhandler

import "github.com/coscms/webcore/library/config"

func Editable(fileName string) (string, bool) {
	return config.FromFile().Sys.Editable(fileName)
}

func Playable(fileName string) (string, bool) {
	return config.FromFile().Sys.Playable(fileName)
}
