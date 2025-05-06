package utils_test

import (
	"testing"

	"biathlon-competitions-prototype/configs"
	"biathlon-competitions-prototype/lib/utils"
)

func TestCompetitorInitialization(t *testing.T) {
	competitor := &utils.Competitor{
		ID:              1,
		ShootingResults: make(map[int][]bool),
	}

	if competitor.ID != 1 {
		t.Errorf("Expected ID 1, got %d", competitor.ID)
	}

	if len(competitor.ShootingResults) != 0 {
		t.Error("ShootingResults should be empty")
	}
}

func TestResultInitialization(t *testing.T) {
	result := &utils.Result{
		CompetitorID: 1,
		Laps:         3,
	}

	if result.CompetitorID != 1 {
		t.Errorf("Expected CompetitorID 1, got %d", result.CompetitorID)
	}

	if result.Laps != 3 {
		t.Errorf("Expected 3 laps, got %d", result.Laps)
	}
}

func TestProcessEvents(t *testing.T) {
	cfg := &configs.Config{
		StartDelta:    "00:00:30.000",
		Laps:          2,
		LapLength:     4000,
		PenaltyLength: 150,
	}

	tests := []struct {
		name     string
		events   []utils.Event
		expected []string
	}{
		{
			name: "Registration and start",
			events: []utils.Event{
				{RawTime: "[10:00:00.000]", CompetitorID: 1, ID: 1},
				{RawTime: "[10:00:01.000]", CompetitorID: 1, ID: 2, ExtraParams: "10:00:30.000"},
				{RawTime: "[10:00:29.000]", CompetitorID: 1, ID: 3},
				{RawTime: "[10:00:30.000]", CompetitorID: 1, ID: 4},
			},
			expected: []string{
				"[10:00:00.000] The competitor(1) registered",
				"[10:00:01.000] The start time for the competitor(1) was set by a draw to 10:00:30.000",
				"[10:00:29.000] The competitor(1) is on the start line",
				"[10:00:30.000] The competitor(1) has started",
			},
		},
		{
			name: "Disqualification for late start",
			events: []utils.Event{
				{RawTime: "[10:00:00.000]", CompetitorID: 1, ID: 1},
				{RawTime: "[10:00:01.000]", CompetitorID: 1, ID: 2, ExtraParams: "10:00:30.000"},
				{RawTime: "[10:01:01.000]", CompetitorID: 1, ID: 3},
			},
			expected: []string{
				"[10:00:00.000] The competitor(1) registered",
				"[10:00:01.000] The start time for the competitor(1) was set by a draw to 10:00:30.000",
				"[10:01:01.000] The competitor(1) is disqualified",
			},
		},
		{
			name: "Complete race with two laps",
			events: []utils.Event{
				{RawTime: "[10:00:00.000]", CompetitorID: 1, ID: 1},
				{RawTime: "[10:00:01.000]", CompetitorID: 1, ID: 2, ExtraParams: "10:00:30.000"},
				{RawTime: "[10:00:29.000]", CompetitorID: 1, ID: 3},
				{RawTime: "[10:00:30.000]", CompetitorID: 1, ID: 4},
				{RawTime: "[10:02:30.000]", CompetitorID: 1, ID: 5, ExtraParams: "1"},
				{RawTime: "[10:02:31.000]", CompetitorID: 1, ID: 6, ExtraParams: "1"},
				{RawTime: "[10:02:32.000]", CompetitorID: 1, ID: 6, ExtraParams: "2"},
				{RawTime: "[10:02:40.000]", CompetitorID: 1, ID: 7},
				{RawTime: "[10:03:30.000]", CompetitorID: 1, ID: 10},
				{RawTime: "[10:05:30.000]", CompetitorID: 1, ID: 5, ExtraParams: "2"},
				{RawTime: "[10:05:31.000]", CompetitorID: 1, ID: 6, ExtraParams: "1"},
				{RawTime: "[10:05:40.000]", CompetitorID: 1, ID: 7},
				{RawTime: "[10:06:30.000]", CompetitorID: 1, ID: 10},
				{RawTime: "[10:06:31.000]", CompetitorID: 1, ID: 2, ExtraParams: "10:07:00.000"},
				{RawTime: "[10:06:51.000]", CompetitorID: 1, ID: 3},
				{RawTime: "[10:07:01.000]", CompetitorID: 1, ID: 4},
				{RawTime: "[10:08:21.000]", CompetitorID: 1, ID: 10},
			},
			expected: []string{
				"[10:00:00.000] The competitor(1) registered",
				"[10:00:01.000] The start time for the competitor(1) was set by a draw to 10:00:30.000",
				"[10:00:29.000] The competitor(1) is on the start line",
				"[10:00:30.000] The competitor(1) has started",
				"[10:02:30.000] The competitor(1) is on the firing range(1)",
				"[10:02:31.000] The target(1) has been hit by competitor(1)",
				"[10:02:32.000] The target(2) has been hit by competitor(1)",
				"[10:02:40.000] The competitor(1) left the firing range",
				"[10:03:30.000] The competitor(1) ended the main lap",
				"[10:05:30.000] The competitor(1) is on the firing range(2)",
				"[10:05:31.000] The target(1) has been hit by competitor(1)",
				"[10:05:40.000] The competitor(1) left the firing range",
				"[10:06:30.000] The competitor(1) ended the main lap",
				"[10:06:30.000] The competitor(1) has finished",
				"[10:06:31.000] The start time for the competitor(1) was set by a draw to 10:07:00.000",
				"[10:06:51.000] The competitor(1) is on the start line",
				"[10:07:01.000] The competitor(1) has started",
				"[10:08:21.000] The competitor(1) ended the main lap",
				"[10:08:21.000] The competitor(1) has finished",
			},
		},
		{
			name: "Complete race with shooting",
			events: []utils.Event{
				{RawTime: "[10:00:00.000]", CompetitorID: 1, ID: 1},
				{RawTime: "[10:00:01.000]", CompetitorID: 1, ID: 2, ExtraParams: "10:00:30.000"},
				{RawTime: "[10:00:29.000]", CompetitorID: 1, ID: 3},
				{RawTime: "[10:00:30.000]", CompetitorID: 1, ID: 4},
				{RawTime: "[10:02:30.000]", CompetitorID: 1, ID: 5, ExtraParams: "1"},
				{RawTime: "[10:02:31.000]", CompetitorID: 1, ID: 6, ExtraParams: "1"},
				{RawTime: "[10:02:32.000]", CompetitorID: 1, ID: 6, ExtraParams: "2"},
				{RawTime: "[10:02:40.000]", CompetitorID: 1, ID: 7},
				{RawTime: "[10:03:30.000]", CompetitorID: 1, ID: 10},
				{RawTime: "[10:05:30.000]", CompetitorID: 1, ID: 5, ExtraParams: "2"},
				{RawTime: "[10:05:31.000]", CompetitorID: 1, ID: 6, ExtraParams: "1"},
				{RawTime: "[10:05:40.000]", CompetitorID: 1, ID: 7},
				{RawTime: "[10:06:30.000]", CompetitorID: 1, ID: 10},
			},
			expected: []string{
				"[10:00:00.000] The competitor(1) registered",
				"[10:00:01.000] The start time for the competitor(1) was set by a draw to 10:00:30.000",
				"[10:00:29.000] The competitor(1) is on the start line",
				"[10:00:30.000] The competitor(1) has started",
				"[10:02:30.000] The competitor(1) is on the firing range(1)",
				"[10:02:31.000] The target(1) has been hit by competitor(1)",
				"[10:02:32.000] The target(2) has been hit by competitor(1)",
				"[10:02:40.000] The competitor(1) left the firing range",
				"[10:03:30.000] The competitor(1) ended the main lap",
				"[10:05:30.000] The competitor(1) is on the firing range(2)",
				"[10:05:31.000] The target(1) has been hit by competitor(1)",
				"[10:05:40.000] The competitor(1) left the firing range",
				"[10:06:30.000] The competitor(1) ended the main lap",
				"[10:06:30.000] The competitor(1) has finished",
			},
		},
		{
			name: "Race with penalty loops",
			events: []utils.Event{
				{RawTime: "[10:00:00.000]", CompetitorID: 1, ID: 1},
				{RawTime: "[10:00:01.000]", CompetitorID: 1, ID: 2, ExtraParams: "10:00:30.000"},
				{RawTime: "[10:00:29.000]", CompetitorID: 1, ID: 3},
				{RawTime: "[10:00:30.000]", CompetitorID: 1, ID: 4},
				{RawTime: "[10:02:30.000]", CompetitorID: 1, ID: 5, ExtraParams: "1"},
				{RawTime: "[10:02:40.000]", CompetitorID: 1, ID: 7},
				{RawTime: "[10:02:45.000]", CompetitorID: 1, ID: 8},
				{RawTime: "[10:02:55.000]", CompetitorID: 1, ID: 9},
				{RawTime: "[10:03:30.000]", CompetitorID: 1, ID: 10},
			},
			expected: []string{
				"[10:00:00.000] The competitor(1) registered",
				"[10:00:01.000] The start time for the competitor(1) was set by a draw to 10:00:30.000",
				"[10:00:29.000] The competitor(1) is on the start line",
				"[10:00:30.000] The competitor(1) has started",
				"[10:02:30.000] The competitor(1) is on the firing range(1)",
				"[10:02:40.000] The competitor(1) left the firing range",
				"[10:02:45.000] The competitor(1) entered the penalty laps",
				"[10:02:55.000] The competitor(1) left the penalty laps",
				"[10:03:30.000] The competitor(1) ended the main lap",
			},
		},
		{
			name: "Race not finished",
			events: []utils.Event{
				{RawTime: "[10:00:00.000]", CompetitorID: 1, ID: 1},
				{RawTime: "[10:00:01.000]", CompetitorID: 1, ID: 2, ExtraParams: "10:00:30.000"},
				{RawTime: "[10:00:29.000]", CompetitorID: 1, ID: 3},
				{RawTime: "[10:00:30.000]", CompetitorID: 1, ID: 4},
				{RawTime: "[10:02:30.000]", CompetitorID: 1, ID: 5, ExtraParams: "1"},
				{RawTime: "[10:02:40.000]", CompetitorID: 1, ID: 7},
				{RawTime: "[10:02:45.000]", CompetitorID: 1, ID: 8},
				{RawTime: "[10:02:55.000]", CompetitorID: 1, ID: 9},
				{RawTime: "[10:03:30.000]", CompetitorID: 1, ID: 10},
				{RawTime: "[10:03:30.000]", CompetitorID: 1, ID: 11, ExtraParams: "Lost in the forest"},
			},
			expected: []string{
				"[10:00:00.000] The competitor(1) registered",
				"[10:00:01.000] The start time for the competitor(1) was set by a draw to 10:00:30.000",
				"[10:00:29.000] The competitor(1) is on the start line",
				"[10:00:30.000] The competitor(1) has started",
				"[10:02:30.000] The competitor(1) is on the firing range(1)",
				"[10:02:40.000] The competitor(1) left the firing range",
				"[10:02:45.000] The competitor(1) entered the penalty laps",
				"[10:02:55.000] The competitor(1) left the penalty laps",
				"[10:03:30.000] The competitor(1) ended the main lap",
				"[10:03:30.000] The competitor(1) can`t continue: Lost in the forest",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputEvents, _, _ := utils.ProcessEvents(cfg, tt.events)

			if len(outputEvents) != len(tt.expected) {
				t.Errorf("Expected %d events, got %d", len(tt.expected), len(outputEvents))
			}

			for i := range outputEvents {
				if outputEvents[i] != tt.expected[i] {
					t.Errorf("Event %d mismatch:\nExpected: %s\nGot:      %s", i, tt.expected[i], outputEvents[i])
				}
			}
		})
	}
}

