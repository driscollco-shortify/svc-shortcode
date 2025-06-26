package deleteExpired

import (
	"errors"
	"fmt"
	fireStore "github.com/driscollco-core/firestore"
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
)

func Action(shortCode entities.ShortCode, bundle serviceComponents.Bundle, logger log.Log) {
	if !shortCode.IsExpired() {
		return
	}

	span := bundle.Span("deleting shortCode because it is expired")
	for _, path := range []string{conf.Config.GCP.FireStore.Paths.ShortCodes, conf.Config.GCP.FireStore.Paths.ShortCodesLegacy} {
		if err := bundle.Db().Delete(fmt.Sprintf("%s/%s", path, shortCode.Id)); err != nil {
			if !errors.Is(err, fireStore.ErrorNoResults) {
				logger.Error("error deleting expired shortCode", "error", err.Error())
				span.Error(fmt.Sprintf("error deleting shortCode : %s", err.Error()))
				return
			}
		} else {
			span.Success()
			return
		}
	}
	logger.Error("could not find shortCode for deletion")
	span.Success()
}
