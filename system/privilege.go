package system

import (
	"fmt"

	"github.com/r2dtools/agent/logger"
	"golang.org/x/sys/unix"
)

type Privilege struct {
	uid, gid   int
	initilised bool
}

var privilige *Privilege

// IncreasePrivilege increases privileges for the current process
func (p *Privilege) IncreasePrivilege() error {
	euid := unix.Geteuid()
	logger.Debug(fmt.Sprintf("increse privilege: current EUID: %d", euid))

	if err := unix.Setuid(euid); err != nil {
		return fmt.Errorf("could not increase privilege: %v", err)
	}

	egid := unix.Getegid()
	logger.Debug(fmt.Sprintf("increse privilege: current EGID: %d", egid))

	if err := unix.Setgid(egid); err != nil {
		unix.Setuid(p.uid) // try to rollback uid
		return fmt.Errorf("could not increase privilege: %v", err)
	}

	return nil
}

// DropPrivilege drops privileges for the current process
func (p *Privilege) DropPrivilege() error {
	if unix.Getuid() != p.uid {
		if err := unix.Setuid(p.uid); err != nil {
			return fmt.Errorf("could not drop privilege: %v", err)
		}
	}

	if unix.Getgid() != p.gid {
		if err := unix.Setgid(p.gid); err != nil {
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

	p.uid = unix.Getuid()
	logger.Debug(fmt.Sprintf("current UID: %d", p.uid))
	p.gid = unix.Getgid()
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
