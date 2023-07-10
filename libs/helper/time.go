package helper

import (
	"database/sql/driver"
	"fmt"
	"time"
)

func ForTime(s time.Duration, p func()) {
	for {
		time.Sleep(s)
		p()
	}
}

type TimeNormalHmd time.Time   // 返回格式 2006-01-02
type TimeNormalHmdH time.Time  // 返回格式 2006-01-02 10:00:00
type TimeNormalHmdHi time.Time // 返回格式 2006-01-02 10:01:00
type TimeNormalHmdHs time.Time // 返回格式 2006-01-02 10:01:01
const TimeFormatYmdHis = "2006-01-02 15:04:05"
const TimeFormatYmdHi = "2006-01-02 15:04"
const TimeFormatYmdH = "2006-01-02 15"
const TimeFormatYmd = "2006-01-02"

// GetNowTime 根据传入时间戳 返回指定日期格式
func GetNowTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("20060102")
}
func (t TimeNormalHmd) MarshalJSON() ([]byte, error) {
	ti := time.Time(t)
	tune := ti.Format(fmt.Sprintf(`"%s"`, TimeFormatYmd))
	return []byte(tune), nil
}

// Value insert timestamp into mysql need this function.
func (t TimeNormalHmd) Value() (driver.Value, error) {
	var zeroTime time.Time
	ti := time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

// Scan valueof time.Time
func (t *TimeNormalHmd) Scan(v interface{}) error {
	ti, ok := v.(time.Time) // NOT directly assertion v.(TimeNormal)
	if ok {
		*t = TimeNormalHmd(ti)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t TimeNormalHmdH) MarshalJSON() ([]byte, error) {
	ti := time.Time(t)
	tune := ti.Format(fmt.Sprintf(`"%s"`, TimeFormatYmdH))
	return []byte(tune), nil
}

// Value insert timestamp into mysql need this function.
func (t TimeNormalHmdH) Value() (driver.Value, error) {
	var zeroTime time.Time
	ti := time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

// Scan valueof time.Time
func (t *TimeNormalHmdH) Scan(v interface{}) error {
	ti, ok := v.(time.Time) // NOT directly assertion v.(TimeNormal)
	if ok {
		*t = TimeNormalHmdH(ti)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t TimeNormalHmdHi) MarshalJSON() ([]byte, error) {
	ti := time.Time(t)
	tune := ti.Format(fmt.Sprintf(`"%s"`, TimeFormatYmdHi))
	return []byte(tune), nil
}

// Value insert timestamp into mysql need this function.
func (t TimeNormalHmdHi) Value() (driver.Value, error) {
	var zeroTime time.Time
	ti := time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

// Scan valueof time.Time
func (t *TimeNormalHmdHi) Scan(v interface{}) error {
	ti, ok := v.(time.Time) // NOT directly assertion v.(TimeNormal)
	if ok {
		*t = TimeNormalHmdHi(ti)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t TimeNormalHmdHs) MarshalJSON() ([]byte, error) {
	ti := time.Time(t)
	tune := ti.Format(fmt.Sprintf(`"%s"`, TimeFormatYmdHis))
	return []byte(tune), nil
}

// Value insert timestamp into mysql need this function.
func (t TimeNormalHmdHs) Value() (driver.Value, error) {
	var zeroTime time.Time
	ti := time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

// Scan valueof time.Time
func (t *TimeNormalHmdHs) Scan(v interface{}) error {
	ti, ok := v.(time.Time) // NOT directly assertion v.(TimeNormal)
	if ok {
		*t = TimeNormalHmdHs(ti)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
