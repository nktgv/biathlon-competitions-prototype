package main

import (
	"biathlon-competitions-prototype/configs"
	"biathlon-competitions-prototype/lib"
	"biathlon-competitions-prototype/lib/logger/sl"
	"biathlon-competitions-prototype/lib/utils"
	"fmt"
	"log/slog"
	"os"
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

	outputEvents, results := lib.ProcessEvents(cfg, events)

	outputLog, err := os.OpenFile("output.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Error("cannot open output file: ", sl.Err(err))
	}
	defer outputLog.Close()

	for _, event := range outputEvents {
		_, err := fmt.Fprintln(outputLog, event)
		if err != nil {
			log.Error("cannot write to output file: ", sl.Err(err))
		}
	}

	resultFile, err := os.OpenFile("result.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Error("cannot open output file: ", sl.Err(err))
	}
	defer resultFile.Close()

	for _, result := range results {
		_, err := fmt.Fprintln(resultFile, lib.FormatResult(result))
		if err != nil {
			log.Error("cannot write to output file: ", sl.Err(err))
		}
	}

	log.Info("finish processing events")
}
