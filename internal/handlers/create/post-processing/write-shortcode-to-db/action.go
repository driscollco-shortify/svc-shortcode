package writeShortcodeToDb

import (
	"fmt"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
)

func Action(shortCode *entities.ShortCode, bundle serviceComponents.Bundle, logger log.Log, parentSpan grafana.Span) {
	span := parentSpan.Child("writing shortCode to database")
	if err := bundle.Db().Write(fmt.Sprintf("%s/%s", conf.Config.GCP.FireStore.Paths.ShortCodes, shortCode.Id), shortCode); err != nil {
		span.Error("error writing shortCode to database")
		logger.Error("error writing shortCode to database", "error", err.Error())
		return
	}
	span.Success()
}
