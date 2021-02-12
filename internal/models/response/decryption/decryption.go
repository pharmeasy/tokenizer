package decryption

// Response represents the successful decryption response
type Response struct {
	Data []Data `json:"data"`
}

// Data is the struct for object
type Data struct {
	Token    string `json:"token"`
	Content  string `json:"content"`
	Metadata string `json:"metadata"`
}
