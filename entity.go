package sdk

type Entity struct {
	Type string `json:"type"`
	Value string `json:"value"`
	Confidence string `json:"confidence"` // WTF, string?
}

type EntityResponse struct {
	Total int `json:"total"`
	Data []Entity `json:"data"`
}