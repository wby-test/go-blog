package model

import (
	"fmt"
	otgorm "github.com/eddycjy/opentracing-gorm"
	"github.com/fishblog/global"
	"github.com/fishblog/pkg/setting"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

const (
	STATE_CLOSE = iota
	STATE_OPEN
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

func NewDBEngine(dbSettings *setting.DatabaseSettings) (*gorm.DB, error) {
	db, err := gorm.Open(dbSettings.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?"+
		"charset=%s&parseTime=%t&loc=Local",
		dbSettings.UserName,
		dbSettings.Password,
		dbSettings.Host,
		dbSettings.DBName,
		dbSettings.Charset,
		dbSettings.ParseTime))
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}
	db.SingularTable(true)
	// 注册回调行为
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetMaxIdleConns(dbSettings.MaxIdleConns)
	db.DB().SetMaxOpenConns(dbSettings.MaxOpenConns)
	otgorm.AddGormCallbacks(db)

	return db, nil
}

// 新增行为的回调
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

// 更新行为的回调
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.Set("ModifiedOn", time.Now().Unix())
	}
}

// 删除行为的回调
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deleteOnField, hasDeleteOnField := scope.FieldByName("DeletedOn")
		isDelField, hasIsDelField := scope.FieldByName("IsDel")
		if !scope.Search.Unscoped && hasDeleteOnField && hasIsDelField {
			now := time.Now().Unix()
			scope.Raw(fmt.Sprintf(
				"UPDATE %v set %v=%v, %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deleteOnField.DBName),
				scope.AddToVars(now),
				scope.Quote(isDelField.DBName),
				scope.AddToVars(1),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption)))
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
