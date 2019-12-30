package apex_monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ApexClient struct {
	//The URI to your apex.
	baseUri string
	cookie string
}

func NewApexMonitor(uri string, cookie string) *ApexClient {
	monitor := ApexClient{baseUri: uri, cookie: cookie}

	return &monitor
}


//Retrieve the status of the apex inputs and outputs at a given time.
func (apexClient *ApexClient) GetStatus(fromTimeInNanoSeconds int64) (*ApexStatus, error) {
	log.Trace("GetStatus")
	millis := fromTimeInNanoSeconds / 1000000

	//This URI has been plucked from the series of http calls used to load the
	//APEX classic console main dashboard. It loads everything.
	uri := fmt.Sprintf("%s/rest/status?_=%d", apexClient.baseUri, millis)
	data, _ := RequestWithCookies(uri, apexClient.cookie)

	//Unmarshal the ApexStatus from the byte stream.
	status, err := apexClient.getApexStatusFromBytes(data)
	if err != nil {
		log.Error("Failed to get status from apex. ", err)
		return nil, err
	}
	return status, nil
}

////Get the ApexStatus from a window in time.
//func (apexClient *ApexClient) GetStatusDuringTimeWindow(timeWindow internal.TimeWindow) (*ApexStatus, error) {
//	log.Trace("GetStatusDuringTimeWindow")
//	//Use todays data
//	toDate := time.Unix(0, timeWindow.ToNanos)
//	//Always subtract a day for the stupid apex
//	dateString := toDate.AddDate(0, 0, -1).Format("060102")
//
//	uri := fmt.Sprintf("%s/rest/status?_=%s&_=%d", apexClient.baseUri, dateString, toDate.Unix())
//	data, _ := RequestWithCookies(uri)
//
//	apexStatus, err := apexClient.getApexStatusFromBytes(data)
//	if err != nil {
//		log.Error("Failed to get status from apex, "+uri, err)
//		return nil, err
//	}
//	return apexStatus, nil
//}

//func (apexClient *ApexClient) GetStatusDuringDates(fromTime time.Time, toTime time.Time) []ApexStatus {
//	//toTime := fromTime.Add(time.Hour * 24)
//	log.Trace("GetStatusDuringDates")
//	var apexLogs []ApexStatus
//
//	diffDuration := toTime.Sub(fromTime)
//	diffDays := int(diffDuration.Hours() / 24)
//
//	from := fromTime
//	for d := 0; d < diffDays; d++ {
//		timeWindow := internal.TimeWindow{
//			FromNanos: from.UnixNano(),
//			ToNanos:   from.Add(time.Hour * 24).UnixNano(),
//		}
//		apexStatus, err := apexClient.GetStatusDuringTimeWindow(timeWindow)
//		if err != nil {
//			//fmt.Println("Hopefully a transient error")
//			//fmt.Println(err)
//			log.Warn(err)
//			from.Add(time.Hour * 24)
//			continue
//		}
//		apexLogs = append(apexLogs, *apexStatus)
//		from.Add(time.Hour * 24)
//	}
//
//	return apexLogs
//}

//Returns a summary that contains all of the readings from the specified date.
//APEX stores ~2 months of data
func (apexClient *ApexClient) GetSummaryFrom(fromTime time.Time) []ApexLog {
	apexSummaryList := apexClient.getLog("ilog", fromTime)

	//apexSummaryList := apexClient.getLog("olog", fromTime)

	////Need a generic slice add all function
	ologs := apexClient.getLog("olog", fromTime)
	for _, olog := range ologs {
		apexSummaryList = append(apexSummaryList, olog)
	}

	return apexSummaryList
}

//Turns out the "graph" page loads everything
//ilog(inputs), olog(outputs), dlog(dos), tlog(??)
//We only care about ilog and olog
func (apexClient *ApexClient) getLog(logName string, fromTime time.Time) []ApexLog {
	log.Trace("Getting the " + logName)
	currentTime := time.Now()
	currentTimeSeconds := currentTime.Unix()

	var apexSummaryList []ApexLog

	diffDuration := currentTime.Sub(fromTime)
	diffDays := int(diffDuration.Hours() / 24)
	if diffDays == 0 {
		//You might be looking for a few hours today
		diffDays = 1
	}

	for i := 0; i < diffDays; i++ {
		alog := apexClient.getOneDayApexLog(logName, fromTime, currentTimeSeconds)
		apexSummaryList = append(apexSummaryList, *alog)
		fromTime = fromTime.Add(time.Hour * 24) //Add one day and try again
	}

	return apexSummaryList
}

//Get the apexLog for a single day
func (apexClient *ApexClient) getOneDayApexLog(logName string, fromTime time.Time, currentTimeSeconds int64) *ApexLog {

	//Based on the calls to the graph page on the classic apex
	uri := fmt.Sprintf("%s/rest/%s?days=%d&sdate=%d&_=%d", apexClient.baseUri, logName, 1, FormatDate(fromTime), currentTimeSeconds)

	data, _ := RequestWithCookies(uri, apexClient.cookie)

	apexComboLog := ApexComboLog{}
	err := json.Unmarshal(data, &apexComboLog)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Some funny stuff to convert the ApexOutputLog to an ApexLog
	var apexLog ApexLog
	if len(apexComboLog.ILog.Record) > 0 {
		apexLog = ApexLog(apexComboLog.ILog)
	} else if len(apexComboLog.OLog.Record) > 1 {
		apexLog = *ConvertOLogToApexLog(apexComboLog.OLog) //Begging to blow up with nil ptr exception
	} else {
		log.Warn("No Logs Found")
		return nil
	}

	return &apexLog
}


func (apexClient *ApexClient) getApexStatusFromBytes(data []byte) (*ApexStatus, error) {
	log.Trace("getApexStatusFromBytes")
	var statusLog ApexStatus
	err := json.Unmarshal([]byte(data), &statusLog)
	if err != nil {
		log.WithField("JSON", string(data)).Error("Failed to unmarshal apex status.", err)

		return nil, err

	}
	return &statusLog, nil
}

func RequestWithCookies(url string, apexCookie string) (b []byte, err error) {
	log.Trace("requestWithCookies")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Failed request with cookies.", err)
		return
	}
	req.AddCookie(&http.Cookie{Name: "connect.sid", Value: apexCookie})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Warn("Request with cookies status code not 200.", err)
		err = errors.New(url +
			"\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return nil, err
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	return bytes, nil
}

//Running this on a raspberry pi gave me all sorts of heck with int64 times overflowing.
//not sure how this helped, but it did. (Hey its late)
func FormatDate(t time.Time) int {
	timeString := t.Format("060102")
	time, _ := strconv.Atoi(timeString)
	return time
}
