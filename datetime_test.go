package dt

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDateTimeOf(t *testing.T) {
	cases := []struct {
		name string
		req  time.Time
		want DateTime
	}{
		{
			name: "Valid date and time",
			req:  time.Date(2014, 8, 20, 15, 8, 0, 0, time.Local),
			want: DateTime{Date{2014, 8, 20, true}, Time{15, 8, true}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := DateTimeOf(tt.req); got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestParseDateTime(t *testing.T) {
	cases := []struct {
		name    string
		req     string
		want    DateTime
		wantErr bool
	}{{
		name:    "Empty string",
		wantErr: true,
	},
		{
			name:    "Date only",
			req:     "2016-03-22",
			wantErr: true,
		},
		{
			name:    "wrong separating character",
			req:     "2016-03-22-13:26:33",
			wantErr: true,
		},
		{
			name:    "Extra char at end",
			req:     "2016-03-22T13:26:33x",
			wantErr: true,
		},
		{
			name: "Valid dateTime",
			req:  "2019-08-22T13:26:33",
			want: DateTime{
				Date: Date{2019, 8, 22, true},
				Time: Time{13, 26, true},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDateTime(tt.req)
			if (err != nil) != tt.wantErr {
				t.Error("expected and got error are different")
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestStringDateTime(t *testing.T) {
	cases := []struct {
		name string
		req  DateTime
		want string
	}{
		{
			name: "Valid dateTime",
			req:  DateTime{Time: Time{15, 35, true}, Date: Date{2019, 12, 31, true}},
			want: "2019-12-31T15:35",
		},
		{
			name: "Invalid dateTime",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.req.String(); got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestDateTimeIn(t *testing.T) {
	dt := DateTime{Date{2016, 1, 2, true}, Time{3, 4, true}}
	got := dt.In(time.UTC)
	want := time.Date(2016, 1, 2, 3, 4, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDateTimeBefore(t *testing.T) {
	d1 := Date{2016, 12, 31, true}
	d2 := Date{2017, 1, 1, true}
	t1 := Time{5, 6, true}
	t2 := Time{5, 7, true}
	for _, test := range []struct {
		dt1, dt2 DateTime
		want     bool
	}{
		{DateTime{d1, t1}, DateTime{d2, t1}, true},
		{DateTime{d1, t1}, DateTime{d1, t2}, true},
		{DateTime{d2, t1}, DateTime{d1, t1}, false},
		{DateTime{d2, t1}, DateTime{d2, t1}, false},
	} {
		if got := test.dt1.Before(test.dt2); got != test.want {
			t.Errorf("%v.Before(%v): got %t, want %t", test.dt1, test.dt2, got, test.want)
		}
	}
}

func TestMarshalDateTime(t *testing.T) {
	dt := DateTime{Date{2014, 8, 20, true}, Time{15, 8, true}}
	bts, err := json.Marshal(dt)
	if err != nil {
		t.Errorf("expected success but got error: %v", err)
	}

	str := string(bts)

	if exp := `"` + dt.String() + `"`; exp != str {
		t.Errorf("expected %s but got %s", exp, str)
	}
}

type dateTimeReq struct {
	DT DateTime `json:"date_time"`
}

func TestUnmarshalDateTime(t *testing.T) {
	cases := []struct {
		name    string
		req     []byte
		want    DateTime
		wantErr bool
	}{
		{
			name:    "Invalid date",
			req:     []byte(`"date_time":"ABCD"`),
			wantErr: true,
		},
		{
			name: "Valid date",
			req:  []byte(`{"date_time":"2019-11-04 15:35"}`),
			want: DateTime{Date{2019, 11, 4, true}, Time{15, 35, true}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var dateTime dateTimeReq
			err := json.Unmarshal(tt.req, &dateTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected and got error are different")
			}

			if dateTime.DT != tt.want {
				t.Errorf("expected dateTime %v, got %v", tt.want, dateTime)
			}
		})
	}
}

func TestValueDateTime(t *testing.T) {
	cases := []struct {
		name      string
		req       DateTime
		wantValue bool
	}{{
		name:      "Valid datetime",
		req:       DateTime{Time: Time{15, 35, true}, Date: Date{2019, 12, 31, true}},
		wantValue: true,
	},
		{
			name: "Invalid date",
			req:  DateTime{Time: Time{15, 35, true}},
		},
		{
			name: "Invalid time",
			req:  DateTime{Date: Date{2019, 12, 31, true}},
		},
		{
			name: "Invalid date and time",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			val, _ := tt.req.Value()
			if tt.wantValue != (val != nil) {
				t.Error("value returned different from expected")
			}
		})
	}
}

func TestScanDateTime(t *testing.T) {
	cases := []struct {
		name    string
		value   interface{}
		want    DateTime
		wantErr bool
	}{
		{
			name: "Nil value",
			want: DateTime{},
		},
		{
			name:  "Bytes value",
			value: []byte("2019-12-31 15:35"),
			want:  DateTime{Time: Time{15, 35, true}, Date: Date{2019, 12, 31, true}},
		},
		{
			name:    "Bytes error",
			value:   []byte("2019-21-41"),
			wantErr: true,
		},
		{
			name:  "String value",
			value: "2019-12-31 15:35",
			want:  DateTime{Time: Time{15, 35, true}, Date: Date{2019, 12, 31, true}},
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
			d := &DateTime{}
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
