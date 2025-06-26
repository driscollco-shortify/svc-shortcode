package generateShortCode

import (
	fireStore "github.com/driscollco-core/firestore"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
)

func Action(bundle serviceComponents.Bundle, shortCode *entities.ShortCode, logger log.Log, parentSpan grafana.Span) error {
	var err error
	exists := true
	for exists {
		if err = generateShortCode(shortCode, parentSpan, logger); err != nil {
			return err
		}

		exists, err = checkIfShortcodeExists(shortCode, parentSpan, logger, bundle.Db())
		if err != nil {
			return err
		}
	}
	return nil
}

func generateShortCode(shortCode *entities.ShortCode, parentSpan grafana.Span, logger log.Log) error {
	span := parentSpan.Child("generating a new shortCode")
	if err := shortCode.ChangeShortCode(conf.Config.Service.Handlers.Create.ShortCode.Length, conf.Config.Service.Handlers.Create.Domains); err != nil {
		span.Error("error generating shortCode")
		logger.Error("error generating shortCode", "error", err.Error())
		return err
	}
	span.Success()
	return nil
}

func checkIfShortcodeExists(shortCode *entities.ShortCode, parentSpan grafana.Span, logger log.Log, db fireStore.Client) (bool, error) {
	span := parentSpan.Child("checking if new shortCode is already assigned")
	results, err := db.Search(conf.Config.GCP.FireStore.Paths.ShortCodes,
		fireStore.Query().Where("ShortCode.URL", "==", shortCode.ShortCode.URL))
	if err != nil {
		logger.Error("error searching for existing shortCode in database", "error", err.Error())
		span.Error("error searching for existing shortCode in database")
		return false, err
	}
	span.Success()
	return len(results) > 0, nil
}
