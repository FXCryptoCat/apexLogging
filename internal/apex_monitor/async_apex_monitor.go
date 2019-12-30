package apex_monitor

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

//This is used when loading previous data. Too much copied code from the apex_monitor.
//Todo: refactor.

//The AsyncApexMonitor writes all ApexLogs to the logChannel as they are created. See getOneDayApexLog
//TODO: Refactor the channel into the constructor
type AsyncApexMonitor struct {
	baseUri string
	cookie string
}

func NewAsyncApexMonitor(uri string, cookie string) *AsyncApexMonitor {
	monitor := AsyncApexMonitor{baseUri: uri, cookie: cookie}
	return &monitor
}


func (asyncClient *AsyncApexMonitor) GetSummaryFrom(fromTime time.Time, logChannel chan ApexLog) {
	defer close(logChannel)
	asyncClient.getLog("ilog", fromTime, logChannel)
	asyncClient.getLog("olog", fromTime, logChannel)

}

func (asyncClient *AsyncApexMonitor) getLog(logName string, fromTime time.Time, logChannel chan ApexLog) {
	log.Trace("Getting the " + logName)
	currentTime := time.Now()
	currentTimeSeconds := currentTime.Unix()


	diffDuration := currentTime.Sub(fromTime)
	diffDays := int(diffDuration.Hours() / 24)
	if diffDays == 0 {
		//You might be looking for a few hours today
		diffDays = 1
	}

	for i := 0; i < diffDays; i++ {
		asyncClient.getOneDayApexLog(logName, fromTime, currentTimeSeconds, logChannel)
		fromTime = fromTime.Add(time.Hour * 24) //Add one day and try again
	}

}

func (asyncClient *AsyncApexMonitor) getOneDayApexLog(logName string, fromTime time.Time, currentTimeSeconds int64, logChannel chan ApexLog)  {
	uri := fmt.Sprintf("%s/rest/%s?days=%d&sdate=%d&_=%d", asyncClient.baseUri, logName, 1, FormatDate(fromTime), currentTimeSeconds)

	data, _ := RequestWithCookies(uri, asyncClient.cookie)

	//Each summary is for one day
	apexSummary := ApexComboLog{}
	//apexSummary := ApexILog{}
	err := json.Unmarshal(data, &apexSummary)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var apexLog ApexLog
	if len(apexSummary.ILog.Record) > 0 {
		apexLog = ApexLog(apexSummary.ILog)
	} else if len(apexSummary.OLog.Record) > 1 {
		apexLog = *ConvertOLogToApexLog(apexSummary.OLog) //Begging to blow up with nil ptr exception
	} else {
		log.Warn("No Logs Found")
		return
	}

	logChannel <- apexLog

}
