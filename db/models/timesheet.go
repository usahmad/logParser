package models

import (
	"gorm.io/gorm"
)

func (TimeSheet) TableName() string {
	return "timesheets"
}

type TimeSheet struct {
	SIP        int    `json:"sip" gorm:"column:sip" form:"sip"`
	TimeWorked int    `json:"time_worked" gorm:"column:time_worked" form:"time_worked"`
	Date       string `json:"date" gorm:"column:date" form:"date"`
}

func (m *TimeSheet) Create(db *gorm.DB, ivrs TimeSheet) (err error) {
	err = db.Create(&ivrs).Error
	if err != nil {
		return err
	}
	return nil
}
