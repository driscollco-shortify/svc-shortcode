package interfaces

import serviceComponents "github.com/driscollco-core/service/components"

//go:generate go run go.uber.org/mock/mockgen -destination=../mocks/mock-request.go -package=mocks . Request
type Request interface {
	serviceComponents.Request
}
