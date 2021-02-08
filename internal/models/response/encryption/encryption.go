package encryption

// Response represents the successful encryption response
type Response struct {
	Data []Data `json:"data"`
}

// Data is the struct for object
type Data struct {
	ID     int    `json:"id"`
	Token  string `json:"token"`
	Cipher []byte `json:"cipher"`
}
