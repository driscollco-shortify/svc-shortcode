package cacheShortcode

import (
	"github.com/driscollco-core/grafana"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	"time"
)

func Action(shortCode entities.ShortCode, parentSpan grafana.Span, bundle serviceComponents.Bundle) {
	span := parentSpan.Child("caching shortCode")
	bundle.Cache().Namespace("shortCodes").Set(shortCode, time.Minute, shortCode.ShortCode.URL)
	span.Success()
}
