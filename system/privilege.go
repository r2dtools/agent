package system

import (
	"fmt"
	"syscall"

	"github.com/r2dtools/agent/logger"
)

type Privilege struct {
	uid, gid   int
	initilised bool
}

var privilige *Privilege

// IncreasePrivilege increases privileges for the current process
func (p *Privilege) IncreasePrivilege() error {
	euid := syscall.Geteuid()

	if euid != syscall.Getuid() {
		logger.Debug(fmt.Sprintf("increse privilege: current EUID: %d", euid))

		if err := syscall.Setuid(euid); err != nil {
			return fmt.Errorf("could not increase privilege (uid): %v", err)
		}
	}

	egid := syscall.Getegid()

	if egid != syscall.Getgid() {
		logger.Debug(fmt.Sprintf("increse privilege: current EGID: %d", egid))

		if err := syscall.Setgid(egid); err != nil {
			syscall.Setuid(p.uid) // try to rollback uid
			return fmt.Errorf("could not increase privilege (gid): %v", err)
		}
	}

	return nil
}

// DropPrivilege drops privileges for the current process
func (p *Privilege) DropPrivilege() error {
	if syscall.Getuid() != p.uid {
		logger.Debug(fmt.Sprintf("drop privilege: current UID: %d", p.uid))

		if err := syscall.Setuid(p.uid); err != nil {
			return fmt.Errorf("could not drop privilege: %v", err)
		}
	}

	if syscall.Getgid() != p.gid {
		logger.Debug(fmt.Sprintf("drop privilege: current GID: %d", p.gid))

		if err := syscall.Setgid(p.gid); err != nil {
			return fmt.Errorf("could not drop privilege: %v", err)
		}
	}

	return nil
}

// Init initialises privilege object
func (p *Privilege) Init() {
	if p.initilised {
		return
	}

	p.uid = syscall.Getuid()
	logger.Debug(fmt.Sprintf("current UID: %d", p.uid))
	p.gid = syscall.Getgid()
	logger.Debug(fmt.Sprintf("current GID: %d", p.gid))
	p.initilised = true
}

// GetPrivilege returns privilege structure
func GetPrivilege() *Privilege {
	if privilige == nil {
		privilige = &Privilege{}
	}

	return privilige
}
