package control

import "time"

func (s *ShipController) doFlight(fp *plannedFlight) error {
	s.log("preparing for flight")

	for _, task := range fp.preflightTasks {
		if err := task(); err != nil {
			return err
		}
	}

	flightplan, err := s.fileFlightplan(fp)
	if err != nil {
		return err
	}

	s.log("departing...\nFlightplan ID: %s", flightplan.ID)

	sleepDuration := time.Minute
	totalFlightDuration := flightplan.ArrivesAt.Sub(*flightplan.CreatedAt)
	for {

		flightplan, err = s.core.user.GetFlightplan(flightplan.ID)
		if err != nil {
			return err
		}

		if ut := time.Until(*flightplan.ArrivesAt); ut < sleepDuration {
			sleepDuration = ut + (time.Second * 2)
		}

		var percentageComplete float32
		{
			durationFlown := time.Since(*flightplan.CreatedAt)
			percentageComplete = float32(durationFlown) / float32(totalFlightDuration) * 100
		}

		if flightplan.TerminatedAt != nil {
			s.log("arrived at %s", flightplan.TerminatedAt.Format(time.Kitchen))
			break
		}

		s.log("en route - %.2f%% complete, %ds remaining", percentageComplete, flightplan.FlightTimeRemaining)
		time.Sleep(sleepDuration)
	}

	s.log("grabbing marketplace data...")
	if err = s.grabMarketplaceData(); err != nil {
		s.log("WARNING: failed to grab marketplace data after flight\n%s", err.Error()) // TODO: nice warning
	}

	return nil
}