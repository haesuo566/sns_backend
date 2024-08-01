package topics

import "sync"

type GatewayTopic struct {
}

var gatewayInit sync.Once
var gatewayInstance *GatewayTopic

func NewGateWayTopic() *GatewayTopic {
	gatewayInit.Do(func() {
		gatewayInstance = &GatewayTopic{}
	})
	return gatewayInstance
}

func (g *GatewayTopic) ExecuteEvent(key, value string) {
	switch key {
	case "ChangeUserName":
		// go g.ChangeUserName(context.Background(), value)
		g.ChangeUserName(value)
	case "ChangeUserTag":
		// go g.ChangeUserTag(context.Background(), value)
		g.ChangeUserTag(value)
	default:
	}
}

// 유저 이름 변경
func (g *GatewayTopic) ChangeUserName(value string) error {
	return nil
}

// 유저 아이디(태그) 변경
func (g *GatewayTopic) ChangeUserTag(value string) error {
	return nil
}
