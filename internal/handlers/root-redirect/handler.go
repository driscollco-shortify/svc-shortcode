package handlerRootRedirect

import (
	serviceComponents "github.com/driscollco-core/service/components"
	"net/http"
)

func Handle(bundle serviceComponents.Bundle, c serviceComponents.Request) error {
	return c.Redirect("https://shortify.pro", http.StatusPermanentRedirect)
}
