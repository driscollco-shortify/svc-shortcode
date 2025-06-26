package flagUnsafe

import (
	"errors"
	"fmt"
	fireStore "github.com/driscollco-core/firestore"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
	dbShortcodeUpdate "github.com/driscollco-shortify/svc-shortcode/internal/db-shortcode-update"
	"time"
)

func Action(shortCode entities.ShortCode, bundle serviceComponents.Bundle, logger log.Log, parentSpan grafana.Span) {
	span := parentSpan.Child("flagging site as unsafe")
	defer span.Success()

	logger.Info("site is found to be unsafe")

	if len(shortCode.Security.CreatorIp) < 1 {
		deleteShortcode(shortCode, bundle.Db(), logger, span)
		return
	}

	dbSpan := parentSpan.Child("updating shortCode record in database")
	if err := dbShortcodeUpdate.Action(shortCode, bundle.Db(), map[string]fireStore.UpdateValue{
		"Clicks.Remaining": {Increment: conf.Config.Behaviours.ShortCodes.ExtendLifetime.RemainingClicks.TopUp},
		"Version":          {Increment: 1},
		"Timeline.Updated": {Value: time.Now().Unix()},
	}); err != nil {
		logger.Error("error updating shortCode", "error", err.Error())
		dbSpan.Error("error updating fireStore")
		return
	}
	dbSpan.Success()
	span.Success()
}

func deleteShortcode(shortCode entities.ShortCode, db fireStore.Client, logger log.Log, parentSpan grafana.Span) {
	span := parentSpan.Child("deleting shortCode from fireStore")
	if err := db.Delete(fmt.Sprintf("%s/%s", conf.Config.GCP.FireStore.Paths.ShortCodes, shortCode.Id)); err != nil {
		if !errors.Is(err, fireStore.ErrorNoResults) {
			logger.Error("error deleting shortCode from fireStore", "error", err.Error())
			span.Error("error deleting shortCode from fireStore")
			return
		}

		if err = db.Delete(fmt.Sprintf("%s/%s", conf.Config.GCP.FireStore.Paths.ShortCodesLegacy, shortCode.Id)); err != nil {
			if !errors.Is(err, fireStore.ErrorNoResults) {
				logger.Error("error deleting shortCode from fireStore", "error", err.Error())
				span.Error("error deleting shortCode from fireStore")
				return
			}
			logger.Error("could not find shortCode in database when trying to delete")
			span.Success()
			return
		}
	}
	span.Success()
}
