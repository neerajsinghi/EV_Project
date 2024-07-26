package entity

type IotBikeDB struct {
	Alarm                    string   `bson:"alarm" json:"alarm"`
	BatteryLevel             float64  `bson:"batteryLevel" json:"batteryLevel"`
	Course                   string   `bson:"course" json:"course"`
	Dealer                   string   `bson:"dealer" json:"dealer"`
	DeviceFixTime            string   `bson:"deviceFixTime" json:"deviceFixTime"`
	DeviceId                 int      `bson:"deviceId" json:"deviceId"`
	DeviceImei               string   `bson:"deviceImei" json:"deviceImei"`
	HarshAccelerationHistory []string `bson:"harshAccelerationHistory" json:"harshAccelerationHistory"`
	HarshBrakingHistory      []string `bson:"harshBrakingHistory" json:"harshBrakingHistory"`
	Ignition                 string   `bson:"ignition" json:"ignition"`
	LastUpdate               string   `bson:"lastUpdate" json:"lastUpdate"`
	Latitude                 string   `bson:"latitude" json:"latitude"`
	Longitude                string   `bson:"longitude" json:"longitude"`
	Location                 Location `bson:"location" json:"location"`
	Name                     string   `bson:"name" json:"name"`
	Phone                    string   `bson:"phone" json:"phone"`
	PosId                    int      `bson:"posId" json:"posId"`
	Speed                    string   `bson:"speed" json:"speed"`
	Status                   string   `bson:"status" json:"status"`
	TotalDistance            string   `bson:"totalDistance" json:"totalDistance"`
	Type                     string   `bson:"type" json:"type"`
	Valid                    int      `bson:"valid" json:"valid"`
}
type Location struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}
