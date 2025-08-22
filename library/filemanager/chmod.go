package filemanager

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
