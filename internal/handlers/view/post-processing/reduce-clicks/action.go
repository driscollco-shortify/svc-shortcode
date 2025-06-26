package reduceClicks

import (
	fireStore "github.com/driscollco-core/firestore"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	"github.com/driscollco-shortify/entities"
	dbShortcodeUpdate "github.com/driscollco-shortify/svc-shortcode/internal/db-shortcode-update"
	"time"
)

func Action(parentSpan grafana.Span, db fireStore.Client, logger log.Log, shortCode entities.ShortCode) {
	span := parentSpan.Child("decrementing the remaining clicks for shortcode").
		Attribute("remaining clicks", shortCode.Clicks.Remaining).
		Attribute("expiry", shortCode.Timeline.Expiry.ISO())

	if err := dbShortcodeUpdate.Action(shortCode, db, map[string]fireStore.UpdateValue{
		"Clicks.Total":     {Increment: 1},
		"Clicks.Remaining": {Increment: -1},
		"Version":          {Increment: 1},
		"Timeline.Updated": {Value: time.Now().Unix()},
	}); err != nil {
		logger.Error("error updating shortCode record in fireStore", "error", err.Error())
		span.Error("error updating shortCode in fireStore")
		return
	}
	span.Success()
}
