package dt

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDateOf(t *testing.T) {
	year, month, day := 2018, time.Month(12), 31
	td := time.Date(year, month, day, 23, 59, 59, 99, &time.Location{})
	date := DateOf(td)
	if !date.Valid {
		t.Error("expected date to be valid")
	}
	if date.Year != year {
		t.Errorf("expected year %v got %v", year, date.Year)
	}

	if date.Month != month {
		t.Errorf("expected month %v got %v", month, date.Month)
	}

	if date.Day != day {
		t.Errorf("expected day %v got %v", day, date.Day)
	}

}

func TestParseDate(t *testing.T) {
	cases := []struct {
		name    string
		str     string
		want    Date
		wantErr bool
	}{
		{
			name: "Valid date",
			str:  "2016-01-02",
			want: Date{2016, 1, 2, true},
		},
		{
			name: "Valid old date",
			str:  "0003-02-04",
			want: Date{3, 2, 4, true},
		},
		{
			name: "Invalid date",
			str:  "2016-01-02x",
			want: Date{},
		},
		{
			name: "Invalid month",
			str:  "2019-23-11",
			want: Date{},
		},
		{
			name: "Invalid day",
			str:  "2019-23-51",
			want: Date{},
		},
		{
			name: "Empty input",
			want: Date{},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.str)
			if got != tt.want {
				t.Errorf("ParseDate(%q) = %+v, want %+v", tt.str, got, tt.want)
			}
			if err != nil && tt.want != (Date{}) {
				t.Errorf("Unexpected error %v from ParseDate(%q)", err, tt.str)
			}
		})
	}

}

func TestString(t *testing.T) {
	cases := []struct {
		name string
		date Date
		want string
	}{
		{
			name: "Date not Valid",
			date: Date{},
			want: "",
		},
		{
			name: "Date set not Valid",
			date: Date{3, 2, 4, false},
			want: "",
		},
		{
			name: "Date valid",
			date: Date{2019, 12, 31, true},
			want: "2019-12-31",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.date.String()
			if tt.want != got {
				t.Errorf("expected %s got %s", tt.want, got)
			}
		})
	}
}

func TestIn(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	year, month, day := 2018, time.Month(12), 31
	td := Date{year, month, day, true}
	tt := td.In(loc)

	if l := tt.Location(); l != loc {
		t.Errorf("expected location %v got %v", loc, l)
	}

	tYear, tMonth, tDay := tt.Date()

	if tYear != year {
		t.Errorf("expected year %v got %v", year, tYear)
	}

	if tMonth != month {
		t.Errorf("expected month %v got %v", month, tMonth)
	}

	if tDay != day {
		t.Errorf("expected day %v got %v", day, tDay)
	}
}

func TestDateArithmetic(t *testing.T) {
	cases := []struct {
		name  string
		start Date
		end   Date
		days  int
	}{
		{
			name:  "zero days noop",
			start: Date{2014, 5, 9, true},
			end:   Date{2014, 5, 9, true},
			days:  0,
		},
		{
			name:  "crossing a year boundary",
			start: Date{2014, 12, 31, true},
			end:   Date{2015, 1, 1, true},
			days:  1,
		},
		{
			name:  "negative number of days",
			start: Date{2015, 1, 1, true},
			end:   Date{2014, 12, 31, true},
			days:  -1,
		},
		{
			name:  "full leap year",
			start: Date{2004, 1, 1, true},
			end:   Date{2005, 1, 1, true},
			days:  366,
		},
		{
			name:  "full non-leap year",
			start: Date{2001, 1, 1, true},
			end:   Date{2002, 1, 1, true},
			days:  365,
		},
		{
			name:  "crossing a leap second",
			start: Date{1972, 6, 30, true},
			end:   Date{1972, 7, 1, true},
			days:  1,
		},
		{
			name:  "dates before the unix epoch",
			start: Date{101, 1, 1, true},
			end:   Date{102, 1, 1, true},
			days:  365,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.start.AddDays(tt.days); got != tt.end {
				t.Errorf("[%s] %#v.AddDays(%v) = %#v, want %#v", tt.name, tt.start, tt.days, got, tt.end)
			}
			if got := tt.end.DaysSince(tt.start); got != tt.days {
				t.Errorf("[%s] %#v.Sub(%#v) = %v, want %v", tt.name, tt.end, tt.start, got, tt.days)
			}
		})
	}
}

