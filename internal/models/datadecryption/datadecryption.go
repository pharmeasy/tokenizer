package datadecryption

// DecryptionResponse represents the successful decryption response
type DecryptionResponse struct {
	DecryptionResponseData []DecryptionResponseData `json:"data"`
}

// DecryptionResponseData is the struct for object
type DecryptionResponseData struct {
	Token    string `json:"token"`
	Content  string `json:"content"`
	Metadata string `json:"metadata"`
}
