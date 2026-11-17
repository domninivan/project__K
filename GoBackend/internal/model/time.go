package model

import (
	"fmt"
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
	return "p"
}

// Date -> JSON
func (d Date) MarshalJSON() {

}

// JSON -> Date
func (d *Date) UnmarshalJSON() {

}

// Date -> BD
func (d *Date) Value() {

}

// BD -> Date
func (d *Date) Scan() {

}

type ClockTime struct {
	time time.Duration
}

func ParseClockTime(s string) (ClockTime, error) {
	t, err := time.ParseDuration(s)
	if err != nil {
		return ClockTime{}, fmt.Errorf("parse clocktime %q: %w", s, err)
	}
	return ClockTime{time: t}, nil
}

func NewClockTime(h, m int) (ClockTime, error) {
	return ClockTime{}, nil
}

func (c ClockTime) String() string {
	return "p"
}

// ClockTime -> JSON
func (c ClockTime) MarshalJSON() {

}

// JSON -> ClockTime
func (c *ClockTime) UnmarshalJSON() {

}

// ClockTime -> BD
func (c *ClockTime) Value() {

}

// BD -> ClockTime
func (c *ClockTime) Scan() {

}
