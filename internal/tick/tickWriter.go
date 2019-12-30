package tick

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/william1034/apexLogging/internal/apex_monitor"
	"strconv"
	"time"
)

//Some functions for writing tick lines.
//type TickWriter struct {
//}
//
//func NewTickWriter() *TickWriter {
//	tw := TickWriter{}
//	return &tw
//}

//func (tickWriter *TickWriter) CreateTickFileFromIlog(ilog apex_log_monitor.ApexInputLog) {
//	log.Trace("CreateTickFileFromIlog")
//	lines := tickWriter.GetTickRecordFromILog(ilog)
//
//	toDate := time.Unix(0,time. Now().UnixNano()) //Just keeps them in order
//	fileBaseName := toDate.Format("060102_030405")
//	fileName := "Apex" + fileBaseName + ".tick"
//
//	tickWriter.storeTickFile(fileName, lines)
//
//}

func  CreateTickFileFromApexStatus(ilog apex_monitor.ApexInputLog) {
	log.Trace("CreateTickFileFromApexStatus")
}

//Gets  list of tick records (1 per input)
func  GetTickRecord(apexStatus apex_monitor.ApexStatus) string {
	log.Trace("GetTickRecord")
	var records = apexStatus.Inputs
	first := true
	var line string

	for _, input := range records {
		if first {
			line = fmt.Sprintf("%s=%f", input.Name, input.Value)
			first = false
		} else {
			line = fmt.Sprintf("%s,%s=%f", line, input.Name, input.Value)
		}
	}

	outputRecords := apexStatus.Outputs

	for _, output := range outputRecords {
		if output.Type != "outlet" && output.Type != "24v" && output.Type != "alert"{
			continue
		}

		statusString := output.Status[0]
		stateString := output.Status[2]
		status := 0
		if statusString == "ON" || statusString == "AON" {
			status = 1
		}
		errorState := 0
		if stateString != "OK"  {
			errorState = 1
		}
		if first {
			line = fmt.Sprintf("%s=%d", output.Name, status)
			first = false
		} else {
			line = fmt.Sprintf("%s,%s=%d", line, output.Name, status)
		}
		line = fmt.Sprintf("%s,%s=%d", line, output.Name + ".error", errorState)
	}

	nanoTime := SecToNanoSeconds(int64(apexStatus.System.Date))
	line = fmt.Sprintf("apex,type=temp %s %d\n", line, nanoTime)
	return line
}

func  GetTickRecordsFromIlog(ilog apex_monitor.ApexLog) []string {
	log.Trace("GetTickRecordsFromIlog")
	var lines []string
	for _, record := range ilog.Record {
		first := true
		var line string
		for _, data := range record.Data {
			value, _ := strconv.ParseFloat(data.Value, 64)
			if first {
				line = fmt.Sprintf("%s=%f", data.Name, value)
				first = false
			} else {
				line = fmt.Sprintf("%s,%s=%f", line, data.Name, value)
			}
		}
		line = fmt.Sprintf("apex,type=temp %s %d", line, SecToNanoSeconds(record.Date) )
		//log.Debugf("RecordTime: %f, to nano:  %d\n", record.Date, SecToNanoSeconds(record.Date))
		lines = append(lines, line)
	}
	//log.Debugf("GetTickRecordsFromIlog lines: %v\n", lines)
	return lines
}

func GetTickRecordFromILog(ilog apex_monitor.ApexInputLog) []string {
	log.Trace("GetTickRecordFromILog")
	first := true
	var line string
	var lines []string

	for _, record := range ilog.Record {
		for _, input := range record.Data {
			if first {
				line = fmt.Sprintf("%s=%s", input.Name, input.Value)
				first = false
			} else {
				line = fmt.Sprintf("%s,%s=%s", line, input.Name, input.Value)
			}
		}
		line = fmt.Sprintf("apex,type=temp %s %d", line, SecToNanoSeconds(record.Date))
		lines = append(lines, line)
		line = ""
	}
	return lines
	//Date is in seconds
}
/*
func (tickWriter *TickWriter) storeTickFile(filename string, lines []string) {
	fqFileName := path.Join(tickWriter.fileDir, filename)
	file, err := os.OpenFile(fqFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range lines {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	file.Close()

}
func (tickWriter *TickWriter) storeJsonFile(filename string, json []byte) {
	if len(json) > 0 {

		fqfilename := fmt.Sprint(tickWriter.fileDir, "/", filename)
		f, err := os.Create(fqfilename)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		_, _ = f.Write(json)
		_ = f.Sync()

		f.Close()
	}
} */


func SecToNanoSeconds(seconds int64 ) int64 {
	x:= seconds * int64(1000000000) * int64(time.Nanosecond)
	return x
}