package model

import "github.com/webx-top/echo"

const (
	// type
	CloudBackupTypeFull   = `full`
	CloudBackupTypeChange = `change`

	// status
	CloudBackupStatusSuccess = `success`
	CloudBackupStatusFailure = `failure`

	// operation
	CloudBackupOperationCreate = `create`
	CloudBackupOperationUpdate = `update`
	CloudBackupOperationDelete = `delete`
	CloudBackupOperationNone   = `none`
)

var CloudBackupTypes = echo.NewKVData().Add(CloudBackupTypeFull, echo.T(`全量`)).Add(CloudBackupTypeChange, echo.T(`监控`))
var CloudBackupStatuses = echo.NewKVData().Add(CloudBackupStatusSuccess, echo.T(`成功`)).Add(CloudBackupStatusFailure, echo.T(`失败`))
var CloudBackupOperations = echo.NewKVData().Add(CloudBackupOperationCreate, echo.T(`创建`)).
	Add(CloudBackupOperationUpdate, echo.T(`更新`)).
	Add(CloudBackupOperationDelete, echo.T(`删除`)).
	Add(CloudBackupOperationNone, echo.T(`未知`))
