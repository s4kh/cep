package main

import (
	"fmt"
	"testing"
)

func TestParseField(t *testing.T) {
	var tests = []struct {
		field      string
		expression string
		want       string
	}{
		{"minute", "8/16", "8 24 40 56"},
		{"minute", "*/15,33", "0 15 30 33 45"},
		{"minute", "1/10,20-23", "1 11 20 21 22 23 31 41 51"},
		// {"minute", "1/10,20-66", "got range error got 20-66 expected range 0-59"},
		{"hour", "0", "0"},
		{"dayOfMonth", "1,15", "1 15"},
		{"month", "*", "1 2 3 4 5 6 7 8 9 10 11 12"},
		{"dayOfWeek", "1-5", "1 2 3 4 5"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.field, tt.expression)
		t.Run(testname, func(t *testing.T) {
			ans, err := ParseField(tt.field, tt.expression)
			if err != nil {
				t.Errorf("got error %s", err)
			}
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}

func TestParseFieldError(t *testing.T) {
	var tests = []struct {
		field      string
		expression string
		wantErr    string
	}{
		{"minute", "1/10,20-66", "got range error got 20-66 expected range 0-59"},
		{"hour", "zxc", "invalid value"},
		{"dayOfMonth", "33", "max value 31"},
		{"month", "20-66", "got range error got 20-66 expected range 0-59"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.field, tt.expression)
		t.Run(testname, func(t *testing.T) {
			_, err := ParseField(tt.field, tt.expression)
			if err != nil {
				t.Logf("got error: %s", err)
			}
			if err == nil {
				t.Errorf("should throw error %s", tt.wantErr)
			}
		})
	}
}
