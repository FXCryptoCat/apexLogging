package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/william1034/apexLogging/cmd/apex_log_monitor/libs"
	"github.com/william1034/apexLogging/cmd/common_flags"
	"github.com/william1034/apexLogging/internal/apex_monitor"
	"github.com/william1034/apexLogging/internal/influx"
	"github.com/william1034/apexLogging/internal/tick"
	"os"
	. "time"
)

func main() {

	flags := common_flags.GetFlags()
	if *flags.Install {
		log.SetLevel(log.TraceLevel) //Use this until we have gotten past the install
		libs.InstallAsService(flags)
		fmt.Println("Please run the following to complete the installation.")
		fmt.Println("sudo systemctl daemon-reload && ")
		fmt.Println("sudo systemctl enable apexMonitor.service &&")
		fmt.Println("sudo systemctl start apexMonitor.service &&")
		fmt.Println("sudo systemctl status apexMonitor.service ")
		fmt.Println()
		os.Exit(0)
	}

	configureLogs(flags.LogLevel, flags.LogFile)
	realTimeMonitor(flags)
}

func configureLogs(logLevel string, logFile string) {

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.SetLevel(level)

	if len(logFile) > 0 {
		f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal("Unable to create log file " + logFile)
		}
		log.SetOutput(f)
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}


func realTimeMonitor(flags common_flags.CommonFlags) {
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
				tick := tick.GetTickLineFromApexStatus(*apexStatus, flags.TickLogFile)
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



