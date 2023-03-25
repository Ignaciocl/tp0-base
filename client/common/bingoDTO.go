package common

type BingoDTO struct {
	Name string `json:"name"`
	Document int `json:"document"`
	BornDate string `json:"born_date"`
	Number int `json:"number"`
	Surname string `json:"surname"`
}

type BingoCommunication struct {
	Action string `json:"action"`
	Data []BingoDTO `json:"data,omitempty"`
}

type BingoResponse struct {
	AmountProcessed int `json:"amount_processed"`
	Status string `json:"status"`
	Winners []string `json:"winners,omitempty"`
}
