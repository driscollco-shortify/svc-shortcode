package handlerCreate

import (
	serviceComponents "github.com/driscollco-core/service/components"
	existingShortcode "github.com/driscollco-shortify/svc-shortcode/internal/handlers/create/existing-shortcode"
	getJwtForShortCode "github.com/driscollco-shortify/svc-shortcode/internal/handlers/create/get-jwt-for-shortcode"
	postProcessing "github.com/driscollco-shortify/svc-shortcode/internal/handlers/create/post-processing"
	setupShortcode "github.com/driscollco-shortify/svc-shortcode/internal/handlers/create/setup-shortcode"
	translateRequest "github.com/driscollco-shortify/svc-shortcode/internal/handlers/create/translate-request"
	"net/http"
)

func Handle(bundle serviceComponents.Bundle, c serviceComponents.Request) error {
	logger := bundle.Log().Child("creator ip", c.IP())

	cacheItem := bundle.Cache().Namespace("flow-control").Get("has-capacity")
	doesHaveCapacity, ok := cacheItem.(bool)
	if ok {
		if !doesHaveCapacity {
			return c.Status(http.StatusTooManyRequests).SendString("slow down - rate limited")
		}
	}

	span := bundle.Span("processing create request")
	shortCode, err := translateRequest.Action(c.Body(), c.IP(), logger, span)
	if err != nil {
		span.Error("error processing create request")
		return c.Status(http.StatusBadRequest).SendString("you must supply a valid create request")
	}

	existing, err := existingShortcode.Get(shortCode, bundle, logger, span)
	if err != nil {
		span.Error("error checking if shortCode destination exists already")
		return c.Status(http.StatusFailedDependency).SendString("error verifying shortCode with backend storage")
	}
	span.Success()

	if len(existing.Id) > 0 {
		token, err := getJwtForShortCode.Action(bundle, existing, false)
		if err != nil {
			return c.Status(http.StatusFailedDependency).SendString("error verifying shortCode")
		}

		go func(ip string) {
			postProcessing.Action(shortCode, ip, bundle, logger, false)
		}(c.IP())

		return c.SendString(token)
	}

	span = bundle.Span("setting up shortCode")
	if err = setupShortcode.Action(&shortCode, bundle, span, logger); err != nil {
		span.Error("error setting up shortCode")
		return c.Status(http.StatusFailedDependency).SendString("error setting up shortCode")
	}
	span.Success()

	bundle.Attribute("shortCode", shortCode.ShortCode.URL)

	go func(ip string) {
		postProcessing.Action(shortCode, ip, bundle, logger, true)
	}(c.IP())

	token, err := getJwtForShortCode.Action(bundle, shortCode, false)
	if err != nil {
		return c.Status(http.StatusFailedDependency).SendString("processing error")
	}
	return c.SendString(token)
}
