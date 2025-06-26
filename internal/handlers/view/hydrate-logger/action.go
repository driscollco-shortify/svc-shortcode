package hydrateLogger

import (
	"github.com/driscollco-core/log"
	"github.com/driscollco-shortify/entities"
)

func Action(logger log.Log, shortCode entities.ShortCode) log.Log {
	return logger.Child("shortCode", shortCode.ShortCode.URL, "url", shortCode.Get(),
		"remaining clicks", shortCode.Clicks.Remaining-1, "clicks total", shortCode.Clicks.Total+1,
		"expiry", shortCode.Timeline.Expiry.ISO())
}
