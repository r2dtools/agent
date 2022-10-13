package system

import (
	"fmt"
	"syscall"
)

type Privilege struct {
	uid, gid, euid, egid int
	initilised           bool
}

var privilige *Privilege

// IncreasePrivilege increases privileges for the current process
func (p *Privilege) IncreasePrivilege() error {
	if err := syscall.Seteuid(p.uid); err != nil {
		return fmt.Errorf("could not increase privilege (uid): %v", err)
	}

	return nil
}

// DropPrivilege drops privileges for the current process
func (p *Privilege) DropPrivilege() error {
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
	p.gid = syscall.Getgid()
	p.euid = syscall.Geteuid()
	p.egid = syscall.Getegid()
	p.initilised = true
}

// GetPrivilege returns privilege structure
func GetPrivilege() *Privilege {
	if privilige == nil {
		privilige = &Privilege{}
	}

	return privilige
}
