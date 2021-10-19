package user

import (
	"fmt"
)

func Infouser() {
	fmt.Printf("%s为 %d 的用户%s、%s和%s为: %s %d %s\n", Uid, uid, Name, Tel, Addr, name, tel, addr)
}

