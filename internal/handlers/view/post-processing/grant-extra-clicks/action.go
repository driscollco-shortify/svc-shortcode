package grantExtraClicks

import (
	fireStore "github.com/driscollco-core/firestore"
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
	dbShortcodeUpdate "github.com/driscollco-shortify/svc-shortcode/internal/db-shortcode-update"
	"time"
)

func Action(bundle serviceComponents.Bundle, logger log.Log, shortCode entities.ShortCode) {
	if shortCode.Clicks.Remaining >= conf.Config.Behaviours.ShortCodes.ExtendLifetime.RemainingClicks.TriggerPoint {
		return
	}

	span := bundle.Span("incrementing remaining clicks for shortCode")

	logger.Info("granting additional clicks as low point threshold has been met",
		"granted clicks", conf.Config.Behaviours.ShortCodes.ExtendLifetime.RemainingClicks.TopUp)

	if err := dbShortcodeUpdate.Action(shortCode, bundle.Db(), map[string]fireStore.UpdateValue{
		"Clicks.Remaining": {Increment: conf.Config.Behaviours.ShortCodes.ExtendLifetime.RemainingClicks.TopUp},
		"Version":          {Increment: 1},
		"Timeline.Updated": {Value: time.Now().Unix()},
	}); err != nil {
		logger.Error("error updating shortCode record in fireStore", "error", err.Error())
		span.Error("error updating shortCode in fireStore")
		return
	}
	span.Success()
}
