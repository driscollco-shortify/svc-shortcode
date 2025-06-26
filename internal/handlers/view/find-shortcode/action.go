package findShortcode

import (
	"errors"
	fireStore "github.com/driscollco-core/firestore"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
)

func Action(rawCode string, bundle serviceComponents.Bundle, logger log.Log) (entities.ShortCode, error) {
	span := bundle.Span("searching for shortCode")

	for _, path := range []string{
		conf.Config.GCP.FireStore.Paths.ShortCodes,
		conf.Config.GCP.FireStore.Paths.ShortCodesLegacy,
	} {
		var childSpan grafana.Span
		switch path {
		case conf.Config.GCP.FireStore.Paths.ShortCodes:
			childSpan = span.Child("searching database")
		default:
			childSpan = span.Child("searching legacy database")
		}
		results, err := bundle.Db().Search(path, fireStore.Query().Where("ShortCode.URL", "==", rawCode))
		if err != nil {
			logger.Error("error searching for shortCode", "error", err.Error())
			childSpan.Error("error searching for shortCode")
			span.Error("error searching for shortCode")
			return entities.ShortCode{}, err
		}

		shortCode := entities.ShortCode{}
		if len(results) > 0 {
			if err = results[0].DataTo(&shortCode); err != nil {
				logger.Error("error unmarshalling db record into shortCode", "error", err.Error())
				childSpan.Error("error unmarshalling db record")
				span.Error("error unmarshalling db record")
				return entities.ShortCode{}, err
			}
			childSpan.Success()
			span.Success()
			return shortCode, nil
		}
		childSpan.Success()
	}

	span.SuccessWithMsg("could not find shortCode")
	return entities.ShortCode{}, errors.New("not_found")
}
