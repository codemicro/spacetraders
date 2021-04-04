package stapi

// This file contains URLs for accessing various resources

const (
	URLBase = "https://api.spacetraders.io/"
)

var (
	URLUserInfo = func(username string) string { return URLBase + "users/" + username }

	URLSystemLocations     = func(system string) string { return URLBase + "game/systems/" + system + "/locations" }
	URLLocationInformation = func(location string) string { return URLBase + "game/locations/" + location }

	URLMarketplaceAtLocation = func(location string) string { return URLBase + "game/locations/" + location + "/marketplace" }
	URLSubmitPurchaseOrder   = func(username string) string { return URLBase + "users/" + username + "/purchase-orders" }
)
