package logTraffic

import (
	"github.com/driscollco-core/flow-client"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
)

func Action(ip string, parentSpan grafana.Span, logger log.Log) {
	span := parentSpan.Child("logging hit with flowControl")
	rateLimiter := flow.New(conf.Config.Service.Handlers.Create.FlowClient.TargetId)
	if err := rateLimiter.Hit(ip); err != nil {
		span.Error("error recording hit with flowControl")
		logger.Error("error recording hit with flowControl", "error", err.Error())
		return
	}
	span.Success()
}
