package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		input time.Time
		expected string
	}{
		{
			name: "standard date - UTC",
			input: time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			expected: "17 Mar 2022 at 10:15",
		},
		{
			name: "leap year date",
			input: time.Date(2020, 2, 29, 23, 59, 0, 0, time.UTC),
			expected: "29 Feb 2020 at 23:59",
		},
		{
			name: "end of year date",
			input: time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: "31 Dec 2021 at 00:00",
		},
	
		{
			name: "non-UTC date",
			input: time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("EST", -5*60*60)),
			expected: "17 Mar 2022 at 15:15",
		},
		{
			name: "Empty date",
			input: time.Time{},
			expected: "",
		},
		{
			name: "CET date",
			input: time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			expected: "17 Mar 2022 at 09:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := humanDate(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q but got %q", tt.expected, result)
			}
		})
	}
}