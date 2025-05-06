package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

type Event struct {
	ID           int
	RawTime      string
	CompetitorID int
	ExtraParams  string
}

func parseEvents(fileIn *os.File) ([]Event, error) {
	in := bufio.NewReader(fileIn)

	var rawEventTime, extraParams string
	var competitorID, eventID int
	var events []Event

	for {
		if _, err := fmt.Fscan(in, &rawEventTime, &eventID, &competitorID); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("error while scanning: %w", err)
		}
		if eventID == 2 || eventID == 5 || eventID == 6 || eventID == 11 {
			_, err := fmt.Fscan(in, &extraParams)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				if extraParams != "" {
					return nil, fmt.Errorf("extra params missing: %w", err)
				}
				return nil, fmt.Errorf("error while scanning extra params: %w", err)
			}
		}
		events = append(events, Event{
			RawTime:      rawEventTime,
			ID:           eventID,
			CompetitorID: competitorID,
			ExtraParams:  extraParams,
		})
	}

	return events, nil
}

func ReadEvents(path string) ([]Event, error) {
	fileIn, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer func() { _ = fileIn.Close() }()

	events, err := parseEvents(fileIn)
	if err != nil {
		return nil, fmt.Errorf("cannot parse events: %w", err)
	}
	return events, nil
}
