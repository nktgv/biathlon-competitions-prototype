package main

import (
	"fmt"
	"log/slog"
	"os"

	"biathlon-competitions-prototype/configs"
	"biathlon-competitions-prototype/lib/logger/sl"
	"biathlon-competitions-prototype/lib/utils"
)

func main() {
	cfg := configs.LoadConfig("./config.json")

	log := configs.ConfigureLogger()

	log.Info("config loaded", slog.Any("config", cfg))

	start, err := utils.ParseDuration(cfg.Start, "15:04:05.000")
	if err != nil {
		log.Error("cannot parse start time", sl.Err(err))
	}

	log.Info("start time", slog.Duration("start", start))

	startDelta, err := utils.ParseDuration(cfg.StartDelta, "15:04:05")
	if err != nil {
		log.Error("cannot parse start time: ", sl.Err(err))
	}

	log.Info("start delta", slog.Duration("startDelta", startDelta))

	filePath := "./events"
	events, err := utils.ReadEvents(filePath)
	if err != nil {
		log.Error("cannot read events: ", sl.Err(err))
	}

	outputEvents, results, order := utils.ProcessEvents(cfg, events)

	outputLog, err := os.OpenFile("output.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Error("cannot open output file: ", sl.Err(err))
	}
	defer func() { _ = outputLog.Close() }()

	for _, event := range outputEvents {
		_, err := fmt.Fprintln(outputLog, event)
		if err != nil {
			log.Error("cannot write to output file: ", sl.Err(err))
		}
	}

	resultFile, err := os.OpenFile("result.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Error("cannot open output file: ", sl.Err(err))
	}
	defer func() { _ = resultFile.Close() }()

	for _, i := range order {
		_, err := fmt.Fprintln(resultFile, utils.FormatResult(results[i]))
		if err != nil {
			log.Error("cannot write to output file: ", sl.Err(err))
		}
	}

	log.Info("finish processing events")
}
