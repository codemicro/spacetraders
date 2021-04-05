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
	URLSubmitSellOrder       = func(username string) string { return URLBase + "users/" + username + "/sell-orders" }

	URLSubmitFlightplan         = func(username string) string { return URLBase + "users/" + username + "/flight-plans" }
	URLGetFlightplanInformation = func(username, flightplanID string) string {
		return URLBase + "users/" + username + "/flight-plans/" + flightplanID
	}

	URLGetShipInfo = func(username, shipID string) string { return URLBase + "users/" + username + "/ships/" + shipID}
)
