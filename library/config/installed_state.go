package config

import (
	"database/sql"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/setup"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

const LockFileName = `installed.lock`

var (
	Installed             sql.NullBool
	installedSchemaVer    float64
	installedPkgSchemaVer = map[string]float64{}
	installedTime         time.Time
	sqlCollection         = NewSQLCollection().RegisterInstall(`nging`, setup.InstallSQL)
)

// GetInstalledPkgSchemaVer returns the version of the given package in the installed lock file.
// If the package is not found in the installed lock file, it returns -1.
func GetInstalledPkgSchemaVer(pkg string) float64 {
	v, ok := installedPkgSchemaVer[pkg]
	if !ok {
		return -1
	}
	return v
}

// GetSQLCollection returns the SQLCollection that is used to store the install SQLs,
// the insert SQLs, and the pre-upgrade SQLs of the system.
// The SQLCollection is used by the system to execute the SQLs during the installation and upgrade process.
// The SQLCollection is also used by the system to get the install SQLs, the insert SQLs, and the pre-upgrade SQLs of the system.
// The SQLCollection is registered with the install SQLs, the insert SQLs, and the pre-upgrade SQLs of the system using the RegisterInstallSQL, RegisterInsertSQL, and RegisterPreupgradeSQL functions.
func GetSQLCollection() *SQLCollection {
	return sqlCollection
}

// RegisterInstallSQL registers an install SQL string to the given project.
// The install SQL string will be executed when installing the project.
// The project name must be unique, and the install SQL string must be a valid SQL string.
// The install SQL string will be appended to the existing install SQLs of the project.
// If the project does not exist, it will be created.
// If the project already exists, the install SQL string will be appended to the existing install SQLs of the project.
// The RegisterInstallSQL function is used by the system to register the install SQLs of the system.
// The RegisterInstallSQL function is also used by the plugins to register the install SQLs of the plugins.
func RegisterInstallSQL(project string, installSQL string) {
	sqlCollection.RegisterInstall(project, installSQL)
}

// RegisterInsertSQL registers an insert SQL string to the given project.
// The insert SQL string will be executed when inserting data into the project.
// The project name must be unique, and the insert SQL string must be a valid SQL string.
// The insert SQL string will be appended to the existing insert SQLs of the project.
// If the project does not exist, it will be created.
// If the project already exists, the insert SQL string will be appended to the existing insert SQLs of the project.
// The RegisterInsertSQL function is used by the system to register the insert SQLs of the system.
// The RegisterInsertSQL function is also used by the plugins to register the insert SQLs of the plugins.
func RegisterInsertSQL(project string, insertSQL string) {
	sqlCollection.RegisterInsert(project, insertSQL)
}

// RegisterPreupgradeSQL registers a pre-upgrade SQL string to the given project.
// The pre-upgrade SQL string will be executed when pre-upgrading the project.
// The project name must be unique, the version string must be a valid version string, and the pre-upgrade SQL string must be a valid SQL string.
// The pre-upgrade SQL string will be appended to the existing pre-upgrade SQLs of the project.
// If the project does not exist, it will be created.
// If the project already exists, the pre-upgrade SQL string will be appended to the existing pre-upgrade SQLs of the project.
// The RegisterPreupgradeSQL function is used by the system to register the pre-upgrade SQLs of the system.
// The RegisterPreupgradeSQL function is also used by the plugins to register the pre-upgrade SQLs of the plugins.
func RegisterPreupgradeSQL(project string, version, preupgradeSQL string) {
	sqlCollection.RegisterPreupgrade(project, version, preupgradeSQL)
}

// GetInsertSQLs returns the map of insert SQLs.
// The key of the returned map is the project name, and the value is a slice of insert SQLs.
// The insert SQLs are registered using the RegisterInsertSQL function.
// The insert SQLs are executed when inserting data into the project.
// The insert SQLs are used by the system to initialize the project database.
// The insert SQLs are also used by the plugins to initialize the project database.
// The insert SQLs are appended to the existing insert SQLs of the project.
// If the project does not exist, it will be created.
func GetInsertSQLs() map[string][]string {
	return sqlCollection.Insert
}

// GetInstallSQLs returns the map of install SQLs.
// The key of the returned map is the project name, and the value is a slice of install SQLs.
func GetInstallSQLs() map[string][]string {
	return sqlCollection.Install
}

// GetPreupgradeSQLs returns the map of pre-upgrade SQLs.
// The key of the returned map is the project name, and the value is a map of version strings to pre-upgrade SQLs.
//
//	For example, if the returned map is map[string]map[string][]string{
//		"nging": {
//			"1.0.0": []string{"sql1", "sql2"},
//			"2.0.0": []string{"sql3", "sql4"},
//		},
//	}, it means that when upgrading the "nging" project from version "1.0.0" to "2.0.0", the system will execute the pre-upgrade SQL "sql3" and "sql4".
func GetPreupgradeSQLs() map[string]map[string][]string {
	return sqlCollection.Preupgrade
}

// genInstalledLockFileContent  generates a content string for the installed lock file.
//
// The generated string consists of three lines:
// - The first line is the current time in the format of `2006-01-02 15:04:05`.
// - The second line is the DB schema version.
// - The third line is the JSON string of the package DB schema version map.
//
// If there is an error when generating the JSON string, the function will return the error.
//
// Parameters:
// - now: The current time.
// - verInfo: The version information.
//
// Returns:
// - A string containing the content of the installed lock file.
// - An error if there is an error when generating the JSON string.
func genInstalledLockFileContent(now time.Time, verInfo *VersionInfo) (string, error) {
	content := now.Format(`2006-01-02 15:04:05`) + "\n" + com.String(verInfo.DBSchema)
	jsonV, err := com.JSONEncodeToString(verInfo.PkgDBSchemas)
	if err != nil {
		return content, err
	}
	content += "\n" + jsonV
	return content, nil
}

// parseInstalledLockFileContent parses the content of the installed lock file into the installed time, the DB schema version, and the package DB schema version map.
//
// The function takes a string content as the parameter, and returns the parsed installed time, the DB schema version, the package DB schema version map, and an error if there is an error when parsing the content.
//
// The content is expected to have three lines:
// - The first line is the current time in the format of `2006-01-02 15:04:05`.
// - The second line is the DB schema version.
// - The third line is the JSON string of the package DB schema version map.
//
// If there is an error when parsing the content, the function will return the error.
func parseInstalledLockFileContent(content string) (time.Time, float64, map[string]float64, error) {
	content = strings.TrimSpace(content)
	var installedTime, installedSchemaVer, installedPkgSchemaVer string
	lines := strings.Split(content, "\n")
	com.SliceExtract(lines, &installedTime, &installedSchemaVer, &installedPkgSchemaVer)
	var t time.Time
	var schemaV float64
	var pkgSchemaV = map[string]float64{}
	var err error
	if len(installedSchemaVer) > 0 {
		schemaV, _ = strconv.ParseFloat(strings.TrimSpace(installedSchemaVer), 64)
	}
	if len(installedPkgSchemaVer) > 0 {
		err = com.JSONDecodeString(installedPkgSchemaVer, &pkgSchemaV)
	}
	if len(installedTime) > 0 {
		t, _ = time.Parse(`2006-01-02 15:04:05`, strings.TrimSpace(installedTime))
	}
	return t, schemaV, pkgSchemaV, err
}

// SetInstalled writes the current time, the DB schema version, and the package DB schema version map to a file.
//
// The function takes a file path as the parameter, and writes the content to the file.
//
// If there is an error when writing the file, the function will return the error.
//
// The function also sets the installed time, the DB schema version, and the installed state to true.
func SetInstalled(lockFile string) error {
	now := time.Now()
	content, err := genInstalledLockFileContent(now, Version)
	if err != nil {
		return err
	}
	err = os.WriteFile(lockFile, []byte(content), os.ModePerm)
	if err != nil {
		return err
	}
	installedTime = now
	installedSchemaVer = Version.DBSchema
	installedPkgSchemaVer = Version.PkgDBSchemas
	Installed.Valid = true
	Installed.Bool = true
	return nil
}

// InstalledLockFile returns the path of the installed lock file if it exists and is valid.
// It checks the lock file in the configuration directory and the current working directory.
// If it does not find a valid lock file, it returns an empty string.
// The function is used to check if the installed lock file exists and is valid.
func InstalledLockFile() string {
	for _, lockFile := range []string{
		filepath.Join(FromCLI().ConfDir(), LockFileName),
		filepath.Join(echo.Wd(), LockFileName),
	} {
		if info, err := os.Stat(lockFile); err == nil && !info.IsDir() {
			return lockFile
		}
	}
	return ``
}

// IsInstalled returns whether the installed lock file exists and is valid.
//
// It first checks the Installed field to see if it is valid. If it is not valid,
// it tries to read the installed lock file. If the file does not exist or is invalid,
// it sets the Installed field to false and returns false.
//
// Otherwise, it sets the Installed field to true and returns true.
func IsInstalled() bool {
	if Installed.Valid {
		return Installed.Bool
	}
	lockFile := InstalledLockFile()
	if len(lockFile) > 0 {
		if b, e := os.ReadFile(lockFile); e == nil && len(b) > 0 {
			content := string(b)
			installedTime, installedSchemaVer, installedPkgSchemaVer, e = parseInstalledLockFileContent(content)
			if e != nil {
				log.Error(e)
				Installed.Valid = true
				Installed.Bool = false
				return Installed.Bool
			}
		}
		Installed.Valid = true
		Installed.Bool = true
	}
	return Installed.Bool
}

// GetSQLInstallFiles returns a list of SQL files to be executed when installing the project.
// The function first checks if the `install.sql` file exists in the configuration directory.
// If it does, it adds the file to the list.
// Then it uses the `filepath.Glob` function to find all files in the configuration directory that match the pattern `install.*.sql`.
// If any matches are found, it adds them to the list.
// Finally, it returns the list of files and any error that occurred.
func GetSQLInstallFiles() ([]string, error) {
	confDIR := FromCLI().Confd
	sqlFile := filepath.Join(confDIR, `install.sql`)
	var sqlFiles []string
	if com.FileExists(sqlFile) {
		sqlFiles = append(sqlFiles, sqlFile)
	}
	matches, err := filepath.Glob(confDIR + echo.FilePathSeparator + `install.*.sql`)
	if len(matches) > 0 {
		sqlFiles = append(sqlFiles, matches...)
	}
	return sqlFiles, err
}

// GetPreupgradeSQLFiles returns a list of SQL files to be executed when pre-upgrading the project.
// The function first checks if the `preupgrade.*.sql` files exist in the configuration directory.
// If they do, it adds them to the list.
// Finally, it returns the list of files and any error that occurred.
func GetPreupgradeSQLFiles() []string {
	confDIR := FromCLI().Confd
	sqlFiles := []string{}
	matches, _ := filepath.Glob(confDIR + echo.FilePathSeparator + `preupgrade.*.sql`)
	if len(matches) > 0 {
		sqlFiles = append(sqlFiles, matches...)
	}
	return sqlFiles
}

// GetSQLInsertFiles returns a list of SQL files to be executed when inserting data into the project.
// The function first checks if the `insert.sql` file exists in the configuration directory.
// If it does, it adds the file to the list.
// Then it uses the `filepath.Glob` function to find all files in the configuration directory that match the pattern `insert.*.sql`.
// If any matches are found, it adds them to the list.
// Finally, it returns the list of files and any error that occurred.
func GetSQLInsertFiles() []string {
	confDIR := FromCLI().Confd
	sqlFile := filepath.Join(confDIR, `insert.sql`)
	sqlFiles := []string{}
	if com.FileExists(sqlFile) {
		sqlFiles = append(sqlFiles, sqlFile)
	}
	matches, _ := filepath.Glob(confDIR + echo.FilePathSeparator + `insert.*.sql`)
	if len(matches) > 0 {
		sqlFiles = append(sqlFiles, matches...)
	}
	return sqlFiles
}
