package encryption

// Response represents the successful decryption response
// type Response struct {
// 	Data []struct {
// 		Token   string `json:"token"`
// 		Content string `json:"content"`
// 	} `json:"data"`
// }

type Data struct {
	Token   string
	Content string
}

type Response struct {
	Data []Data
}

// var response encryption.Response
