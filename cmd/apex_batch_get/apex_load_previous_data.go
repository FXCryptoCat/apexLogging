package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/william1034/apexLogging/cmd/common_flags"
	"github.com/william1034/apexLogging/internal/apex_monitor"
	"github.com/william1034/apexLogging/internal/influx"
	"github.com/william1034/apexLogging/internal/tick"
	"os"
	"time"
)


func main() {

	flags := common_flags.GetFlags()
	configureLogs(flags.LogLevel)

	apexClient := apex_monitor.NewAsyncApexMonitor("http://" + flags.ApexIp, flags.ApexCookie)
	ic :=influx.NewInfluxClient(flags.InfluxIp, flags.InfluxUser, flags.InfluxPassword)

	from := time.Date(2019, 11, 11, 0, 0, 0, 0, time.UTC)

	apexLogChannel := make(chan apex_monitor.ApexLog)
	go apexClient.GetSummaryFrom(from, apexLogChannel)

	for {
		apexLog, ok := <- apexLogChannel
		if ok == false {
			fmt.Println("-- DONE --")
			break
		} else {
			lines := tick.GetTickRecordsFromApexLog(apexLog)

			if !*flags.DisableDataSource {
				ic.WriteRecords(lines)
			}
			fmt.Println(lines)
		}

	}
}


func configureLogs(logLevel string) {

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.SetLevel(level)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

