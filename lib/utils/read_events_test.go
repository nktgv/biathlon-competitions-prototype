package utils_test

import (
	"biathlon-competitions-prototype/lib/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEvents(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []utils.Event
	}{
		{
			name:     "empty input",
			input:    "",
			expected: []utils.Event(nil),
		},
		{
			name:  "single basic event",
			input: "[10:00:00.000] 1 101",
			expected: []utils.Event{
				{
					RawTime:      "[10:00:00.000]",
					ID:           1,
					CompetitorID: 101,
					ExtraParams:  "",
				},
			},
		},
		{
			name:  "multiple events with extra params",
			input: "[10:00:00.000] 1 101\n[10:01:00.000] 2 102 10:05:00.000\n[10:02:00.000] 5 103 range1",
			expected: []utils.Event{
				{
					RawTime:      "[10:00:00.000]",
					ID:           1,
					CompetitorID: 101,
					ExtraParams:  "",
				},
				{
					RawTime:      "[10:01:00.000]",
					ID:           2,
					CompetitorID: 102,
					ExtraParams:  "10:05:00.000",
				},
				{
					RawTime:      "[10:02:00.000]",
					ID:           5,
					CompetitorID: 103,
					ExtraParams:  "range1",
				},
			},
		},
		{
			name:  "event with all fields",
			input: "[10:00:00.000] 6 101 target1",
			expected: []utils.Event{
				{
					RawTime:      "[10:00:00.000]",
					ID:           6,
					CompetitorID: 101,
					ExtraParams:  "target1",
				},
			},
		},
		{
			name:     "incomplete event with extra params",
			input:    "[10:00:00.000] 2 101",
			expected: []utils.Event(nil),
		},
		{
			name:     "event with missing competitor ID",
			input:    "[10:00:00.000] 1",
			expected: []utils.Event(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp(t.TempDir(), "testevents")
			require.NoError(t, err)
			defer os.Remove(tmpfile.Name())

			_, err = tmpfile.WriteString(tt.input)
			require.NoError(t, err)
			err = tmpfile.Close()
			require.NoError(t, err)

			events, err := utils.ReadEvents(tmpfile.Name())
			require.NoError(t, err)
			assert.Len(t, events, len(tt.expected))

			assert.Equal(t, tt.expected, events)
		})
	}
}

func TestReadEvents(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		expected    []utils.Event
		expectError bool
	}{
		{
			name:        "successful read",
			fileContent: "[10:00:00.000] 1 101\n[10:01:00.000] 2 102 10:05:00.000",
			expected: []utils.Event{
				{
					RawTime:      "[10:00:00.000]",
					ID:           1,
					CompetitorID: 101,
					ExtraParams:  "",
				},
				{
					RawTime:      "[10:01:00.000]",
					ID:           2,
					CompetitorID: 102,
					ExtraParams:  "10:05:00.000",
				},
			},
			expectError: false,
		},
		{
			name:        "file not found",
			fileContent: "",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			var err error

			if !tt.expectError {
				tmpfile, err := os.CreateTemp(t.TempDir(), "testevents")
				require.NoError(t, err)
				defer os.Remove(tmpfile.Name())

				_, err = tmpfile.WriteString(tt.fileContent)
				require.NoError(t, err)
				err = tmpfile.Close()
				require.NoError(t, err)

				path = tmpfile.Name()
			} else {
				path = "nonexistent_file.txt"
			}

			result, err := utils.ReadEvents(path)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseEventsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []utils.Event
	}{
		{
			name:  "malformed time format",
			input: "[bad_time] 1 101",
			expected: []utils.Event{
				{
					RawTime:      "[bad_time]",
					ID:           1,
					CompetitorID: 101,
					ExtraParams:  "",
				},
			},
		},
		{
			name:  "extra whitespace",
			input: "  [10:00:00.000]   1   101   \n  [10:01:00.000]   2   102   10:05:00.000  ",
			expected: []utils.Event{
				{
					RawTime:      "[10:00:00.000]",
					ID:           1,
					CompetitorID: 101,
					ExtraParams:  "",
				},
				{
					RawTime:      "[10:01:00.000]",
					ID:           2,
					CompetitorID: 102,
					ExtraParams:  "10:05:00.000",
				},
			},
		},
		{
			name:     "missing extra param when expected",
			input:    "[10:00:00.000] 6 101",
			expected: []utils.Event(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp(t.TempDir(), "testevents")
			require.NoError(t, err)
			defer os.Remove(tmpfile.Name())

			_, err = tmpfile.WriteString(tt.input)
			require.NoError(t, err)
			err = tmpfile.Close()
			require.NoError(t, err)

			events, err := utils.ReadEvents(tmpfile.Name())
			require.NoError(t, err)
			assert.Len(t, events, len(tt.expected))

			assert.Equal(t, tt.expected, events)
		})
	}
}
