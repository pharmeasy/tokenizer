package decryption

// Response represents the successful decryption response
type Response struct {
	Data []Data `json:"data"`
}

// Data is the struct for object
type Data struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}
