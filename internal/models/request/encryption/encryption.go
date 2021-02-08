package encryption

// Request represents the successful encryption request
type Request struct {
	RequestID int    `json:"requestId"`
	Source    string `json:"source"`
	Level     int    `json:"level"`
	Data      []Data `json:"data"`
}

// Data represents the data field in the incoming encryption request
type Data struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Salt     string `json:"salt"`
	MetaData string `json:"metadata"`
}
