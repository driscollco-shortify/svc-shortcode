package interfaces

import serviceComponents "github.com/driscollco-core/service/components"

//go:generate go run go.uber.org/mock/mockgen -destination=../mocks/mock-bundle.go -package=mocks . Bundle
type Bundle interface {
	serviceComponents.Bundle
}
