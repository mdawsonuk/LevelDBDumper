// +build windows

package main

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func isAdmin() bool {

	var sid *windows.SID

	err := windows.AllocateAndInitializeSid(&windows.SECURITY_NT_AUTHORITY, 2, windows.SECURITY_BUILTIN_DOMAIN_RID, windows.DOMAIN_ALIAS_RID_ADMINS, 0, 0, 0, 0, 0, 0, &sid)
	if err != nil {
		printLine(fmt.Sprintf("SID Error: %s", err), Fatal)
		return false
	}

	token := windows.Token(0)

	member, err := token.IsMember(sid)
	if err != nil {
		printLine(fmt.Sprintf("Token Membership Error: %s", err), Fatal)
		return false
	}
	return member
}
