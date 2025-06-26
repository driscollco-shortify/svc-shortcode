package hasCapacity

import (
	"github.com/driscollco-core/flow-client"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
	logTraffic "github.com/driscollco-shortify/svc-shortcode/internal/handlers/create/flowcontrol/log-traffic"
)

func Action(logger log.Log, parentSpan grafana.Span, ip string) bool {
	span := parentSpan.Child("checking capacity with flowControl")
	rateLimiter := flow.New(conf.Config.Service.Handlers.Create.FlowClient.TargetId)
	hasCapacity := rateLimiter.HasCapacity(ip, conf.Config.Service.Handlers.Create.FlowClient.RateLimitId)
	span.Success()

	go func() {
		logTraffic.Action(ip, parentSpan, logger)
	}()

	if hasCapacity {
		return true
	}

	go func() {
		innerSpan := parentSpan.Child("checking hard capacity with flowControl")
		hasHardCapacity := rateLimiter.HasCapacity(ip, conf.Config.Service.Handlers.Create.FlowClient.HardRateLimitId)
		innerSpan.Success()

		if hasHardCapacity {
			return
		}

		innerSpan = parentSpan.Child("banning ip address with CloudFlare")
		if err := rateLimiter.BanIp(ip, conf.Config.Service.Handlers.Create.FlowClient.TargetId, conf.Config.CloudFlare.Credentials.ApiKey, conf.Config.CloudFlare.Credentials.ZoneId); err != nil {
			innerSpan.Error("error banning ip with CloudFlare")
			logger.Error("error banning ip address with cloudflare", "error", err.Error(), "ip", ip)
			return
		}
		innerSpan.Success()
	}()
	return false
}
