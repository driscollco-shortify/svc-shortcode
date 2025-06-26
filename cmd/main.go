package main

import (
	fireStore "github.com/driscollco-core/firestore"
	gcpPubSub "github.com/driscollco-core/gcp-pub-sub"
	"github.com/driscollco-core/service"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
	handlerCreate "github.com/driscollco-shortify/svc-shortcode/internal/handlers/create"
	handlerRootRedirect "github.com/driscollco-shortify/svc-shortcode/internal/handlers/root-redirect"
	handlerView "github.com/driscollco-shortify/svc-shortcode/internal/handlers/view"
	"os"
)

func main() {
	s := service.New("shortCode")
	s.Router().Get("/", handlerRootRedirect.Handle, "redirectFromShortDomain")
	s.Router().Get("/:shortCode", handlerView.Handle, "viewShortcode")
	s.Router().Post("/create", handlerCreate.Handle, "createShortcode")

	if err := s.Config().Populate(&conf.Config); err != nil {
		s.Log().Error("unable to populate config", "error", err.Error())
		os.Exit(0)
	}

	s.Log().Info("test message", "grafana metrics userId", os.Getenv("Grafana_Metrics_UserId"))

	s.WithDb(fireStore.New(conf.Config.GCP.ProjectId, conf.Config.GCP.Credentials))

	sender, err := gcpPubSub.NewSender(conf.Config.GCP.ProjectId, conf.Config.Service.Handlers.Create.PubSub.Topics.Created)
	if err != nil {
		s.Log().Error("unable to create pub/sub sender for topic", "topic", conf.Config.Service.Handlers.Create.PubSub.Topics.Created, "error", err.Error())
		os.Exit(0)
	}
	s.WithPubSubTopic(conf.Config.Service.Handlers.Create.PubSub.Topics.Created, sender)
	s.Run()
}
