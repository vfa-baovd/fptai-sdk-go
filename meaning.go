package fptai

type Meaning struct {
	Intents []struct{
		Name string `json:"label"`
		Confidence float64
	}
}