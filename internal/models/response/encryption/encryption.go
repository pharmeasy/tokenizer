package encryption

// EncryptResponse represents the successful encryption response
type EncryptResponse struct {
	ResponseData []ResponseData `json:"data"`
}

// ResponseData is the struct for object
type ResponseData struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

// Request represents the successful encryption request
type Request struct {
	RequestID  string `json:"requestId"`
	Identifier string `json:"identifier"`
	Level      string `json:"level"`
	Data       []Data `json:"data"`
}

// Data represents the data field in the incoming encryption request
type Data struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	Salt     string `json:"salt"`
	MetaData string `json:"metadata"`
}
