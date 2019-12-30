package apex_monitor

import (
	"github.com/sirupsen/logrus"
)

//The ApexInputLog and Output Log only differ by a little. We convert all of the
//ApexOutputLog to ApexInputLog to make logging easier. Use the ApexLog alias
//to ensure we hav the right thing.
type ApexLog ApexInputLog

//The Input Log structure returned from the ilog call
//used when loading Chart page on the classic dashboard.

type ApexInputLog struct {
	Hostname string `json:"hostname"`
	Software string `json:"software"`
	Hardware string `json:"hardware"`
	Serial   string `json:"serial"`
	Type     string `json:"type"`
	Extra    struct {
		Sdver string `json:"sdver"`
	} `json:"extra"`

	Timezone string          `json:"timezone"`
	Date     int             `json:"date,omitempty"`
	Record   []ApexLogRecord `json:"record"`
}

//Converts a ApexOutputLog to an ApexLog (A ApexInputLog in disguise)
func ConvertOLogToApexLog(olog ApexOutputLog) *ApexLog {
	//For our purposes we don't care about records that are missing records.
	if len(olog.Record) == 0 {
		logrus.Warn("ApexOutputLog being converted to ApexLog has no records. Skipping")
		return nil
	}
	//The only difference is the ApexInputLog:: Record is a list of ApexLogRecord
	//ApexOutputLog is a list of ApexRecordData

	alog := ApexLog{
		Hostname: olog.Hostname,
		Software: olog.Software,
		Hardware: olog.Hardware,
		Serial:   olog.Serial,
		Type:     olog.Type,
		Extra:    olog.Extra,
		Timezone: olog.Timezone,
		Date:     olog.Date,
	}

	var apexLogRecordList []ApexLogRecord

	//Build the new Record from the olog
	for _, apexRecordData := range olog.Record {
		v := "0"
		if apexRecordData.Value == "ON" || apexRecordData.Value == "AON" {
			v = "1"
		}
		newData := ApexRecordData{
			Did:   apexRecordData.Did,
			Name:  apexRecordData.Name,
			Value: v,
		}

		apexLogRecord := ApexLogRecord{
			Date: apexRecordData.Date,
			Data: []ApexRecordData{newData},
		}
		apexLogRecordList = append(apexLogRecordList, apexLogRecord)
	}
	alog.Record = apexLogRecordList

	return &alog
}
