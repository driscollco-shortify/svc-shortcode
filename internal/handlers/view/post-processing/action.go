package postProcessing

import (
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	cacheShortcode "github.com/driscollco-shortify/svc-shortcode/internal/handlers/view/post-processing/cache-shortcode"
	extendExpiry "github.com/driscollco-shortify/svc-shortcode/internal/handlers/view/post-processing/extend-expiry"
	"github.com/driscollco-shortify/svc-shortcode/internal/handlers/view/post-processing/grant-extra-clicks"
	"github.com/driscollco-shortify/svc-shortcode/internal/handlers/view/post-processing/reduce-clicks"
	sendMetric "github.com/driscollco-shortify/svc-shortcode/internal/handlers/view/post-processing/send-metric"
)

const (
	postProcessSpanTitle = "post processing (goroutines)"
)

func Action(bundle serviceComponents.Bundle, logger log.Log, shortCode entities.ShortCode, doCache bool) {
	postProcessSpan := bundle.Span(postProcessSpanTitle)
	defer postProcessSpan.Success()

	if doCache {
		go func() {
			cacheShortcode.Action(shortCode, postProcessSpan, bundle)
		}()
	}

	go func() {
		reduceClicks.Action(postProcessSpan, bundle.Db(), logger, shortCode)
	}()

	go func() {
		grantExtraClicks.Action(bundle, logger, shortCode)
	}()

	go func() {
		extendExpiry.Action(shortCode, postProcessSpan, bundle.Db(), logger)
	}()

	go func() {
		sendMetric.Action(postProcessSpan, bundle, logger)
	}()
}
