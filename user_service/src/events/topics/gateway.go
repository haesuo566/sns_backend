package topics

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/haesuo566/sns_backend/user_service/src/events/impls"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/domain/gateway"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/dto"
)

type gatewayTopic struct {
	gatewayService *gateway.Service
}

var gatewayInit sync.Once
var gatewayInstance impls.Topic

func NewGateWayTopic(gatewayService *gateway.Service) impls.Topic {
	gatewayInit.Do(func() {
		gatewayInstance = &gatewayTopic{
			gatewayService,
		}
	})
	return gatewayInstance
}

func (g *gatewayTopic) ExecuteEvent(correlationId, key string, value []byte) error {
	var err error

	switch key {
	case "SaveUser":
		var data dto.JwtTokenInfo
		if err := json.Unmarshal(value, &data); err != nil {
			return err
		}
		err = g.gatewayService.SaveUser(data)
	case "ChangeUserName":
		// go g.ChangeUserName(context.Background(), value)
		// g.ChangeUserName(value.(string))
	case "ChangeUserTag":
		// go g.ChangeUserTag(context.Background(), value)
		// g.ChangeUserTag(value.(string))
	default:
		// return fmt.Errorf("")
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}
