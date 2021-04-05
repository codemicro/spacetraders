package control

import (
	"github.com/codemicro/spacetraders/internal/analysis"
	"github.com/codemicro/spacetraders/internal/stapi"
)

// This file contains ship controller functionality that scouts for locations without a known marketplace

func (s *ShipController) doScout() error {
	return nil
}

func (s *ShipController) grabMarketplaceData() error {
	marketplace, err := stapi.GetMarketplaceAtLocation(s.ship.Location)
	if err != nil {
		return err
	}
	return analysis.RecordMarketplaceAtLocation(s.ship.Location, marketplace)
}