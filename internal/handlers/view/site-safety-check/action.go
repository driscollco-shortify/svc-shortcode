package siteSafetyCheck

import (
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
	flagUnsafe "github.com/driscollco-shortify/svc-shortcode/internal/handlers/view/site-safety-check/flag-unsafe"
	recordSafetyCheck "github.com/driscollco-shortify/svc-shortcode/internal/handlers/view/site-safety-check/record-safety-check"
	safeUrlChecker "github.com/driscollco-shortify/svc-shortcode/internal/safe-url-checker"
	"time"
)

func IsSiteSafe(shortCode entities.ShortCode, bundle serviceComponents.Bundle, logger log.Log) bool {
	if shortCode.Security.LastSafetyChecked.Time().Add(conf.Config.Behaviours.SafeUrlChecks.CheckInterval).After(time.Now()) {
		return true
	}

	span := bundle.Span("carrying out url safety checks")
	defer span.Success()

	childSpan := span.Child("carrying out safety check on destination url")
	isUnsafe, err := safeUrlChecker.IsUnsafe(shortCode.Get())
	if err != nil {
		logger.Error("error checking if shortCode destination is safe", "error", err.Error())
		childSpan.Error("error checking if site is safe")
		recordSafetyCheck.Action(shortCode, bundle.Db(), logger, span)
		return true
	}
	childSpan.Success()

	if isUnsafe {
		flagUnsafe.Action(shortCode, bundle, logger, span)
		return false
	}

	childSpan = span.Child("carrying out safety check on shortCode url")
	isUnsafe, err = safeUrlChecker.IsUnsafe(shortCode.ShortCode.URL)
	if err != nil {
		logger.Error("error checking if shortCode url is safe", "error", err.Error())
		childSpan.Error("error checking if site is safe")
		recordSafetyCheck.Action(shortCode, bundle.Db(), logger, span)
		return true
	}
	childSpan.Success()

	if isUnsafe {
		flagUnsafe.Action(shortCode, bundle, logger, span)
		return false
	}

	go func() {
		recordSafetyCheck.Action(shortCode, bundle.Db(), logger, span)
	}()
	return true
}
