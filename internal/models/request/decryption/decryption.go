package decryption

// Request represents the successful decryption request
type Request struct {
	RequestID  string `json:"requestId"`
	Identifier string `json:"identifier"`
	Level      string `json:"level"`
	Data       Data   `json:"data"`
}

// Data represents the data field in the incoming decryption request
type Data struct {
	Tokens []Token `json:"tokens"`
}

// Token represents an object which will contain the tokenized string and salt
type Token struct {
	ID   string `json:"id"`
	Salt string `json:"salt"`
}
