package config

import (
	"strings"

	"github.com/webx-top/echo"
)

var (
	actGroups          = []string{`base`, `smtp`, `log`}
	onKeySetSettings   = map[string][]func(Diff) error{}
	onGroupSetSettings = map[string][]func(Diffs) error{}
)

func SettingsInitialized() bool {
	return echo.Has(SettingName)
}

// OnGroupSetSettings 注册配置组的变动事件
// OnGroupSetSettings(`base`,fn)
func OnGroupSetSettings(group string, fn func(Diffs) error) {
	if SettingsInitialized() {
		values := FromDB(group)
		diffs := Diffs{}
		for k, v := range values {
			diffs[k] = &Diff{
				New:    v,
				IsDiff: true,
			}
		}
		fn(diffs)
	}
	if _, ok := onGroupSetSettings[group]; !ok {
		onGroupSetSettings[group] = []func(Diffs) error{}
	}
	onGroupSetSettings[group] = append(onGroupSetSettings[group], fn)
}

// OnKeySetSettings 注册配置组中某个配置的变动事件
// OnKeySetSettings(`base.debug`,fn)
func OnKeySetSettings(groupAndKey string, fn func(Diff) error) {
	if SettingsInitialized() {
		args := strings.SplitN(groupAndKey, `.`, 2)
		values := FromDB(args[0])
		var val interface{}
		if len(args) == 2 {
			val = values.Get(args[1])
		} else {
			val = values
		}
		fn(Diff{
			New:    val,
			IsDiff: true,
		})
	}
	if _, ok := onKeySetSettings[groupAndKey]; !ok {
		onKeySetSettings[groupAndKey] = []func(Diff) error{}
	}
	onKeySetSettings[groupAndKey] = append(onKeySetSettings[groupAndKey], fn)
}

func FireInitSettings(configs echo.H) error {
	for group, fnList := range onGroupSetSettings {
		values := configs.GetStore(group)
		diffs := Diffs{}
		for k, v := range values {
			diffs[k] = &Diff{
				New:    v,
				IsDiff: true,
			}
		}
		for _, fn := range fnList {
			if err := fn(diffs); err != nil {
				return err
			}
		}
	}
	for groupAndKey, fnList := range onKeySetSettings {
		args := strings.SplitN(groupAndKey, `.`, 2)
		values := configs.GetStore(args[0])
		var val interface{}
		if len(args) == 2 {
			val = values.Get(args[1])
		} else {
			val = values
		}
		for _, fn := range fnList {
			if err := fn(Diff{
				New:    val,
				IsDiff: true,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func FireSetSettings(group string, diffs Diffs) error {
	if fnList, ok := onGroupSetSettings[group]; ok {
		for _, fn := range fnList {
			if err := fn(diffs); err != nil {
				return err
			}
		}
	}
	for key, diff := range diffs {
		k := group + `.` + key
		if fnList, ok := onKeySetSettings[k]; ok {
			for _, fn := range fnList {
				if err := fn(*diff); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
