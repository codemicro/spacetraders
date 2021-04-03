package stapi

// This file contains URLs for accessing various resources

const (
	URLBase = "https://api.spacetraders.io/"
)

var (
	URLUserInfo = func (username string) string { return URLBase + "users/" + username }
)
