package checkCreatorSafety

import (
	fireStore "github.com/driscollco-core/firestore"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
)

func IsSafe(bundle serviceComponents.Bundle, logger log.Log, parentSpan grafana.Span, ip string) bool {
	span := parentSpan.Child("checking if creator ip is safe")
	results, err := bundle.Db().Search(conf.Config.GCP.FireStore.Paths.ShortCodes,
		fireStore.Query().Where("Security.CreatorIp", "==", ip).
			Where("Security.FlaggedUnsafe", "==", true))

	if err != nil {
		span.Error("error checking if ip is safe")
		logger.Error("error checking if ip is safe in fireStore", "error", err.Error())
		return true
	}
	span.Success()
	if len(results) > 0 {
		logger.Info("creator ip address has been flagged as the creator of an unsafe url")
		return false
	}
	return true
}
