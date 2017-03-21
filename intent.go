package sdk

type Intent struct {
	Label string `json:"label"`
	Description string `json:"description"`
	Code string `json:"intent_code"`
	CreatedTime string `json:"created_time"`
	ApplicationCode string `json:"application_code"`
}