func TestFormatResult(t *testing.T) {
	tests := []struct {
		name     string
		result   *utils.Result
		expected string
	}{
		{
			name: "Completed race",
			result: &utils.Result{
				CompetitorID:  1,
				Status:        "10:30:15.500",
				Laps:          2,
				LapTimes:      []string{"03:45.123", "03:30.456"},
				AvgSpeeds:     []string{"15.123", "15.456"},
				PenaltyTimes:  []string{"00:30.123"},
				PenaltySpeeds: []string{"4.987"},
				ShootingStats: "8/10",
			},
			expected: "[10:30:15.500] 1 [{03:45.123, 15.123}, {03:30.456, 15.456}] [{00:30.123, 4.987}] 8/10",
		},
		{
			name: "Disqualified",
			result: &utils.Result{
				CompetitorID:  2,
				Status:        "NotStarted",
				Laps:          2,
				LapTimes:      []string{},
				AvgSpeeds:     []string{},
				ShootingStats: "0/0",
			},
			expected: "[NotStarted] 2 [{,}, {,}] [] 0/0",
		},
		{
			name: "Not finished",
			result: &utils.Result{
				CompetitorID:  3,
				Status:        "NotFinished",
				Laps:          2,
				LapTimes:      []string{"03:45.123"},
				AvgSpeeds:     []string{"15.123"},
				ShootingStats: "5/5",
			},
			expected: "[NotFinished] 3 [{03:45.123, 15.123}, {,}] [] 5/5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := utils.FormatResult(tt.result)
			if formatted != tt.expected {
				t.Errorf("Format mismatch:\nExpected: %s\nGot:      %s", tt.expected, formatted)
			}
		})
	}
}