func TestDateBefore(t *testing.T) {
	for _, tt := range []struct {
		d1, d2 Date
		want   bool
	}{
		{Date{2016, 12, 31, true}, Date{2017, 1, 1, true}, true},
		{Date{2016, 1, 1, true}, Date{2016, 3, 1, true}, true},
		{Date{2016, 12, 30, true}, Date{2016, 12, 31, true}, true},
	} {
		if got := tt.d1.Before(tt.d2); got != tt.want {
			t.Errorf("%v.Before(%v): got %t, want %t", tt.d1, tt.d2, got, tt.want)
		}
	}
}

func TestToGoTime(t *testing.T) {
	d := Date{2016, 12, 31, true}
	goTime := d.ToTime()
	if td := time.Date(2016, time.December, 31, 0, 0, 0, 0, time.UTC); td != goTime {
		t.Errorf("expected %v got %v", goTime, td)
	}
}

func TestMarshalDate(t *testing.T) {
	d := Date{2016, 12, 31, true}
	bts, err := json.Marshal(d)
	if err != nil {
		t.Errorf("expected success but got error: %v", err)
	}

	str := string(bts)

	if exp := `"` + d.String() + `"`; exp != str {
		t.Errorf("expected %s but got %s", exp, str)
	}
}

type dateReq struct {
	D Date `json:"date"`
}

func TestUnmarshalDate(t *testing.T) {
	cases := []struct {
		name string
		req  []byte
		want Date
	}{
		{
			name: "Invalid date",
			req:  []byte(`"date":"ABCD"`),
		},
		{
			name: "Valid date",
			req:  []byte(`{"date":"2019-11-04"}`),
			want: Date{2019, 11, 04, true},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var date dateReq
			err := json.Unmarshal(tt.req, &date)
			if (err == nil) != tt.want.Valid {
				t.Errorf("expected date to be empty, got err: %v", err)
			}
			if date.D != tt.want {
				t.Errorf("expected date %v, got %v", tt.want, date)
			}
		})
	}
}

func TestValueDate(t *testing.T) {
	cases := []struct {
		name      string
		req       Date
		wantValue bool
	}{{
		name:      "Valid date",
		req:       Date{2019, 12, 31, true},
		wantValue: true,
	},
		{
			name: "Invalid date",
			req:  Date{2019, 12, 31, false},
		}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			val, _ := tt.req.Value()
			if tt.wantValue != (val != nil) {
				t.Error("value returned different from expected")
			}
		})
	}
}

func TestScanDate(t *testing.T) {
	cases := []struct {
		name    string
		value   interface{}
		want    Date
		wantErr bool
	}{
		{
			name: "Nil value",
			want: Date{},
		},
		{
			name:  "Bytes value",
			value: []byte("2019-07-15"),
			want:  Date{2019, 07, 15, true},
		},
		{
			name:    "Bytes error",
			value:   []byte("2019-21-41"),
			wantErr: true,
		},
		{
			name:  "String value",
			value: "2019-07-15",
			want:  Date{2019, 07, 15, true},
		},
		{
			name:    "String error",
			value:   "2019-21-41",
			wantErr: true,
		},
		{
			name:    "Invalid type",
			value:   8,
			wantErr: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			d := &Date{}
			err := d.Scan(tt.value)
			if (err != nil) != tt.wantErr {
				t.Error("expected error and got error do not match")
			}
			if *d != tt.want {
				t.Errorf("expected %v, got %v", tt.want, *d)
			}
		})
	}
}
