package utils_test

import (
	"slices"
	"testing"
	"time"

	"biathlon-competitions-prototype/lib/utils"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		input    map[int]*utils.Result
		expected []int
	}{
		{
			name:     "empty map",
			input:    map[int]*utils.Result{},
			expected: []int{},
		},
		{
			name: "single element",
			input: map[int]*utils.Result{
				1: {Status: "[Finished]", TotalTime: 10 * time.Second},
			},
			expected: []int{1},
		},
		{
			name: "sorted by time descending",
			input: map[int]*utils.Result{
				1: {Status: "[Finished]", TotalTime: 30 * time.Second},
				2: {Status: "[Finished]", TotalTime: 20 * time.Second},
				3: {Status: "[Finished]", TotalTime: 10 * time.Second},
			},
			expected: []int{3, 2, 1},
		},
		{
			name: "NotStarted and NotFinished come last",
			input: map[int]*utils.Result{
				1: {Status: "[NotStarted]", TotalTime: 10 * time.Second},
				2: {Status: "[Finished]", TotalTime: 20 * time.Second},
				3: {Status: "[NotFinished]", TotalTime: 30 * time.Second},
				4: {Status: "[Finished]", TotalTime: 40 * time.Second},
			},
			expected: []int{1, 3, 2, 4},
		},
		{
			name: "mixed statuses",
			input: map[int]*utils.Result{
				1: {Status: "[Finished]", TotalTime: 10 * time.Second},
				2: {Status: "[NotStarted]", TotalTime: 20 * time.Second},
				3: {Status: "[Finished]", TotalTime: 30 * time.Second},
				4: {Status: "[NotFinished]", TotalTime: 40 * time.Second},
				5: {Status: "[Finished]", TotalTime: 50 * time.Second},
			},
			expected: []int{2, 4, 1, 3, 5},
		},
		{
			name: "equal times, different statuses",
			input: map[int]*utils.Result{
				1: {Status: "[Finished]", TotalTime: 10 * time.Second},
				2: {Status: "[NotStarted]", TotalTime: 10 * time.Second},
				3: {Status: "[Finished]", TotalTime: 10 * time.Second},
			},
			expected: []int{2, 1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := utils.Sort(tt.input)

			var keys []int
			for k := range order {
				keys = append(keys, k+1)
			}

			if len(keys) != len(tt.expected) {
				t.Errorf("expected %d elements, got %d", len(tt.expected), len(keys))
			}

			if !slices.Equal(order, tt.expected) {
				t.Errorf("got %v, want %v", order, tt.expected)
			}
		})
	}
}
