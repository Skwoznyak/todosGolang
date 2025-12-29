package responce

//использовался бы во всех хендлерах, вынес отдельно
type Responce struct{
	Status string `json:"status"`
	Error string `json:"error,omitempty"` // теги не работаю в стандартной либе(, просто как пример
}

const (
	StatusOK = "OK"
	StatusError = "Error"
)

func OK() Responce {
	return Responce{
		Status: StatusOK,
	}
}

func Error(msg string) Responce{
	return Responce{
		Status: StatusError,
		Error: msg,
	}
}