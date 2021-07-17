package user

import "strconv"

func Inputuser(u map[string]string) {
	uid, _ = strconv.Atoi(u[Uid])
	name = u[Name]
	tel, _ = strconv.Atoi(u[Tel])
	addr = u[Addr]
	Infouser()
}
