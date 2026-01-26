package systemd

import (
	"log"
	"os"
	"os/exec"
)

type Manager struct {
	DryRun bool
}

func NewManager() Manager {
	return Manager{DryRun: os.Getenv("DRY_RUN") == "1"}
}

func (m Manager) Start(unit string) error {
	if m.DryRun {
		log.Printf("[DRY_RUN] systemctl start %s", unit)
		return nil
	}
	cmd := exec.Command("systemctl", "start", unit)
	_, err := cmd.CombinedOutput()
	return err
}

func (m Manager) Stop(unit string) error {
	if m.DryRun {
		log.Printf("[DRY_RUN] systemctl stop %s", unit)
		return nil
	}
	cmd := exec.Command("systemctl", "stop", unit)
	_, err := cmd.CombinedOutput()
	return err
}

func (m Manager) IsActive(unit string) bool {
	if m.DryRun {
		// In dry-run mode we’ll pretend it’s active to avoid flapping logic.
		// (Later we can simulate state if needed.)
		return true
	}
	cmd := exec.Command("systemctl", "is-active", "--quiet", unit)
	return cmd.Run() == nil
}
