package dto

type AlarmCenter struct {
	Lat    float64 `json:"lat" validate:"required,latitude"`
	Lng    float64 `json:"lng" validate:"required,longitude"`
	Radius int     `json:"radius" validate:"required,min=0"` // meters
}

type AlarmDispatchRequest struct {
	AlarmID string      `json:"alarm_id" validate:"required"`
	Urgency string      `json:"urgency" validate:"required,oneof=immediate high normal low"`
	Center  AlarmCenter `json:"center" validate:"required"`
	Signal  string      `json:"signal" validate:"required"`
	Content string      `json:"content" validate:"required"`
}
