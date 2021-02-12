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

// EncryptRequest represents the successful encryption request
type EncryptRequest struct {
	RequestID   string        `json:"requestId"`
	Identifier  string        `json:"identifier"`
	Level       int           `json:"level"`
	RequestData []RequestData `json:"data"`
}

// RequestData represents the data field in the incoming encryption request
type RequestData struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Salt     string `json:"salt"`
	Metadata string `json:"metadata"`
}
