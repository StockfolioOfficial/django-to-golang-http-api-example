package supporter

import (
	"database/sql/driver"
	"errors"
	"time"
)


const timeLayout = "2006-01-02 15:04:05.999999Z07:00"

type Time time.Time

func (t *Time) Scan(value interface{}) error {
	if value == nil {
		return errors.New("") // todo error message
	}

	switch data := value.(type) {
	case string:
		src, err := time.Parse(timeLayout, data)
		*t = Time(src)
		return err
	case []byte:
		src, err := time.Parse(timeLayout, string(data))
		*t = Time(src)
		return err
	default:
		return errors.New("") // todo error message
	}
}

func (t Time) Value() (driver.Value, error) {
	return time.Time(t), nil
}


func (t Time) AsTime() time.Time {
	return time.Time(t)
}