package system

import (
	"fmt"
	"syscall"

	"github.com/r2dtools/agent/pkg/logger"
)

type PrivilegeManagerInterface interface {
	Increase() error
	Drop() error
}

type PrivilegeManager struct {
	uid, gid, euid, egid int
	logger               logger.LoggerInterface
}

func (p *PrivilegeManager) Increase() error {
	p.logger.Debug("increase sys privilege: euid: %d", p.uid)

	if err := syscall.Seteuid(p.uid); err != nil {
		return fmt.Errorf("could not increase sys privilege (uid): %v", err)
	}

	return nil
}

func (p *PrivilegeManager) Drop() error {
	p.logger.Debug("drop sys privilege: euid: %d", p.euid)

	if err := syscall.Seteuid(p.euid); err != nil {
		return fmt.Errorf("could not drop sys privilege: %v", err)
	}

	return nil
}

func GetPrivilegeManager(logger logger.LoggerInterface) PrivilegeManagerInterface {
	return &PrivilegeManager{
		logger: logger,
		uid:    syscall.Getuid(),
		gid:    syscall.Getgid(),
		euid:   syscall.Geteuid(),
		egid:   syscall.Getegid(),
	}
}
