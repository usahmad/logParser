package models

import (
	"gorm.io/gorm"
)

func (TimeSheetLog) TableName() string {
	return "timesheet_logs"
}

type TimeSheetLog struct {
	SIP  int    `json:"sip" gorm:"column:sip" form:"sip"`
	Date string `json:"date" gorm:"column:date" form:"date"`
	Type string `json:"type" gorm:"column:type" form:"type"`
}

func (m *TimeSheetLog) Create(db *gorm.DB, log TimeSheetLog) (err error) {
	err = db.Create(&log).Error
	if err != nil {
		return err
	}
	return nil
}
