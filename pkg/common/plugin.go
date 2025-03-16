package common

type Plugin interface {
	OnLoad() error
}
