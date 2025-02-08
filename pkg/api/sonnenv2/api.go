package sonnenv2

type Status struct {
	ApparentOutput            float32     `json:"Apparent_output"`
	BackupBuffer              string      `json:"BackupBuffer"`
	BatteryCharging           bool        `json:"BatteryCharging"`
	BatteryDischarging        bool        `json:"BatteryDischarging"`
	ConsumptionW              float32     `json:"Consumption_W"`
	Fac                       float64     `json:"Fac"`
	FlowConsumptionBattery    bool        `json:"FlowConsumptionBattery"`
	FlowConsumptionGrid       bool        `json:"FlowConsumptionGrid"`
	FlowConsumptionProduction bool        `json:"FlowConsumptionProduction"`
	FlowGridBattery           bool        `json:"FlowGridBattery"`
	FlowProductionBattery     bool        `json:"FlowProductionBattery"`
	FlowProductionGrid        bool        `json:"FlowProductionGrid"`
	GridFeedInW               float32     `json:"GridFeedIn_W"`
	IsSystemInstalled         float32     `json:"IsSystemInstalled"`
	OperatingMode             string      `json:"OperatingMode"`
	PacTotalW                 float32     `json:"Pac_total_W"`
	ProductionW               float32     `json:"Production_W"`
	Rsoc                      float32     `json:"RSOC"`
	Sac1                      float32     `json:"Sac1"`
	Sac2                      interface{} `json:"Sac2"`
	Sac3                      interface{} `json:"Sac3"`
	SystemStatus              string      `json:"SystemStatus"`
	Timestamp                 string      `json:"Timestamp"`
	Usoc                      float32     `json:"USOC"`
	Uac                       float32     `json:"Uac"`
	Ubat                      float32     `json:"Ubat"`
	DischargeNotAllowed       bool        `json:"dischargeNotAllowed"`
	GeneratorAutostart        bool        `json:"generator_autostart"`
}
