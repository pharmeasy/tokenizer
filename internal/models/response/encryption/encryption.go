package encryption

// Response represents the successful decryption response
type Response struct {
	Data []struct {
		Token   string `json:"token"`
		Content string `json:"content"`
	} `json:"data"`
}