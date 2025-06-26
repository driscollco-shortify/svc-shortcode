package handlerView

import (
	"fmt"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	deleteExpired "github.com/driscollco-shortify/svc-shortcode/internal/handlers/view/delete-expired"
	findShortcode "github.com/driscollco-shortify/svc-shortcode/internal/handlers/view/find-shortcode"
	hydrateLogger "github.com/driscollco-shortify/svc-shortcode/internal/handlers/view/hydrate-logger"
	postProcessing "github.com/driscollco-shortify/svc-shortcode/internal/handlers/view/post-processing"
	siteSafetyCheck "github.com/driscollco-shortify/svc-shortcode/internal/handlers/view/site-safety-check"
	"net/http"
)

func Handle(bundle serviceComponents.Bundle, c serviceComponents.Request) error {
	logger := bundle.Log()
	if c.Params("shortCode") == "" {
		return c.Redirect("https://shortify.pro", http.StatusPermanentRedirect)
	}

	if len(c.Params("shortCode")) > 3 {
		span := bundle.Span("redirecting non-shortcode back to shortify.pro").Attribute("url", c.Params("shortCode"))
		span.Success()
		return c.Redirect("https://shortify.pro", http.StatusPermanentRedirect)
	}

	potentialShortCode := fmt.Sprintf("%s/%s", c.Hostname(), c.Params("shortCode"))
	var shortCode entities.ShortCode

	bundle.Attribute("shortCode", potentialShortCode)
	span := bundle.Span("checking if shortCode is cached")
	item := bundle.Cache().Namespace("shortCodes").Get(potentialShortCode)
	span.Success()

	var ok bool
	shortCode, ok = item.(entities.ShortCode)
	if ok {
		bundle.Attribute("clicks.remaining", shortCode.Clicks.Remaining)
		bundle.Attribute("clicks.total", shortCode.Clicks.Total)
		bundle.Attribute("expiry", shortCode.Timeline.Expiry.ISO())
		postProcessing.Action(bundle, logger, shortCode, false)
		return c.Redirect(shortCode.RawOriginal, http.StatusPermanentRedirect)
	}

	shortCode, err := findShortcode.Action(potentialShortCode, bundle, logger)
	if err != nil {
		return c.Redirect("https://shortify.pro", http.StatusPermanentRedirect)
	}

	logger = hydrateLogger.Action(logger, shortCode)
	bundle.Attribute("clicks.remaining", shortCode.Clicks.Remaining)
	bundle.Attribute("clicks.total", shortCode.Clicks.Total)
	bundle.Attribute("expiry", shortCode.Timeline.Expiry.ISO())

	if !shortCode.CanView() {
		deleteExpired.Action(shortCode, bundle, logger)
		return c.Redirect("https://shortify.pro", http.StatusPermanentRedirect)
	}

	go func() {
		siteSafetyCheck.IsSiteSafe(shortCode, bundle, logger)
	}()

	postProcessing.Action(bundle, logger, shortCode, true)
	return c.Redirect(shortCode.RawOriginal, http.StatusPermanentRedirect)
}
