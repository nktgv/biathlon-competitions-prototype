package utils

import (
	"bufio"
	"fmt"
	"os"
)

type Event struct {
	ID           int
	RawTime      string
	CompetitorID int
	ExtraParams  string
}

func parseEvents(fileIn *os.File) ([]Event, error) {
	var in *bufio.Reader
	in = bufio.NewReader(fileIn)

	var rawEventTime, extraParams string
	var competitorID, eventID int
	var events []Event

	for {
		if _, err := fmt.Fscan(in, &rawEventTime, &eventID, &competitorID); err != nil {
			break
		}
		if eventID == 2 || eventID == 5 || eventID == 6 || eventID == 11 {
			_, err := fmt.Fscan(in, &extraParams)
			if err != nil {
				break
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
		return nil, fmt.Errorf("failed to open input file", err)
	}
	defer fileIn.Close()

	events, err := parseEvents(fileIn)
	if err != nil {
		return nil, fmt.Errorf("error in events parsing", err)
	}

	return events, nil
}
