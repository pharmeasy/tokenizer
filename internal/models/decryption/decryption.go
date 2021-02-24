package decryption

// DecryptResponse represents the successful decryption response
type DecryptResponse struct {
	DecryptionResponseData []DecryptResponseData `json:"data"`
}

// DecryptResponseData is the struct for object
type DecryptResponseData struct {
	Token    string `json:"token"`
	Content  string `json:"content"`
	Metadata string `json:"metadata"`
}

// DecryptRequest represents the successful decryption request
type DecryptRequest struct {
	RequestID          string               `json:"requestId"`
	Identifier         string               `json:"identifier"`
	DecryptRequestData []DecryptRequestData `json:"data"`
}

// DecryptRequestData represents the data field in the incoming decryption request
type DecryptRequestData struct {
	Token string `json:"token"`
	Salt  string `json:"salt"`
}
