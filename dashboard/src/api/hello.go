package api

type HelloModel struct {
	Key int
}

var _ = (Model)(&HelloModel{})

func (self *HelloModel) Access(appKey string, token string) bool {
	return true
}

type ResponseHello struct {
	Msg string `json:"msg"`
	Key int    `json:"key"`
}

func (self *HelloModel) GetApiHello(res *ResponseHello) error {
	res.Msg = "Hello, World"
	res.Key = self.Key
	self.Key++
	return nil
}

type RequestHello struct {
    Key int `json:"key"`
}

func (self *HelloModel) PostApiHello(res *ResponseHello, req *RequestHello) error {
    res.Key = req.Key
    self.Key = req.Key
    res.Msg = "ok"

    return nil
}
