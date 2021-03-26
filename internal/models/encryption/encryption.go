package encryption

// EncryptResponse represents the successful encryption response
type EncryptResponse struct {
	ResponseData []ResponseData `json:"data"`
}

// ResponseData is the struct for object
type ResponseData struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// EncryptRequest represents the successful encryption request
type EncryptRequest struct {
	RequestID          string               `json:"requestId"`
	Identifier         string               `json:"identifier"`
	EncryptRequestData []EncryptRequestData `json:"data"`
}

// EncryptRequestData represents the data field in the incoming encryption request
type EncryptRequestData struct {
	ID       string            `json:"id"`
	Level    string            `json:"level"`
	Content  string            `json:"content"`
	Salt     string            `json:"salt"`
	Metadata map[string]string `json:"metadata"`
}
