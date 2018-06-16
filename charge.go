package wechat

type Charge interface {
	Handle(conf map[string]interface{})(interface{}, error)
}