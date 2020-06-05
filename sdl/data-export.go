package main

import (
	"encoding/json"
	"os"

	"github.com/neophob/ganother-world/logger"
)

const (
	framesToCapture int16 = 10
)

type FrameData struct {
	Variables [256]int16
	ChannelPC [64]uint16
}

//DataExport holds the raw data
type DataExport struct {
	frames       [framesToCapture]FrameData
	currentFrame int16
	active       bool
}

type SerializedJsonDataFormat struct {
	version         uint8
	filename        string
	framesPerSecond uint8
	dump            [framesToCapture]FrameData
}

//Initialize Datastructure for data export
func InitDataExport() DataExport {
	return DataExport{active: true, currentFrame: 0}
}

func (dataExport *DataExport) addDataFrame(variables [256]int16, channelPC [64]uint16) {
	if dataExport.active == false {
		return
	}

	frameData := FrameData{variables, channelPC}
	dataExport.frames[dataExport.currentFrame] = frameData
	dataExport.currentFrame++

	if dataExport.currentFrame%501 == 500 {
		logger.Info("dumped frame %d", dataExport.currentFrame)
	}

	if dataExport.currentFrame >= framesToCapture {
		dataExport.active = false

		logger.Info("Create JSON Export")
		enc := json.NewEncoder(os.Stdout)
		out := SerializedJsonDataFormat{2, "filename", 25, dataExport.frames}
		//d := map[string]int{"apple": 5, "lettuce": 7}
		//enc.Encode(d)
		enc.Encode(out)
	}
}
