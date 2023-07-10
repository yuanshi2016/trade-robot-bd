//package dto
/**
    * @Author: YuanShi
    * @Date: 2022/7/8 8:18 PM
    * @Desc: //TODO
**/
package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type TimeNormalYm time.Time    // 返回格式 2006-01
type TimeNormalYmd time.Time   // 返回格式 2006-01-02
type TimeNormalYmdH time.Time  // 返回格式 2006-01-02 10:00:00
type TimeNormalYmdHi time.Time // 返回格式 2006-01-02 10:01:00
type TimeNormalYmdHs time.Time // 返回格式 2006-01-02 10:01:01
const TimeFormatYmdHis = "2006-01-02 15:04:05"
const TimeFormatYmdHi = "2006-01-02 15:04"
const TimeFormatYmdH = "2006-01-02 15"
const TimeFormatYmd = "2006-01-02"
const TimeFormatYm = "2006-01-02"

type FormatDate int64

// GetNowYTime 根据传入时间戳 返回指定日期格式
func (m FormatDate) GetNowYTime() string {
	return time.Unix(int64(m), 0).Format("2006")
}

// GetNowYmTime 根据传入时间戳 返回指定日期格式
func (m FormatDate) GetNowYmTime() string {
	return time.Unix(int64(m), 0).Format("200601")
}

// GetNowYmdTime 根据传入时间戳 返回指定日期格式
func (m FormatDate) GetNowYmdTime() string {
	return time.Unix(int64(m), 0).Format("20060102")
}

// GetNowYmdHTime 根据传入时间戳 返回指定日期格式
func (m FormatDate) GetNowYmdHTime() string {
	return time.Unix(int64(m), 0).Format("2006010215")
}

// GetNowYmdhiTime 根据传入时间戳 返回指定日期格式
func (m FormatDate) GetNowYmdhiTime() string {
	return time.Unix(int64(m), 0).Format("200601021504")
}
func (t TimeNormalYmd) MarshalJSON() ([]byte, error) {
	ti := time.Time(t)
	tune := ti.Format(fmt.Sprintf(`"%s"`, TimeFormatYmd))
	return []byte(tune), nil
}

// Value insert timestamp into mysql need this function.
func (t TimeNormalYmd) Value() (driver.Value, error) {
	var zeroTime time.Time
	ti := time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}
func (t TimeNormalYm) MarshalJSON() ([]byte, error) {
	ti := time.Time(t)
	tune := ti.Format(fmt.Sprintf(`"%s"`, TimeFormatYm))
	return []byte(tune), nil
}

// Value insert timestamp into mysql need this function.
func (t TimeNormalYm) Value() (driver.Value, error) {
	var zeroTime time.Time
	ti := time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

// Scan valueof time.Time
func (t *TimeNormalYmd) Scan(v interface{}) error {
	ti, ok := v.(time.Time) // NOT directly assertion v.(TimeNormal)
	if ok {
		*t = TimeNormalYmd(ti)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t TimeNormalYmdH) MarshalJSON() ([]byte, error) {
	ti := time.Time(t)
	tune := ti.Format(fmt.Sprintf(`"%s"`, TimeFormatYmdH))
	return []byte(tune), nil
}

// Value insert timestamp into mysql need this function.
func (t TimeNormalYmdH) Value() (driver.Value, error) {
	var zeroTime time.Time
	ti := time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

// Scan valueof time.Time
func (t *TimeNormalYmdH) Scan(v interface{}) error {
	ti, ok := v.(time.Time) // NOT directly assertion v.(TimeNormal)
	if ok {
		*t = TimeNormalYmdH(ti)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t TimeNormalYmdHi) MarshalJSON() ([]byte, error) {
	ti := time.Time(t)
	tune := ti.Format(fmt.Sprintf(`"%s"`, TimeFormatYmdHi))
	return []byte(tune), nil
}

// Value insert timestamp into mysql need this function.
func (t TimeNormalYmdHi) Value() (driver.Value, error) {
	var zeroTime time.Time
	ti := time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

// Scan valueof time.Time
func (t *TimeNormalYmdHi) Scan(v interface{}) error {
	ti, ok := v.(time.Time) // NOT directly assertion v.(TimeNormal)
	if ok {
		*t = TimeNormalYmdHi(ti)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t TimeNormalYmdHs) MarshalJSON() ([]byte, error) {
	ti := time.Time(t)
	tune := ti.Format(fmt.Sprintf(`"%s"`, TimeFormatYmdHis))
	return []byte(tune), nil
}

// Value insert timestamp into mysql need this function.
func (t TimeNormalYmdHs) Value() (driver.Value, error) {
	var zeroTime time.Time
	ti := time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

// Scan valueof time.Time
func (t *TimeNormalYmdHs) Scan(v interface{}) error {
	ti, ok := v.(time.Time) // NOT directly assertion v.(TimeNormal)
	if ok {
		*t = TimeNormalYmdHs(ti)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
