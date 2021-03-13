package system

import (
	"fmt"
	"syscall"

	"github.com/r2dtools/agent/logger"
)

type Privilege struct {
	uid, gid, euid, egid int
	initilised           bool
}

var privilige *Privilege

// IncreasePrivilege increases privileges for the current process
func (p *Privilege) IncreasePrivilege() error {
	logger.Debug(fmt.Sprintf("increase privilege: EUID: %d", p.uid))

	if err := syscall.Seteuid(p.uid); err != nil {
		return fmt.Errorf("could not increase privilege (uid): %v", err)
	}

	return nil
}

// DropPrivilege drops privileges for the current process
func (p *Privilege) DropPrivilege() error {
	logger.Debug(fmt.Sprintf("drop privilege: EUID: %d", p.euid))

	if err := syscall.Seteuid(p.euid); err != nil {
		return fmt.Errorf("could not drop privilege: %v", err)
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
	p.euid = syscall.Geteuid()
	logger.Debug(fmt.Sprintf("current EUID: %d", p.euid))
	p.egid = syscall.Getegid()
	logger.Debug(fmt.Sprintf("current EGID: %d", p.egid))
	p.initilised = true
}

// GetPrivilege returns privilege structure
func GetPrivilege() *Privilege {
	if privilige == nil {
		privilige = &Privilege{}
	}

	return privilige
}
