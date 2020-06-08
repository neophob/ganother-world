package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/neophob/ganother-world/logger"
)

const (
	framesToCapture int16 = 2000
)

type FrameData struct {
	Variables     [256]int16 `json:"variables"`
	ChannelPC     [64]uint16 `json:"channelPC"`
	ChannelPaused [64]bool   `json:"channelPaused"`
}

//DataExport holds the raw data
type DataExport struct {
	frames       [framesToCapture]FrameData
	currentFrame int16
	active       bool
}

type SerializedJsonDataFormat struct {
	Version         int                        `json:"version"`
	Filename        string                     `json:"filename"`
	FramesPerSecond int                        `json:"framesPerSecond"`
	Dump            [framesToCapture]FrameData `json:"dump"`
}

//Initialize Datastructure for data export
func InitDataExport() DataExport {
	return DataExport{
		active:       true,
		currentFrame: 0,
	}
}

func (dataExport *DataExport) addDataFrame(variables [256]int16, channelPC [64]uint16, channelPaused [64]bool) {
	if dataExport.active == false {
		return
	}

	frameData := FrameData{variables, channelPC, channelPaused}
	dataExport.frames[dataExport.currentFrame] = frameData
	dataExport.currentFrame++

	if dataExport.currentFrame%501 == 500 {
		logger.Info("dumped frame %d", dataExport.currentFrame)
	}

	if dataExport.currentFrame >= framesToCapture {
		dataExport.active = false

		logger.Info("Create JSON Export...")
		out := &SerializedJsonDataFormat{
			Version:         2,
			Filename:        "filename",
			FramesPerSecond: 25,
			Dump:            dataExport.frames,
		}
		b, err := json.Marshal(out)
		if err != nil {
			logger.Error("json.Marshal: %v", err)
			return
		}
		err = ioutil.WriteFile("./aw-dump.json", b, 0644)
		if err != nil {
			panic(err)
		}
		logger.Info("Export DONE!")
	}
}
