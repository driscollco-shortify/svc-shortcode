package postProcessing

import (
	"encoding/json"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
	hasCapacity "github.com/driscollco-shortify/svc-shortcode/internal/handlers/create/flowcontrol/has-capacity"
	checkCreatorSafety "github.com/driscollco-shortify/svc-shortcode/internal/handlers/create/post-processing/check-creator-safety"
	writeShortcodeToDb "github.com/driscollco-shortify/svc-shortcode/internal/handlers/create/post-processing/write-shortcode-to-db"
	"time"
)

func Action(shortCode entities.ShortCode, ip string, bundle serviceComponents.Bundle, logger log.Log, postCreate bool) {
	span := bundle.Span("post processing")
	defer span.Success()

	if postCreate {
		go func() {
			postCreateMetric(bundle, span, logger)
		}()

		if !checkCreatorSafety.IsSafe(bundle, logger, span, ip) {
			logger.Info("creator ip was found to be unsafe; aborting shortCode save to db")
			return
		}

		go func() {
			writeShortcodeToDb.Action(&shortCode, bundle, logger, span)
		}()

		go func() {
			sendPubSubEvent(shortCode, bundle, span, logger)
		}()
	}

	go func(ipAddr string) {
		flowControlSpan := span.Child("managing traffic with flowControl")
		doesHaveCapacity := hasCapacity.Action(logger, flowControlSpan, ipAddr)
		flowControlSpan.Success()
		bundle.Cache().Namespace("flow-control").Set(doesHaveCapacity, time.Minute, "has-capacity")
	}(ip)
}

func postCreateMetric(bundle serviceComponents.Bundle, parentSpan grafana.Span, logger log.Log) {
	span := parentSpan.Child("sending post create metric to grafana")
	if err := bundle.Metric("shortify-metrics").Label("metric", "shortcodes_created").Send(1); err != nil {
		span.Error("error sending metric to grafana")
		logger.Error("error sending metric to grafana", "error", err.Error())
		return
	}
	span.Success()
}

func sendPubSubEvent(shortCode entities.ShortCode, bundle serviceComponents.Bundle, parentSpan grafana.Span, logger log.Log) {
	span := parentSpan.Child("marshalling shortCode to json")
	bytes, _ := json.Marshal(shortCode)
	span.Success()

	span = parentSpan.Child("sending pub/sub event to announce creation of shortCode")
	if _, err := bundle.PubSubTopic(conf.Config.Service.Handlers.Create.PubSub.Topics.Created).Publish(bytes); err != nil {
		span.Error("error sending pub/sub message")
		logger.Error("error sending pub/sub message", "error", err.Error())
		return
	}
	span.Success()
}
