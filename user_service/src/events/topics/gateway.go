package topics

import (
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

func NewGateWayTopic() impls.Topic {
	gatewayInit.Do(func() {
		gatewayInstance = &gatewayTopic{
			gatewayService: gateway.NewService(),
		}
	})
	return gatewayInstance
}

func (g *gatewayTopic) ExecuteEvent(correlationId, key string, value []byte) error {
	var err error
	var data interface{}

	switch key {
	case "SaveUser":
		if data, err = impls.Marshal[dto.JwtTokenInfo](value); err == nil {
			// 근데 이거 goroutine으로 돌아야 할 것 같은데..
			// 그러면 error를 알 수 가 없음 어차피 return 안해도 logic 에러만 체크하면 되지 않나?
			// context나 channel로 받아서 문제 있으면
			// 근데 그 이전에 event에서 goroutine을 돌려야하나???
			// 그냥 다 때려치우고 원래는 transaction outbox pattern으로 가야하는데 그건 나중에 migration하고
			// 지금은 그냥 goroutine으로 돌려서 안에서 error 던지면 producer로 error 쏘는걸로
			err = g.gatewayService.SaveUser(data.(dto.JwtTokenInfo))
		}
	case "ChangeUserName":
		// go g.ChangeUserName(context.Background(), value)
		// g.ChangeUserName(value.(string))
	case "ChangeUserTag":
		// go g.ChangeUserTag(context.Background(), value)
		// g.ChangeUserTag(value.(string))
	case "ChangeProfile":
	default:
		// return fmt.Errorf("")
	}

	// error 대해서 log만 찍도록 해야할 듯
	// producer로 error 날리는것 까지 해야하나?
	// 원래는 transaction outbox pattern으로 처리해야함
	if err != nil {
		// logrus.WithFields(logrus.Fields{
		// 	"correlationId": correlationId,
		// 	"key":           key,
		// }).Error(err)
		return err
	} else {
		return nil
	}
}
