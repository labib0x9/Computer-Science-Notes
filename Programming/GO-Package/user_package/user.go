package main

import "os/user"

// Incomplete

func main() {
	User, _ := user.Current()

	User.GroupIds()
	_ = User.Username

	user.Lookup()
	user.LookupGroup()
	user.LookupGroupId()
	user.LookupId()

	user.UnknownGroupError
	user.UnknownGroupIdError
	user.UnknownUserError
	user.UnknownUserIdError
}
