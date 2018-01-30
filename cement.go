package cement

import (
	"github.com/desertbit/glue"
	"encoding/json"
)

const (
	MsgInvalidPayload = `{"error":"INVALID_PAYLOAD"}`

	CodeOk    = 0
	CodeError = 1
)

type payload struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

type response struct {
	Id     string `json:"id"`
	Status int    `json:"status"`
	Data   string `json:"data"`
}

type OnReadFunc func(channel *glue.Socket, messageId string, data string) (int, string)

func Glue(channel *glue.Channel, f OnReadFunc) glue.OnReadFunc {
	return func(data string) {
		p := &payload{}
		err := json.Unmarshal([]byte(data), p)
		if err != nil {
			writeJson(channel, &response{
				Id:     "",
				Status: CodeError,
				Data:   MsgInvalidPayload,
			})
			return
		}

		status, res := f(channel.Socket(), p.Id, p.Data)
		writeJson(channel, &response{
			Id:     p.Id,
			Status: status,
			Data:   res,
		})
	}
}

func writeJson(channel *glue.Channel, res *response) {
	b, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	channel.Write(string(b))
}
