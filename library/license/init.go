package license

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func init() {
	tplfunc.TplFuncMap[`HasFeature`] = HasFeature
	tplfunc.TplFuncMap[`HasAnyFeature`] = HasAnyFeature
	tplfunc.TplFuncMap[`LicenseDomain`] = Domain
	tplfunc.TplFuncMap[`LicensePackage`] = Package
	tplfunc.TplFuncMap[`LicenseProductURL`] = ProductURL
	tplfunc.TplFuncMap[`LicenseSkipCheck`] = func() bool { return SkipLicenseCheck }
	tplfunc.TplFuncMap[`TrackerURL`] = TrackerURL
	tplfunc.TplFuncMap[`TrackerHTML`] = TrackerHTML
	tplfunc.TplFuncMap[`ProductURL`] = ProductURL

	navigate.FeatureChecker = func(feature string) bool {
		return HasFeature(feature)
	}
}
