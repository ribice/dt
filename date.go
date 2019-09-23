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

// A Date represents a date (year, month, day).
//
// This type does not include location information, and therefore does not
// describe a unique 24-hour timespan.
type Date struct {
	Year  int        // Year (e.g., 2014).
	Month time.Month // Month of the year (January = 1, ...).
	Day   int        // Day of the month, starting at 1.
	Valid bool
}

// DateOf returns the Date in which a time occurs in that time's location.
func DateOf(t time.Time) Date {
	d := Date{Valid: !t.IsZero()}
	d.Year, d.Month, d.Day = t.Date()
	return d
}

// ParseDate parses a string in RFC3339 full-date format and returns the date value it represents.
func ParseDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, err
	}
	return DateOf(t), nil
}

// String returns the date in RFC3339 full-date format.
func (d Date) String() string {
	if d.Valid {
		return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
	}
	return ""
}

// In returns the time corresponding to time 00:00:00 of the date in the location.
//
// In is always consistent with time.Date, even when time.Date returns a time
// on a different day. For example, if loc is America/Indiana/Vincennes, then both
//     time.Date(1955, time.May, 1, 0, 0, 0, 0, loc)
// and
//     dt.Date{Year: 1955, Month: time.May, Day: 1}.In(loc)
// return 23:00:00 on April 30, 1955.
//
// In panics if loc is nil.
func (d Date) In(loc *time.Location) time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, loc)
}

// AddDays returns the date that is n days in the future.
// n can also be negative to go into the past.
func (d Date) AddDays(n int) Date {
	return DateOf(d.In(time.UTC).AddDate(0, 0, n))
}

// DaysSince returns the signed number of days between the date and s, not including the end day.
// This is the inverse operation to AddDays.
func (d Date) DaysSince(s Date) (days int) {
	// We convert to Unix time so we do not have to worry about leap seconds:
	// Unix time increases by exactly 86400 seconds per day.
	deltaUnix := d.In(time.UTC).Unix() - s.In(time.UTC).Unix()
	return int(deltaUnix / 86400)
}

// Before reports whether d1 occurs before d2.
func (d Date) Before(d2 Date) bool {
	if d.Year != d2.Year {
		return d.Year < d2.Year
	}
	if d.Month != d2.Month {
		return d.Month < d2.Month
	}
	return d.Day < d2.Day
}

// ToTime converts Date to Go time
func (d Date) ToTime() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of d.String().
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be a string in a format accepted by ParseDate.
func (d *Date) UnmarshalText(data []byte) error {
	var err error
	*d, err = ParseDate(string(data))
	if err == nil {
		d.Valid = true
	}
	return err
}

// Value implements valuer interface
func (d Date) Value() (driver.Value, error) {
	if d.Valid {
		return driver.Value(d.String()), nil
	}
	return nil, nil
}

// Scan implements sql scannner interface
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = Date{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		dt, err := ParseDate(string(v))
		if err != nil {
			return err
		}
		*d = dt
		return nil
	case string:
		dt, err := ParseDate(v)
		if err != nil {
			return err
		}
		*d = dt
		return nil
	}
	return fmt.Errorf("Can't convert %T to Date", value)
}
