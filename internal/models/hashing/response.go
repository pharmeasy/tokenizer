package hashing

type GenerateHashResponse struct {
	Data []Response `json:"data"`
}

type Response struct {
	Id   string `json:"id"`
	Hash string `json:"hash"`
}
