package setupShortcode

import (
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
	generateShortCode "github.com/driscollco-shortify/svc-shortcode/internal/handlers/create/setup-shortcode/generate-short-code"
	"time"
)

func Action(shortCode *entities.ShortCode, bundle serviceComponents.Bundle, parentSpan grafana.Span, logger log.Log) error {
	span := parentSpan.Child("generating a new shortCode URL")
	if err := generateShortCode.Action(bundle, shortCode, logger, span); err != nil {
		span.Error("error generating a new shortCode URL")
		return err
	}
	span.Success()

	if !shortCode.Clicks.UserDefined {
		shortCode.Clicks.Total = 0
		shortCode.Clicks.Remaining = conf.Config.Service.Handlers.Create.Clicks.Max
	}

	if shortCode.Timeline.Expiry.Time().IsZero() {
		shortCode.Timeline.Expiry.Set(time.Now().Add(conf.Config.Service.Handlers.Create.Lifetime.Expiry))
	}
	return nil
}
