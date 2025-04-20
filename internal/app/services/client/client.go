package client

type Client struct {
	Id      string      `json:"id"`
	MsgChan chan string `json:"msg_chan"`
}

func NewClient(id string) *Client {
	return &Client{
		Id:      id,
		MsgChan: make(chan string),
	}
}
