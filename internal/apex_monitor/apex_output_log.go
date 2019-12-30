package apex_monitor

type ApexOutputLog struct {
	Hostname string `json:"hostname"`
	Software string `json:"software"`
	Hardware string `json:"hardware"`
	Serial   string `json:"serial"`
	Type     string `json:"type"`
	Extra    struct {
		Sdver string `json:"sdver"`
	} `json:"extra"`

	Timezone string `json:"timezone"`
	Date     int    `json:"date,omitempty"`
	//This is the differences between the ApexOutputLog and ApexInputLog.
	//The input log has a list of ApexLogRecord
	Record   []ApexRecordData `json:"record"`
}


