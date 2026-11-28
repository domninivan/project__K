package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Date struct {
	time time.Time
}

const dateLayout = "2006-01-02"

func NewDate(y, m, d int) (Date, error) {
	if time.Month(m) < time.January || time.Month(m) > time.December {
		return Date{}, ErrInvalidDate
	}
	if d < 1 || d > 31 {
		return Date{}, ErrInvalidDate
	}
	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	if t.Year() != y || t.Month() != time.Month(m) || t.Day() != d {
		return Date{}, ErrInvalidDate
	}
	return Date{time: t}, nil
}

func ParseDate(s string) (Date, error) {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return Date{}, fmt.Errorf("parse date %q: %w", s, err)
	}
	return Date{t}, nil
}

func (d Date) String() string {
	return d.time.Format(dateLayout)
}

func (d Date) Time() time.Time { return d.time }

func (d Date) IsZero() bool { return d.time.IsZero() }

// Date -> JSON
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.time.Format(dateLayout))
}

// JSON -> Date
func (d *Date) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		*d = Date{}
		return nil
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return fmt.Errorf("unmarshal date %q: %w", s, err)
	}
	*d = Date{t}
	return nil
}

// Date -> BD
func (d Date) Value() (driver.Value, error) {
	if d.time.IsZero() {
		return nil, nil
	}
	return d.time.Format(dateLayout), nil
}

// BD -> Date
func (d *Date) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*d = Date{}
		return nil
	case string:
		return d.scanString(v)
	case []byte:
		return d.scanString(string(v))
	case time.Time:
		*d = Date{time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.UTC)}
		return nil
	default:
		return fmt.Errorf("scan date: unsupported type %T", src)
	}
}

func (d *Date) scanString(s string) error {
	parsed, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

type ClockTime struct {
	time time.Duration
}

func ParseClockTime(s string) (ClockTime, error) {
	if parts := strings.Split(s, ":"); len(parts) == 2 || len(parts) == 3 {
		h, err := strconv.Atoi(parts[0])
		if err != nil {
			return ClockTime{}, fmt.Errorf("parse clocktime %q: %w", s, err)
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return ClockTime{}, fmt.Errorf("parse clocktime %q: %w", s, err)
		}
		sec := 0
		if len(parts) == 3 {
			sec, err = strconv.Atoi(parts[2])
			if err != nil {
				return ClockTime{}, fmt.Errorf("parse clocktime %q: %w", s, err)
			}
		}
		return NewClockTime(h, m, sec)
	}
	t, err := time.ParseDuration(s)
	if err != nil || t < 0 || t >= 24*time.Hour {
		return ClockTime{}, fmt.Errorf("parse clocktime %q: %w", s, err)
	}
	return ClockTime{t}, nil
}

func NewClockTime(h, m, sec int) (ClockTime, error) {
	if h < 0 || h > 23 || m < 0 || m > 59 || sec < 0 || sec > 59 {
		return ClockTime{}, ErrInvalidDate
	}
	return ClockTime{time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(sec)*time.Second}, nil
}

func (c ClockTime) String() string {
	total := int(c.time / time.Second)
	return fmt.Sprintf("%02d:%02d:%02d", total/3600, (total%3600)/60, total%60)
}

func (c ClockTime) Minutes() int { return int(c.time / time.Minute) }

func (c ClockTime) IsZero() bool { return c.time == 0 }

// ClockTime -> JSON
func (c ClockTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// JSON -> ClockTime
func (c *ClockTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		*c = ClockTime{}
		return nil
	}
	ct, err := ParseClockTime(s)
	if err != nil {
		return err
	}
	*c = ct
	return nil
}

// ClockTime -> BD
func (c ClockTime) Value() (driver.Value, error) {
	return c.String(), nil
}

// BD -> ClockTime
func (c *ClockTime) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*c = ClockTime{}
		return nil
	case string:
		return c.scanString(v)
	case []byte:
		return c.scanString(string(v))
	case time.Time:
		*c = ClockTime{time.Duration(v.Hour())*time.Hour + time.Duration(v.Minute())*time.Minute + time.Duration(v.Second())*time.Second}
		return nil
	default:
		return fmt.Errorf("scan clocktime: unsupported type %T", src)
	}
}

func (c *ClockTime) scanString(s string) error {
	parsed, err := ParseClockTime(s)
	if err != nil {
		return err
	}
	*c = parsed
	return nil
}
