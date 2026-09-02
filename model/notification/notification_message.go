package notification

type NotificationMessage struct {
	CustomerNo     string  `json:"customer_no"`
	Channel        string  `json:"channel"`
	Site           string  `json:"site"`
	Category       string  `json:"category"`
	Title          string  `json:"title"`
	Content        string  `json:"content"`
	Timestamp      string  `json:"timestamp"`
	Event          string  `json:"event"`
	DateLog        string  `json:"date_log"`
	TimeLog        string  `json:"time_log"`
	Shift          string  `json:"shift"`
	ShiftSequence  int64   `json:"shift_sequence"`
	EquipmentNo    string  `json:"equipment_no"`
	EquipmentType  string  `json:"equipment_type"`
	EquipmentModel string  `json:"equipment_model"`
	Username       string  `json:"username"`
	Name           string  `json:"name"`
	Fleet          string  `json:"fleet"`
	State          string  `json:"state"`
	Type           string  `json:"type"`
	Reason         string  `json:"reason"`
	Activity       string  `json:"activity"`
	StatusActivity string  `json:"status_activity"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Altitude       float64 `json:"altitude"`
	Bearing        float64 `json:"bearing"`
}
