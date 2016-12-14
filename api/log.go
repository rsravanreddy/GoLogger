package api

// type Log struct {
// 	LogPayload LogPayloads
// }

type Log struct {
	NodePath    string `json:"node_path"`
	ProcessPath string `json:"process_path"`
	TimeStamp   string `json:"time_stamp"`
	NodeId      int    `json:"node_id"`
	Action      string `json:"action"`
}
