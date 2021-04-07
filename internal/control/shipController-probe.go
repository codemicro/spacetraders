package control

import (
	"github.com/codemicro/spacetraders/internal/analysis"
	"github.com/codemicro/spacetraders/internal/stapi"
	"time"
)

func (s *ShipController) probeAction() {
	for {
		// s.log("recording marketplace at %s", s.ship.Location)
		marketplace, err := stapi.GetMarketplaceAtLocation(s.ship.Location)
		if err != nil {
			s.error(err)
			return
		}

		err = analysis.RecordMarketplaceAtLocation(s.ship.Location, marketplace)
		if err != nil {
			s.error(err)
			return
		}
		time.Sleep(time.Second * 30)
	}

}
