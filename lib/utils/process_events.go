package utils

import (
	"biathlon-competitions-prototype/configs"
	"fmt"
	"strings"
	"time"
)

type Competitor struct {
	ID                   int
	Registered           bool
	PlannedStart         time.Time
	ActualStart          time.Time
	PenaltyStart         time.Time
	IsFinishedCompletely bool
	IsDisqualified       bool
	IsNotFinished        bool
	FinishTime           time.Time
	LapTimes             []time.Duration
	PenaltyTimes         []time.Duration
	ShootingResults      map[int][]bool
	CurrentLap           int
	OnFiringRange        bool
	OnPenaltyLoop        bool
	Comment              string
}

type Result struct {
	CompetitorID  int
	Status        string
	Laps          int
	LapTimes      []string
	AvgSpeeds     []string
	PenaltyTimes  []string
	PenaltySpeeds []string
	ShootingStats string
	TotalTime     time.Duration
}

func ProcessEvents(cfg *configs.Config, events []Event) ([]string, map[int]*Result, []int) {
	competitors := make(map[int]*Competitor)
	results := make(map[int]*Result)
	var outputEvents []string

	startDelta, _ := ParseDuration(cfg.StartDelta, "15:04:05.000")

	for _, event := range events {
		eventTime, _ := ParseTime(event.RawTime[1 : len(event.RawTime)-1])
		competitor, exists := competitors[event.CompetitorID]
		if !exists {
			competitor = &Competitor{
				ID:              event.CompetitorID,
				ShootingResults: make(map[int][]bool),
			}
			competitors[event.CompetitorID] = competitor
		}

		if _, exists := results[event.CompetitorID]; !exists {
			results[event.CompetitorID] = &Result{CompetitorID: competitor.ID, Laps: cfg.Laps}
		}

		switch event.ID {
		case 1:
			competitor.Registered = true
			outputEvents = append(outputEvents, fmt.Sprintf("%s The competitor(%d) registered", event.RawTime, event.CompetitorID))
		case 2:
			plannedTime, _ := ParseTime(event.ExtraParams)
			competitor.PlannedStart = plannedTime
			outputEvents = append(outputEvents, fmt.Sprintf("%s The start time for the competitor(%d) was set by a draw to %s", event.RawTime, event.CompetitorID, event.ExtraParams))
		case 3:
			if eventTime.After(competitor.PlannedStart.Add(startDelta)) {
				competitor.IsDisqualified = true
				outputEvents = append(outputEvents, fmt.Sprintf("%s The competitor(%d) is disqualified", event.RawTime, event.CompetitorID))
				break
			}
			outputEvents = append(outputEvents, fmt.Sprintf("%s The competitor(%d) is on the start line", event.RawTime, event.CompetitorID))
		case 4:
			competitor.ActualStart = eventTime
			outputEvents = append(outputEvents, fmt.Sprintf("%s The competitor(%d) has started", event.RawTime, event.CompetitorID))
		case 5:
			competitor.OnFiringRange = true
			outputEvents = append(outputEvents, fmt.Sprintf("%s The competitor(%d) is on the firing range(%s)", event.RawTime, event.CompetitorID, event.ExtraParams))
		case 6:
			if competitor.ShootingResults[competitor.CurrentLap] == nil {
				competitor.ShootingResults[competitor.CurrentLap] = make([]bool, 0)
			}
			competitor.ShootingResults[competitor.CurrentLap] = append(competitor.ShootingResults[competitor.CurrentLap], true)
			outputEvents = append(outputEvents, fmt.Sprintf("%s The target(%s) has been hit by competitor(%d)", event.RawTime, event.ExtraParams, event.CompetitorID))
		case 7:
			competitor.OnFiringRange = false
			outputEvents = append(outputEvents, fmt.Sprintf("%s The competitor(%d) left the firing range", event.RawTime, event.CompetitorID))
		case 8:
			competitor.OnPenaltyLoop = true
			competitor.PenaltyStart = eventTime
			outputEvents = append(outputEvents, fmt.Sprintf("%s The competitor(%d) entered the penalty laps", event.RawTime, event.CompetitorID))
		case 9:
			competitor.OnPenaltyLoop = false
			penaltyTime := eventTime.Sub(competitor.PenaltyStart)
			competitor.PenaltyTimes = append(competitor.PenaltyTimes, penaltyTime)
			results[event.CompetitorID].PenaltyTimes = append(results[event.CompetitorID].PenaltyTimes, FormatDurationToTime(penaltyTime))
			speed := float64(cfg.PenaltyLength) / penaltyTime.Seconds()
			results[event.CompetitorID].PenaltySpeeds = append(results[event.CompetitorID].PenaltySpeeds, fmt.Sprintf("%.3f", speed))
			outputEvents = append(outputEvents, fmt.Sprintf("%s The competitor(%d) left the penalty laps", event.RawTime, event.CompetitorID))
		case 10:
			lapTime := eventTime.Sub(competitor.ActualStart)
			competitor.LapTimes = append(competitor.LapTimes, lapTime)
			results[event.CompetitorID].TotalTime = eventTime.Sub(competitor.ActualStart)
			results[event.CompetitorID].LapTimes = append(results[event.CompetitorID].LapTimes, FormatDurationToTime(lapTime))
			speed := float64(cfg.LapLength) / lapTime.Seconds()
			results[event.CompetitorID].AvgSpeeds = append(results[event.CompetitorID].AvgSpeeds, fmt.Sprintf("%.3f", speed))
			competitor.CurrentLap++
			outputEvents = append(outputEvents, fmt.Sprintf("%s The competitor(%d) ended the main lap", event.RawTime, event.CompetitorID))
			if competitor.CurrentLap >= cfg.Laps {
				competitor.IsFinishedCompletely = true
				competitor.FinishTime = eventTime
				outputEvents = append(outputEvents, fmt.Sprintf("%s The competitor(%d) has finished", event.RawTime, event.CompetitorID))
			}
		case 11:
			competitor.IsNotFinished = true
			competitor.Comment = event.ExtraParams
			outputEvents = append(outputEvents, fmt.Sprintf("%s The competitor(%d) can`t continue: %s", event.RawTime, event.CompetitorID, event.ExtraParams))
		}
	}

	for _, competitor := range competitors {
		result := results[competitor.ID]

		if competitor.IsDisqualified {
			result.Status = "[NotStarted]"
		} else if competitor.IsNotFinished {
			result.Status = "[NotFinished]"
		} else {
			result.Status = competitor.FinishTime.Format("15:04:05.000")
		}

		hits := 0
		shots := len(competitor.ShootingResults) * 5
		for _, lapHits := range competitor.ShootingResults {
			hits += len(lapHits)
		}
		result.ShootingStats = fmt.Sprintf("%d/%d", hits, shots)
	}

	order := Sort(results)

	return outputEvents, results, order
}

