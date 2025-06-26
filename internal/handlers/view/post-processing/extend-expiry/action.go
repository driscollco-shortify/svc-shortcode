package extendExpiry

import (
	fireStore "github.com/driscollco-core/firestore"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	"github.com/driscollco-shortify/entities"
	dbShortcodeUpdate "github.com/driscollco-shortify/svc-shortcode/internal/db-shortcode-update"
	"time"
)

func Action(shortCode entities.ShortCode, parentSpan grafana.Span, db fireStore.Client, logger log.Log) {
	if shortCode.Timeline.Expiry.Time().After(time.Now().AddDate(0, 0, 20)) {
		return
	}

	span := parentSpan.Child("extending the expiry date of the shortcode").
		Attribute("remaining clicks", shortCode.Clicks.Remaining).
		Attribute("expiry", shortCode.Timeline.Expiry.ISO())

	if err := dbShortcodeUpdate.Action(shortCode, db, map[string]fireStore.UpdateValue{
		"Timeline.Expiry":  {Value: time.Now().AddDate(0, 1, 0).Unix()},
		"Version":          {Increment: 1},
		"Timeline.Updated": {Value: time.Now().Unix()},
	}); err != nil {
		logger.Error("error updating shortCode record in fireStore", "error", err.Error())
		span.Error("error updating shortCode in fireStore")
		return
	}
	span.Success()
}
