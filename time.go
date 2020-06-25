// Copyright 2016 Google LLC
// Copyright 2019 Emir Ribic
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dt

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// A Time represents a time with nanosecond precision.
//
// This type does not include location information, and therefore does not
// describe a unique moment in time.
//
// This type exists to represent the TIME type in storage-based APIs like BigQuery.
// Most operations on Times are unlikely to be meaningful. Prefer the DateTime type.
type Time struct {
	Hour   int // The hour of the day in 24-hour format; range [0-23]
	Minute int // The minute of the hour; range [0-59]
	Valid  bool
}

// TimeOf returns the Time representing the time of day in which a time occurs
// in that time's location. It ignores the date.
func TimeOf(t time.Time) Time {
	tm := Time{Valid: !t.IsZero()}
	tm.Hour, tm.Minute, _ = t.Clock()
	return tm
}

// ParseTime parses a string and returns the time value it represents.
// ParseTime accepts an extended form of the RFC3339 partial-time format. After
// the HH:MM:SS part of the string, an optional fractional part may appear,
// consisting of a decimal point followed by one to nine decimal digits.
// (RFC3339 admits only one digit after the decimal point).
func ParseTime(s string) (Time, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		t, err := time.Parse("15:04:05", s)
		return TimeOf(t), err
	}
	return TimeOf(t), nil
}

// String returns the date in the format described in ParseTime.
// If Valid is not true, it will return empty string
func (t Time) String() string {
	if t.Valid {
		return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
	}
	return ""
}

// ToDate converts Time into time.Time
func (t Time) ToDate() time.Time {
	return time.Date(0, 0, 0, t.Hour, t.Minute, 0, 0, time.UTC)
}

// After checks if instance of t is after tm
func (t Time) After(tm Time) bool {
	if t.Hour == tm.Hour{
		return t.Minute > tm.Minute
	}
	return t.Hour > tm.Hour
}

// After checks if instance of t is before tm
func (t Time) Before(tm Time) bool {
	if t.Hour == tm.Hour{
		return t.Minute < tm.Minute
	}
	return t.Hour < tm.Hour
}

// Subtract returns difference between t and t2 in minutes
func (t Time) Subtract(t2 Time) int {
	return (t.Hour-t2.Hour)*60 + t.Minute - t2.Minute
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of d.String().
func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be a string in a format accepted by ParseTime.
func (t *Time) UnmarshalText(data []byte) error {
	var err error
	*t, err = ParseTime(string(data))
	return err
}

// Value implements valuer interface
func (t Time) Value() (driver.Value, error) {
	if t.Valid {
		return driver.Value(t.String()), nil
	}
	return nil, nil
}

// Scan implements sql scan interface
func (t *Time) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		tm, err := ParseTime(string(v))
		if err != nil {
			return err
		}
		*t = tm
		return nil
	case string:
		tm, err := ParseTime(v)
		if err != nil {
			return err
		}
		*t = tm
		return nil
	}
	return fmt.Errorf("Can't convert %T to Time", value)
}