func FormatResult(result *Result) string {
	var builder strings.Builder

	builder.WriteString("[")
	builder.WriteString(result.Status)
	builder.WriteString("]")
	builder.WriteString(" ")
	builder.WriteString(fmt.Sprintf("%d", result.CompetitorID))

	builder.WriteString(" [")
	if result.Laps > len(result.LapTimes) {
		for i := 0; i < len(result.LapTimes); i++ {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(fmt.Sprintf("{%s, %s}", result.LapTimes[i], result.AvgSpeeds[i]))
		}
		for i := 0; i < result.Laps-len(result.LapTimes); i++ {
			if len(result.LapTimes) == 0 {
				builder.WriteString("{,}")
				if i < result.Laps-len(result.LapTimes)-1 {
					builder.WriteString(", ")
				}
				continue
			}
			builder.WriteString(", ")
			builder.WriteString("{,}")
		}
	} else {
		for i := 0; i < result.Laps; i++ {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(fmt.Sprintf("{%s, %s}", result.LapTimes[i], result.AvgSpeeds[i]))
		}
	}
	builder.WriteString("]")

	builder.WriteString(" [")
	for i := 0; i < len(result.PenaltyTimes); i++ {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("{%s, %s}", result.PenaltyTimes[i], result.PenaltySpeeds[i]))
	}
	builder.WriteString("]")

	builder.WriteString(" ")
	builder.WriteString(result.ShootingStats)

	return builder.String()
}
