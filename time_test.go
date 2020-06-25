package dt

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeOf(t *testing.T) {
	time := time.Date(2014, 8, 20, 15, 8, 43, 1, time.Local)
	want := Time{15, 8, true}
	if got := TimeOf(time); got != want {
		t.Errorf("TimeOf(%v) = %+v, want %+v", time, got, want)
	}
}

func TestTimeToString(t *testing.T) {
	cases := []struct {
		name    string
		req     string
		want    Time
		wantErr bool
	}{
		{
			name: "Hours and minutes",
			req:  "15:51",
			want: Time{15, 51, true},
		},
		{
			name: "Hours, minutes and seconds",
			req:  "16:52:33",
			want: Time{16, 52, true},
		},
		{
			name:    "Invalid time format",
			req:     "33:33:33",
			wantErr: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTime(tt.req)
			if (err != nil) != tt.wantErr {
				t.Error("expected error and got error do not match")
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}

}

func TestTimeString(t *testing.T) {
	cases := []struct {
		name string
		req  Time
		want string
	}{{
		name: "Valid date",
		req:  Time{15, 12, true},
		want: "15:12",
	},
		{
			name: "Invalid date",
			req:  Time{},
			want: "",
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

func TestTimeToDate(t *testing.T) {
	got := Time{15, 20, true}.ToDate()
	if want := time.Date(0, 0, 0, 15, 20, 0, 0, time.UTC); got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestAfter(t *testing.T) {
	cases := []struct {
		name    string
		t1, t2  Time
		isAfter bool
	}{
		{
			name:    "Hour after",
			t1:      Time{23, 59, true},
			t2:      Time{22, 59, true},
			isAfter: true,
		},
		{
			name:    "Minute after",
			t1:      Time{23, 59, true},
			t2:      Time{23, 45, true},
			isAfter: true,
		},
		{
			name:    "Equal",
			t1:      Time{23, 59, true},
			t2:      Time{23, 59, true},
			isAfter: false,
		},
		{
			name:    "Hour before",
			t1:      Time{22, 59, true},
			t2:      Time{23, 59, true},
			isAfter: false,
		},
		{
			name:    "Minute before",
			t1:      Time{11, 59, true},
			t2:      Time{23, 59, true},
			isAfter: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.t1.After(tt.t2) != tt.isAfter {
				t.Errorf("expected isAfter %v, got %v", tt.isAfter, tt.t1.After(tt.t2))
			}
		})
	}
}

func TestBefore(t *testing.T) {
	cases := []struct {
		name    string
		t1, t2  Time
		isBefore bool
	}{
		{
			name:    "Hour after",
			t1:      Time{23, 59, true},
			t2:      Time{22, 59, true},
			isBefore: false,
		},
		{
			name:    "Minute after",
			t1:      Time{23, 59, true},
			t2:      Time{23, 45, true},
			isBefore: false,
		},
		{
			name:    "Equal",
			t1:      Time{23, 59, true},
			t2:      Time{23, 59, true},
			isBefore: false,
		},
		{
			name:    "Hour before",
			t1:      Time{22, 59, true},
			t2:      Time{23, 59, true},
			isBefore: true,
		},
		{
			name:    "Minute before",
			t1:      Time{11, 59, true},
			t2:      Time{23, 59, true},
			isBefore: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.t1.Before(tt.t2) != tt.isBefore {
				t.Errorf("expected isBefore %v, got %v", tt.isBefore, tt.t1.Before(tt.t2))
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	cases := []struct {
		name   string
		t1, t2 Time
		diff   int
	}{
		{
			name: "Diff positive",
			t1:   Time{23, 59, true},
			t2:   Time{22, 59, true},
			diff: 60,
		},
		{
			name: "Diff negative",
			t1:   Time{23, 59, true},
			t2:   Time{23, 45, true},
			diff: 14,
		},
		{
			name: "No diff",
			t1:   Time{23, 59, true},
			t2:   Time{23, 59, true},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.t1.Subtract(tt.t2) != tt.diff {
				t.Errorf("expected diff %v, got %v", tt.diff, tt.t1.Subtract(tt.t2))
			}
		})
	}
}

func TestMarshalTime(t *testing.T) {
	tm := Time{15, 25, true}
	bts, err := json.Marshal(tm)
	if err != nil {
		t.Errorf("expected success but got error: %v", err)
	}

	str := string(bts)

	if exp := `"` + tm.String() + `"`; exp != str {
		t.Errorf("expected %s but got %s", exp, str)
	}
}

type timeReq struct {
	T Time `json:"time"`
}

func TestUnmarshalTime(t *testing.T) {
	cases := []struct {
		name string
		req  []byte
		want Time
	}{
		{
			name: "Invalid date",
			req:  []byte(`"date":"ABCD"`),
		},
		{
			name: "Valid date",
			req:  []byte(`{"time":"15:25"}`),
			want: Time{15, 25, true},
		},
		{
			name: "Valid date with sevonds",
			req:  []byte(`{"time":"15:25:35"}`),
			want: Time{15, 25, true},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var tm timeReq
			err := json.Unmarshal(tt.req, &tm)
			if (err == nil) != tt.want.Valid {
				t.Errorf("expected time to be empty, got err: %v", err)
			}
			if tm.T != tt.want {
				t.Errorf("expected time %v, got %v", tt.want, tm)
			}
		})
	}
}

func TestValueTime(t *testing.T) {
	cases := []struct {
		name      string
		req       Time
		wantValue bool
	}{{
		name:      "Valid time",
		req:       Time{15, 25, true},
		wantValue: true,
	},
		{
			name: "Invalid time",
			req:  Time{12, 91, false},
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

func TestScanTime(t *testing.T) {
	cases := []struct {
		name    string
		value   interface{}
		want    Time
		wantErr bool
	}{
		{
			name: "Nil value",
			want: Time{},
		},
		{
			name:  "Bytes value",
			value: []byte("15:25"),
			want:  Time{15, 25, true},
		},
		{
			name:    "Bytes error",
			value:   []byte("15:91"),
			wantErr: true,
		},
		{
			name:  "String value",
			value: "15:41",
			want:  Time{15, 41, true},
		},
		{
			name:    "String error",
			value:   "91:12",
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
			tm := &Time{}
			err := tm.Scan(tt.value)
			if (err != nil) != tt.wantErr {
				t.Error("expected error and got error do not match")
			}
			if *tm != tt.want {
				t.Errorf("expected %v, got %v", tt.want, *tm)
			}
		})
	}
}
