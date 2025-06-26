package translateRequest

import (
	"encoding/json"
	"errors"
	"github.com/driscollco-core/grafana"
	"github.com/driscollco-core/log"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/entities/requests"
)

func Action(body []byte, ip string, logger log.Log, parentSpan grafana.Span) (entities.ShortCode, error) {
	if len(body) < 1 {
		return entities.ShortCode{}, errors.New("no request body")
	}

	span := parentSpan.Child("converting request body to request struct")
	req := requests.Create{}
	if err := json.Unmarshal(body, &req); err != nil {
		logger.Error("error unmarshalling request body", "error", err.Error())
		span.Error("error converting request body")
		return entities.ShortCode{}, err
	}
	span.Success()

	if len(req.Url) < 1 {
		return entities.ShortCode{}, errors.New("no url specified")
	}

	span = parentSpan.Child("translating create struct to shortCode")
	shortCode := entities.ShortCode{}
	if err := shortCode.Create(req); err != nil {
		logger.Error("error converting create request to shortCode", "error", err.Error())
		span.Error("error converting create request to shortCode")
		return entities.ShortCode{}, err
	}
	span.Success()
	shortCode.Security.CreatorIp = ip
	return shortCode, nil
}
