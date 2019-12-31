package tick

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/william1034/apexLogging/internal/apex_monitor"
	"strconv"
	"time"
)

//Format of a tick record. Note the {Measurement, Value} pairs are space separated.
//apex, type = {TYPE} {MeasurementName, Value} {MeasurementName, Value}... {TimeInNanoSeconds}

//Gets  list of tick records (1 per input)
func GetTickLineFromApexStatus(apexStatus apex_monitor.ApexStatus) string {
	log.Trace("GetTickLineFromApexStatus")
	var inputRecords = apexStatus.Inputs
	line := ""

	//Add the input measurements to the TickLine. Note the
	for _, inputRecord := range inputRecords {
		line = addMeasurementToTickLine(line, inputRecord.Name, inputRecord.Value)
	}

	outputRecords := apexStatus.Outputs
	for _, output := range outputRecords {
		if output.Type != "outlet" && output.Type != "24v" && output.Type != "alert" {
			continue
		}
		//See README for the meaning of the status array
		statusString := output.Status[0]
		stateString := output.Status[2]
		status := 0
		if statusString == "ON" || statusString == "AON" {
			status = 1
		}
		errorState := 0
		if stateString != "OK" {
			errorState = 1
		}
		line = addMeasurementToTickLine(line, output.Name, status)
		line = addMeasurementToTickLine(line, output.Name+".error", errorState)
	}

	nanoTime := SecToNanoSeconds(int64(apexStatus.System.Date))
	//BUG: TODO: The type is always "temp" fix it.
	line = fmt.Sprintf("apex,type=temp %s %d\n", line, nanoTime)
	return line
}

//The ApexLog contains records with different times. Each time
//requires its own line.
func GetTickRecordsFromApexLog(apexLog apex_monitor.ApexLog) []string {
	log.Trace("GetTickRecordsFromApexLog")
	var lines []string
	for _, record := range apexLog.Record {
		var line string
		for _, data := range record.Data {
			value, _ := strconv.ParseFloat(data.Value, 64)
			line = addMeasurementToTickLine(line, data.Name, value)
		}
		line = fmt.Sprintf("apex,type=temp %s %d", line, SecToNanoSeconds(record.Date))
		lines = append(lines, line)
	}
	//log.Debugf("GetTickRecordsFromApexLog lines: %v\n", lines)
	return lines
}

func SecToNanoSeconds(seconds int64) int64 {
	x := seconds * int64(1000000000) * int64(time.Nanosecond)
	return x
}

func addMeasurementToTickLine(line string, name string, value interface{}) string {
	if len(line) == 0 {
		line = fmt.Sprintf("%s=%v", name, value)
	} else {
		line = fmt.Sprintf("%s,%s=%v", line, name, value)
	}
	return line
}