package apex_monitor


//The ApexInputLog and Output Log only differ by a little. We convert all of the
//ApexOutputLog to ApexInputLog to make logging easier. Use the ApexLog alias
//to ensure we hav the right thing.
type ApexLog ApexInputLog

//This allows you to unmarshal input our output logs with a
//single command. Ultimately the output log will be converted
//to an ApexLog.
type ApexComboLog struct {
	OLog ApexOutputLog `json:"olog,omitempty"`
	ILog ApexInputLog `json:"ilog,omitempty"`
}

//This is the "Record" in the ilog. It has a date, then
//the measurements in the data use that date. Output logs
//Have the date embedded in each record.
type ApexLogRecord struct {
	Date int64            `json:"date,omitempty"`
	Data []ApexRecordData `json:"data"`
}

//Included the date with omitempty to support the conversion from
//ApexOutputLog in ApexLog
type ApexRecordData struct {
	Name  string `json:"name"`
	Date  int64  `json:"date,omitempty"`
	Did   string `json:"did"`
	Value string `json:"value"`
}
