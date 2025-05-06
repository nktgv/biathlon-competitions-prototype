package utils_test

import (
	"testing"
	"time"

	"biathlon-competitions-prototype/lib/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "valid time",
			input:    "12:34:56.789",
			expected: time.Date(0, 1, 1, 12, 34, 56, 789000000, time.UTC),
			wantErr:  false,
		},
		{
			name:     "time with leading zeros",
			input:    "01:02:03.004",
			expected: time.Date(0, 1, 1, 1, 2, 3, 4000000, time.UTC),
			wantErr:  false,
		},
		{
			name:     "midnight",
			input:    "00:00:00.000",
			expected: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:    "invalid format",
			input:   "12:34:56",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ParseTime(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		layout   string
		expected time.Duration
		wantErr  bool
	}{
		{
			name:     "standard time format",
			input:    "12:34:56.789",
			layout:   "15:04:05.000",
			expected: 12*time.Hour + 34*time.Minute + 56*time.Second + 789*time.Millisecond,
			wantErr:  false,
		},
		{
			name:     "time with leading zeros",
			input:    "01:02:03.004",
			layout:   "15:04:05.000",
			expected: 1*time.Hour + 2*time.Minute + 3*time.Second + 4*time.Millisecond,
			wantErr:  false,
		},
		{
			name:     "different layout",
			input:    "3:04PM",
			layout:   "3:04PM",
			expected: 15*time.Hour + 4*time.Minute,
			wantErr:  false,
		},
		{
			name:    "invalid format",
			input:   "12:34:56",
			layout:  "15:04:05.000",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			layout:  "15:04:05.000",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ParseDuration(tt.input, tt.layout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestFormatDurationToTime(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Duration
		expected string
	}{
		{
			name:     "standard duration",
			input:    12*time.Hour + 34*time.Minute + 56*time.Second + 789*time.Millisecond,
			expected: "12:34:56.789",
		},
		{
			name:     "duration with leading zeros",
			input:    1*time.Hour + 2*time.Minute + 3*time.Second + 4*time.Millisecond,
			expected: "01:02:03.004",
		},
		{
			name:     "zero duration",
			input:    0,
			expected: "00:00:00.000",
		},
		{
			name:     "duration less than second",
			input:    123 * time.Millisecond,
			expected: "00:00:00.123",
		},
		{
			name:     "duration less than minute",
			input:    59*time.Second + 999*time.Millisecond,
			expected: "00:00:59.999",
		},
		{
			name:     "duration less than hour",
			input:    59*time.Minute + 59*time.Second + 999*time.Millisecond,
			expected: "00:59:59.999",
		},
		{
			name:     "duration more than 24 hours",
			input:    25*time.Hour + 1*time.Minute + 1*time.Second + 1*time.Millisecond,
			expected: "25:01:01.001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.FormatDurationToTime(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
