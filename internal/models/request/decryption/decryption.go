package decryption

// Request represents the incoming decryption request
type Request struct {
	RequestID	int 	`json:"requestId"`
	Source		string	`json:"source"`
	Data		Data 	`json:"data"`
}

// Data represents the data field in the incoming decryption request
type Data struct {
	Tokens	[]string	`json:"tokens"`
}
