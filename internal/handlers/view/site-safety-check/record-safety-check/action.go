package recordSafetyCheck

import (
	fireStore "github.com/driscollco-core/firestore"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	"github.com/driscollco-shortify/entities"
	dbShortcodeUpdate "github.com/driscollco-shortify/svc-shortcode/internal/db-shortcode-update"
	"time"
)

func Action(shortCode entities.ShortCode, db fireStore.Client, logger log.Log, parentSpan grafana.Span) {
	span := parentSpan.Child("updating last safety check record in database")
	if err := dbShortcodeUpdate.Action(shortCode, db, map[string]fireStore.UpdateValue{
		"Security.LastSafetyChecked": {Value: time.Now().Unix()},
		"Timeline.Updated":           {Value: time.Now().Unix()},
		"Version":                    {Increment: 1},
	}); err != nil {
		logger.Error("error updating last safety check record in fireStore", "error", err.Error())
		span.Error("error updating last safety check in fireStore")
		return
	}
	span.Success()
}
