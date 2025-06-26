package sendMetric

import (
	"fmt"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	serviceComponents "github.com/driscollco-core/service/components"
)

func Action(parentSpan grafana.Span, bundle serviceComponents.Bundle, logger log.Log) {
	span := parentSpan.Child("sending metric to Grafana")
	if err := bundle.Metric("shortify-metrics").Label("metric", "shortcodes-visited").Send(1); err != nil {
		logger.Error("error sending metric to Grafana", "error", err.Error())
		span.Error(fmt.Sprintf("error sending metric to Grafana : %s", err.Error()))
		return
	}
	span.Success()
}
