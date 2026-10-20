package datetime

import (
	"testing"
	stdtime "time"
)

func TestDate(t *testing.T) {
	d, err := NewDate(2024, stdtime.February, 29)
	if err != nil {
		t.Fatal(err)
	}
	if d.Year() != 2024 || d.Month() != stdtime.February || d.Day() != 29 {
		t.Fatalf("date = %v", d)
	}
	if d.ISOFormat() != "2024-02-29" {
		t.Fatalf("ISOFormat = %q", d.ISOFormat())
	}
	if _, err := NewDate(2023, stdtime.February, 29); err == nil {
		t.Fatal("invalid leap date should error")
	}
	next := d.AddDays(1)
	if next.ISOFormat() != "2024-03-01" {
		t.Fatalf("AddDays = %v", next)
	}
	if next.Sub(d) != 1 {
		t.Fatalf("Sub = %d", next.Sub(d))
	}
	if !d.Before(next) || next.Before(d) {
		t.Fatal("Before mismatch")
	}
}

func TestDateTime(t *testing.T) {
	loc := stdtime.FixedZone("CST", 8*3600)
	base := stdtime.Date(2024, 1, 2, 3, 4, 5, 0, loc)
	d := DateTimeFromTime(base)
	if d.Timestamp() != base.Unix() {
		t.Fatalf("Timestamp = %d", d.Timestamp())
	}
	if d.ISOFormat() != "2024-01-02T03:04:05" {
		t.Fatalf("ISOFormat = %q", d.ISOFormat())
	}
	if d.Strftime("%Y/%m/%d %H:%M") != "2024/01/02 03:04" {
		t.Fatalf("Strftime = %q", d.Strftime("%Y/%m/%d %H:%M"))
	}
	if d.Weekday() != 2 {
		t.Fatalf("Weekday = %d", d.Weekday())
	}
	if d.Date().ISOFormat() != "2024-01-02" {
		t.Fatalf("Date = %v", d.Date())
	}
	next := d.AddDays(1)
	if next.Sub(d) != NewTimedelta(24*stdtime.Hour) {
		t.Fatalf("Sub = %v", next.Sub(d))
	}
}

func TestTimedelta(t *testing.T) {
	delta := NewTimedelta(36 * stdtime.Hour)
	if delta.Days() != 1.5 || delta.TotalSeconds() != 36*3600 {
		t.Fatalf("delta = %v", delta)
	}
	if delta.Neg().TotalSeconds() != -36*3600 {
		t.Fatal("Neg failed")
	}
	if delta.Neg().Abs() != delta {
		t.Fatal("Abs failed")
	}
	if delta.Add(NewTimedelta(stdtime.Hour)).TotalSeconds() != 37*3600 {
		t.Fatal("Add failed")
	}
}
