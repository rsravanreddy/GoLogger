package api

type Success struct {
	Success string `json:"success"`
	Length int `json:"length"`

}

func NewSuccess(msg string,len int) *Success {
	return &Success{Success: msg,Length:len}
}
