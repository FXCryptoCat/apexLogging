package main

import (

	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/william1034/apexLogging/cmd/common_flags"
	"github.com/william1034/apexLogging/internal/apex_monitor"
	"github.com/william1034/apexLogging/internal/influx"
	"github.com/william1034/apexLogging/internal/tick"
	"os"
	. "time"
)

//scp go_build_apexlogs_PI_linux willum@192.168.20.169:~

func main() {

	flags := common_flags.GetFlags()
	configureLogs(flags.LogLevel)
	realTimeMonitor(flags)
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

func realTimeMonitor(flags common_flags.CommonFlags ) {
	////The apex actually brings everything back. Slowly
	duration := Duration(*flags.Delay) * Second
	ticker := NewTicker(duration)
	quit := make(chan struct{})

	apexClient := apex_monitor.NewApexMonitor(flags.ApexIp, flags.ApexUserName, flags.ApexPassword)
	ic := influx.NewInfluxClient(flags.InfluxIp, flags.InfluxUser, flags.InfluxPassword)

	go func() {
		for {
			select {
			case <-ticker.C:
				nanos := Now().UnixNano()
				apexStatus, err := apexClient.GetStatus(nanos)

				if err != nil {
					fmt.Println(err)
					//Sleep for a few seconds before continuing....
					Sleep(Millisecond * 5000)

					//Try to reauth
					apexClient.ReAuth()
					continue
				}

				log.Trace("Fetching another record")
				tick := tick.GetTickLineFromApexStatus(*apexStatus)
				if *flags.DisableDataSource {
					log.Debug(tick)
				} else {
					ic.WriteRecord(tick)
				}

				//fmt.Print(tick)

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-quit

}
