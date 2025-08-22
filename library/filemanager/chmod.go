package filemanager

import (
	"os"
	"strconv"
)

const (
	PermRead  = 4
	PermWrite = 2
	PermExec  = 1
)

type Perm struct {
	Read  bool
	Write bool
	Exec  bool
}

func (p *Perm) ToCodes() [3]uint32 {
	v := [3]uint32{0, 0, 0}
	if p.Read {
		v[0] = PermRead
	}
	if p.Write {
		v[1] = PermWrite
	}
	if p.Exec {
		v[2] = PermExec
	}
	return v
}

func toPerm(n byte) Perm {
	switch n {
	case '7':
		return Perm{Read: true, Write: true, Exec: true}
	case '6':
		return Perm{Read: true, Write: true, Exec: false}
	case '5':
		return Perm{Read: true, Write: false, Exec: true}
	case '4':
		return Perm{Read: true, Write: false, Exec: false}
	case '3':
		return Perm{Read: false, Write: true, Exec: true}
	case '2':
		return Perm{Read: false, Write: true, Exec: false}
	case '1':
		return Perm{Read: false, Write: false, Exec: true}
	default:
		return Perm{}
	}
}

func FileModeToPerms(mode os.FileMode) Perms {
	v := strconv.FormatUint(uint64(mode), 8)
	if len(v) == 3 {
		return Perms{
			Owner: toPerm(v[0]),
			Group: toPerm(v[1]),
			Other: toPerm(v[2]),
		}
	}
	return Perms{}
}

type Perms struct {
	Owner Perm
	Group Perm
	Other Perm
}

func ValidatePermCodes(n [3]uint32) bool {
	if n[0] != PermRead && n[0] != 0 {
		return false
	}
	if n[1] != PermWrite && n[1] != 0 {
		return false
	}
	if n[2] != PermExec && n[2] != 0 {
		return false
	}
	return true
}

func ValidatePermNumber(n uint32) bool {
	return n >= 1 && n <= 7
}
