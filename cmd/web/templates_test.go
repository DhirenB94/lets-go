package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	testData := []struct {
		testName           string
		testDate           time.Time
		expectedStringDate string
	}{
		{
			testName:           "UTC",
			testDate:           time.Date(2023, 07, 17, 14, 0, 0, 0, time.UTC),
			expectedStringDate: "17 Jul 2023 at 14:00",
		},
		{
			testName:           "empty time",
			testDate:           time.Time{},
			expectedStringDate: "",
		},
		{
			testName:           "CET",
			testDate:           time.Date(2023, 07, 17, 14, 0, 0, 0, time.FixedZone("CET",  1*60*60)),
			expectedStringDate: "17 Jul 2023 at 13:00",
		},
	}

	for _, tt := range testData {
		t.Run(tt.testName, func(t *testing.T) {
			actualStringDate := humanDate(tt.testDate)
			if actualStringDate != tt.expectedStringDate {
				t.Errorf("got: %q, want %q", actualStringDate, tt.expectedStringDate)
			}
		})
	}
}
