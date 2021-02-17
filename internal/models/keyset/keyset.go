package keysetmodel

type EncryptedKeyset struct {
	EncryptedKeyset string     `json:"encryptedKeyset"`
	KeysetInfo      KeysetInfo `json:"keysetInfo"`
}

type KeysetInfo struct {
	PrimaryKeyId int       `json:"primaryKeyId"`
	KeyInfo      []KeyInfo `json:"keyInfo"`
}

type KeyInfo struct {
	TypeUrl          string `json:"typeUrl"`
	Status           string `json:"status"`
	KeyID            int    `json:"keyId"`
	OutputPrefixType string `json:"outputPrefixType"`
}
