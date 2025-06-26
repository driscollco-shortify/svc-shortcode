package existingShortcode

import (
	fireStore "github.com/driscollco-core/firestore"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
)

func Get(shortCode entities.ShortCode, bundle serviceComponents.Bundle, logger log.Log, parentSpan grafana.Span) (entities.ShortCode, error) {
	span := parentSpan.Child("checking database to see if shortCode exists already")
	results, err := bundle.Db().Search(conf.Config.GCP.FireStore.Paths.ShortCodes, fireStore.Query().Where("Hash", "==", shortCode.Hash))
	if err != nil {
		span.Error("error checking for existing shortCode")
		logger.Error("error checking for existing shortCode", "error", err.Error())
		return entities.ShortCode{}, err
	}
	span.Success()

	if len(results) < 1 {
		return entities.ShortCode{}, nil
	}

	span = parentSpan.Child("populating shortCode struct from db record")
	existing := entities.ShortCode{}
	if err = results[0].DataTo(&existing); err != nil {
		span.Error("error populating shortCode")
		return entities.ShortCode{}, err
	}
	span.Success()
	return existing, nil
}
