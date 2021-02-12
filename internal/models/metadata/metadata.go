package metadata

// MetaUpdateRequest represents the struct for incoming update requests
type MetaUpdateRequest struct {
	RequestID    string         `json:"requestId"`
	Identifier   string         `json:"identifier"`
	Level        int            `json:"level"`
	UpdateParams []UpdateParams `json:"data"`
}

// MetaRequest represents the struct for incoming update requests
type MetaRequest struct {
	RequestID  string   `json:"requestId"`
	Identifier string   `json:"identifier"`
	Level      int      `json:"level"`
	Tokens     []string `json:"tokens"`
}

// MetaResponse represents the successful metadata response
type MetaResponse struct {
	MetaParams []MetaParams `json:"data"`
}

// MetaParams is the struct for object
type MetaParams struct {
	Token    string `json:"token"`
	Metadata string `json:"metadata"`
}

// UpdateParams contain the meta params to be updated
type UpdateParams struct {
	Token    string `json:"token"`
	Metadata string `json:"metadata"`
}

// MetaUpdateResponse represents the successful update response message for successful metadata updates
type MetaUpdateResponse struct {
	Status string
}
