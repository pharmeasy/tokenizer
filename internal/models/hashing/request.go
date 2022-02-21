package hashing

type GenerateHashRequest struct {
	RequestId  string `json:"requestId"`
	Identifier string `json:"identifier"`
	Data       []Data `json:"data"`
}

type Data struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}
