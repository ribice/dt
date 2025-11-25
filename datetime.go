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

// A DateTime represents a date and time.
type DateTime struct {
	Date Date
	Time Time
}

// Note: We deliberately do not embed Date into DateTime, to avoid promoting AddDays and Sub.

// DateTimeOf returns the DateTime in which a time occurs in that time's location.
func DateTimeOf(t time.Time) DateTime {
	return DateTime{
		Date: DateOf(t),
		Time: TimeOf(t),
	}
}

var dtFormats = []string{"2006-01-02T15:04", "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02 15:04"}

// ParseDateTime parses a string and returns the DateTime it represents.
// ParseDateTime accepts a variant of the RFC3339 date-time format that omits
// the time offset but includes an optional fractional time, as described in
// ParseTime. Informally, the accepted format is
//
//	YYYY-MM-DDTHH:MM:SS[.FFFFFFFFF]
//
// where the 'T' may be a lower-case 't'.
func ParseDateTime(s string) (DateTime, error) {
	var t time.Time
	var err error
	for _, f := range dtFormats {
		t, err = time.Parse(f, s)
		if err == nil {
			break
		}
	}

	if err != nil {
		return DateTime{}, err
	}
	return DateTimeOf(t), nil
}

// String returns the date in the format described in ParseDate.
func (dt DateTime) String() string {
	if dt.Date.Valid && dt.Time.Valid {
		return dt.Date.String() + "T" + dt.Time.String()
	}
	return ""
}

// In returns the time corresponding to the DateTime in the given location.
//
// If the time is missing or ambigous at the location, In returns the same
// result as time.Date. For example, if loc is America/Indiana/Vincennes, then
// both
//
//	time.Date(1955, time.May, 1, 0, 30, 0, 0, loc)
//
// and
//
//	dt.DateTime{
//	    dt.Date{Year: 1955, Month: time.May, Day: 1}},
//	    dt.Time{Minute: 30}}.In(loc)
//
// return 23:30:00 on April 30, 1955.
//
// In panics if loc is nil.
func (dt DateTime) In(loc *time.Location) time.Time {
	return time.Date(dt.Date.Year, dt.Date.Month, dt.Date.Day, dt.Time.Hour, dt.Time.Minute, 0, 0, loc)
}

// Before reports whether dt occurs before dt2.
func (dt DateTime) Before(dt2 DateTime) bool {
	return dt.In(time.UTC).Before(dt2.In(time.UTC))
}

// After reports whether dt occurs after dt2.
func (dt DateTime) After(dt2 DateTime) bool {
	return dt2.Before(dt)
}

// Equal reports whether dt equals dt2.
func (dt DateTime) Equal(dt2 DateTime) bool {
	return dt.Date.Equal(dt2.Date) && dt.Time.Equal(dt2.Time)
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of dt.String().
func (dt DateTime) MarshalText() ([]byte, error) {
	return []byte(dt.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The datetime is expected to be a string in a format accepted by ParseDateTime
func (dt *DateTime) UnmarshalText(data []byte) error {
	var err error
	*dt, err = ParseDateTime(string(data))
	return err
}

// MarshalJSON implements the json.Marshaler interface.
// The output is a JSON string in YYYY-MM-DDTHH:MM format, or null if Date or Time is invalid.
func (dt DateTime) MarshalJSON() ([]byte, error) {
	if !dt.Date.Valid || !dt.Time.Valid {
		return []byte("null"), nil
	}
	return fmt.Appendf([]byte{}, "%q", dt.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The datetime is expected to be a JSON string in a format accepted by ParseDateTime,
// or null/empty string for an invalid datetime.
func (dt *DateTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" || string(data) == `""` {
		*dt = DateTime{}
		return
	}
	// Strip quotes from JSON string data
	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	*dt, err = ParseDateTime(s)
	if err == nil {
		dt.Date.Valid = true
		dt.Time.Valid = true
	}
	return
}

// Value implements valuer interface
func (dt DateTime) Value() (driver.Value, error) {
	if dt.Date.Valid && dt.Time.Valid {
		return driver.Value(dt.String()), nil
	}
	//nolint:nilnil
	return nil, nil
}

// Scan implements sql scan interface
func (dt *DateTime) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*dt = DateTimeOf(v)
		return nil
	case []byte:
		str := string(v)
		pdt, err := ParseDateTime(str)
		if err == nil {
			*dt = pdt
		}
		return err
	case string:
		pdt, err := ParseDateTime(v)
		if err == nil {
			*dt = pdt
		}
		return err
	default:
		return fmt.Errorf("can't convert %T to DateTime", value)
	}
}

// Compare compares dt and dt2. If dt is before dt2, it returns -1;
// if dt is after dt2, it returns +1; otherwise it returns 0.
func (dt DateTime) Compare(dt2 DateTime) int {
	return dt.In(time.UTC).Compare(dt2.In(time.UTC))
}
