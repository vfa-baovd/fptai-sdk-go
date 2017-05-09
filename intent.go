package fptai

type IntentResponse struct {
	Name           string `json:"label"`
	Description     string `json:"description"`
	Code            string `json:"intent_code"`
	CreatedTime     string `json:"created_time"`
}

type Intent struct {
	Name string `json:"label"`
	Confidence float64
}