package controller

type Systemd interface {
	Start(unit string) error
	Stop(unit string) error
	IsActive(unit string) bool
}

type Notifier interface {
	Send(event string, activeProver string) error
}
