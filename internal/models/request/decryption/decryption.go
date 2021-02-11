package decryption

// Request represents the successful decryption request
type Request struct {
	RequestID  string `json:"requestId"`
	Identifier string `json:"identifier"`
	Level      int    `json:"level"`
	Data       []Data `json:"data"`
}

// Data represents the data field in the incoming decryption request
type Data struct {
	Token string `json:"token"`
	Salt  string `json:"salt"`
}
