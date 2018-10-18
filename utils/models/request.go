package models

type Request struct {
	IniPoint Point `json:"ini_point"`
	EndPoint Point `json:"end_point"`
}
