package apex_monitor

//A little hackery to unmarshall the ilog and ologs in one place.
//the ologs will be converted to a ApexLog(ApexInputLog in disguise)
type ApexComboLog struct {
	OLog ApexOutputLog `json:"olog,omitempty"`
	ILog ApexInputLog `json:"ilog,omitempty"`
}

type ApexLogRecord struct {
	Date int64            `json:"date,omitempty"`
	Data []ApexRecordData `json:"data"`
}

type ApexRecordData struct {
	Name  string `json:"name"`
	Date  int64  `json:"date,omitempty"`
	Did   string `json:"did"`
	Value string `json:"value"`
}
