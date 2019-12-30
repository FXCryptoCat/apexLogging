package apex_monitor

//The JSON returned from the "http://YOURIP/rest/status?_=TIMEINMILLIS" that is
//called when loading the main dashboard.
type ApexStatus struct {
	System struct {
		Hostname string `json:"hostname"`
		Software string `json:"software"`
		Hardware string `json:"hardware"`
		Serial   string `json:"serial"`
		Type     string `json:"type"`
		Extra    struct {
			Sdver string `json:"sdver"`
		} `json:"extra"`
		Timezone string `json:"timezone"`
		Date     int    `json:"date"`
	} `json:"system"`
	Modules []struct {
		Abaddr  int    `json:"abaddr"`
		Hwtype  string `json:"hwtype"`
		Hwrev   int    `json:"hwrev"`
		Swrev   int    `json:"swrev"`
		Swstat  string `json:"swstat"`
		Pcount  int    `json:"pcount"`
		Pgood   int    `json:"pgood"`
		Perror  int    `json:"perror"`
		Reatt   int    `json:"reatt"`
		Inact   int    `json:"inact"`
		Boot    bool   `json:"boot"`
		Present bool   `json:"present"`
		Extra   struct {
		} `json:"extra"`
	} `json:"modules"`
	Nstat struct {
		Dhcp           bool     `json:"dhcp"`
		Hostname       string   `json:"hostname"`
		Ipaddr         string   `json:"ipaddr"`
		Netmask        string   `json:"netmask"`
		Gateway        string   `json:"gateway"`
		DNS            []string `json:"dns"`
		HTTPPort       int      `json:"httpPort"`
		FusionEnable   bool     `json:"fusionEnable"`
		Quality        int      `json:"quality"`
		Strength       int      `json:"strength"`
		Link           bool     `json:"link"`
		WifiAPLock     bool     `json:"wifiAPLock"`
		WifiEnable     bool     `json:"wifiEnable"`
		WifiAPPassword string   `json:"wifiAPPassword"`
		Ssid           string   `json:"ssid"`
		WifiAP         bool     `json:"wifiAP"`
		EmailPassword  string   `json:"emailPassword"`
		UpdateFirmware bool     `json:"updateFirmware"`
		LatestFirmware string   `json:"latestFirmware"`
	} `json:"nstat"`
	Feed struct {
		Name   int `json:"name"`
		Active int `json:"active"`
	} `json:"feed"`
	Power struct {
		Failed   int `json:"failed"`
		Restored int `json:"restored"`
	} `json:"power"`
	Outputs []struct {
		Status []string `json:"status"`
		Name   string   `json:"name"`
		Gid    string   `json:"gid"`
		Type   string   `json:"type"`
		ID     int      `json:"ID"`
		Did    string   `json:"did"`
	} `json:"outputs"`
	Inputs []struct {
		Did   string  `json:"did"`
		Type  string  `json:"type"`
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	} `json:"inputs"`
	Link struct {
		LinkState int    `json:"linkState"`
		LinkKey   string `json:"linkKey"`
		Link      bool   `json:"link"`
	} `json:"link"`
}

