package license

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/admpub/log"
	"github.com/admpub/pp/ppnocolor"
	"github.com/coscms/webcore/library/config"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"golang.org/x/net/publicsuffix"
)

func init() {
	config.Version.BuildOS = runtime.GOOS
	config.Version.BuildArch = runtime.GOARCH
	config.Version.Package = `free`
	config.Version.Number = `5.0.0`
	/*/
	(&ServerURL{
		Tracker: `http://nging.coscms.com/product/script/nging/tracker.js`,
		Product: `http://nging.coscms.com/product/detail/nging`,
		License: `http://nging.coscms.com/product/license/nging`,
		Version: `http://nging.coscms.com/product/version/nging`,
	}).Apply()
	//*/
}

func TestCpuID(t *testing.T) {
	cpuID, err := CpuID()
	assert.NoError(t, err)
	t.Logf(`cpuID: %s`, cpuID)

	mp := SplitChecksums(`d41d8cd98f00b204e9800998ecf8427e  file1.txt
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855  file2.txt
`)
	//ppnocolor.Println(mp)
	assert.Equal(t, map[string]string{
		"file1.txt": "d41d8cd98f00b204e9800998ecf8427e",
		"file2.txt": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	}, mp)
}

func TestEmptyLicense(t *testing.T) {
	SetEmptyLicenseFeature(`A`, `B`, `C`)
	assert.True(t, HasFeature(`A`))
}

func TestParseVersionInfo(t *testing.T) {
	output := `[DEBUG]88202|2025-05-06 14:47:23|Debug|app|******
[ERROR]88202|2025-05-06 14:47:23|bindata_none.go:55|Error|app|*****
Nging v5.3.3-dev licensed(free)
Schema: v7.7001
Build: 20250424162337`
	ver, err := parseVersionInfo(output)
	assert.NoError(t, err)
	assert.Equal(t, &config.VersionInfo{
		Name:      `Nging`,
		Label:     `dev`,
		Number:    `5.3.3`,
		Package:   `free`,
		DBSchema:  7.7001,
		BuildTime: `20250424162337`,
	}, ver)
}

func TestLicenseDownload(t *testing.T) {
	return
	// dirEntries, _ := os.ReadDir(filepath.Join(echo.Wd(), `.`))
	// for _, dirEntry := range dirEntries {
	// 	t.Log(dirEntry.Name())
	// }
	// return
	err := Download(nil)
	if err != nil {
		panic(err)
	}
}

func TestLicenseLatestVersion(t *testing.T) {
	return
	defer log.Close()
	ctx := defaults.NewMockContext()
	info, err := LatestVersion(ctx, ``, true)
	if err != nil {
		panic(err)
	}
	ppnocolor.Println(info)
	err = info.Extract()
	if err != nil {
		panic(err)
	}
	ppnocolor.Println(info.extractedDir)
	ppnocolor.Println(info.executable)
	ngingDir, err := filepath.Abs(`./testdata`)
	if err != nil {
		panic(err)
	}
	echo.SetWorkDir(ngingDir)
	args := append([]string{}, os.Args[1:]...)
	os.Args = os.Args[0:1]
	os.Args = append(os.Args, `-p`, `29990`)
	err = info.Upgrade(ctx, ngingDir, `default`)
	if err != nil {
		panic(err)
	}
	os.Args = append(os.Args[0:1], args...)
}

func TestLicenseValidateFromOfficial(t *testing.T) {
	return
	err := validateFromOfficial(nil)
	if err != nil {
		panic(err)
	}
}

func TestLicenseEqDomain(t *testing.T) {
	defer log.Close()
	assert.True(t, EqDomain(`www.webx.top`, `webx.top`))

	domain, err := publicsuffix.EffectiveTLDPlusOne(`www.webx.top`)
	assert.Nil(t, err)
	assert.Equal(t, `webx.top`, domain)

	domain, err = publicsuffix.EffectiveTLDPlusOne(`www.abc.com.cn`)
	assert.Nil(t, err)
	assert.Equal(t, `abc.com.cn`, domain)

	domain, err = publicsuffix.EffectiveTLDPlusOne(`com.cn`)
	assert.NotNil(t, err)
	assert.Equal(t, ``, domain)

	publicSuffix, icann := publicsuffix.PublicSuffix(`www.webx.top`)
	assert.True(t, icann)
	assert.Equal(t, `top`, publicSuffix)

	publicSuffix, icann = publicsuffix.PublicSuffix(`www.webx.x`)
	assert.False(t, icann)
	assert.Equal(t, `x`, publicSuffix)
}

func TestUserAgent(t *testing.T) {
	version = `1.2.3-beta`
	packageName = `free`
	emptyLicense.Info.Name = `user`
	ua := MakeUserAgent()
	assert.Equal(t, `nging-free/1.2.3-beta user`, ua)
	t.Log(ua)
	r := ParseUserAgent(ua)
	assert.Equal(t, UserAgentRaw{
		Product: `nging`,
		Package: `free`,
		Version: `1.2.3-beta`,
		User:    `user`,
	}, r)
}