func TestShootingStatsCalculation(t *testing.T) {
	cfg := &configs.Config{
		StartDelta:    "00:00:30.000",
		Laps:          2,
		LapLength:     4000,
		PenaltyLength: 150,
	}

	events := []utils.Event{
		{RawTime: "[10:00:00.000]", CompetitorID: 1, ID: 1},
		{RawTime: "[10:00:01.000]", CompetitorID: 1, ID: 2, ExtraParams: "10:00:30.000"},
		{RawTime: "[10:00:29.000]", CompetitorID: 1, ID: 3},
		{RawTime: "[10:00:30.000]", CompetitorID: 1, ID: 4},
		{RawTime: "[10:02:30.000]", CompetitorID: 1, ID: 5, ExtraParams: "1"},
		{RawTime: "[10:02:31.000]", CompetitorID: 1, ID: 6, ExtraParams: "1"},
		{RawTime: "[10:02:32.000]", CompetitorID: 1, ID: 6, ExtraParams: "2"},
		{RawTime: "[10:02:33.000]", CompetitorID: 1, ID: 6, ExtraParams: "3"},
		{RawTime: "[10:02:40.000]", CompetitorID: 1, ID: 7},
		{RawTime: "[10:03:30.000]", CompetitorID: 1, ID: 10},
		{RawTime: "[10:05:30.000]", CompetitorID: 1, ID: 5, ExtraParams: "2"},
		{RawTime: "[10:05:31.000]", CompetitorID: 1, ID: 6, ExtraParams: "1"},
		{RawTime: "[10:05:40.000]", CompetitorID: 1, ID: 7},
		{RawTime: "[10:06:30.000]", CompetitorID: 1, ID: 10},
	}

	_, results, _ := utils.ProcessEvents(cfg, events)
	result := results[1]

	if result.ShootingStats != "4/10" {
		t.Errorf("Expected shooting stats 4/10, got %s", result.ShootingStats)
	}
}
