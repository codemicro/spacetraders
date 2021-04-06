package control

import (
	"time"
)

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

	// TODO: we probably don't need to be retrieving the flightplan from the API every single time
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

		if flightplan.ArrivesAt.Before(time.Now()) {
			time.Sleep(time.Until(*flightplan.ArrivesAt) + time.Second)
			s.log("arrived at %s", flightplan.ArrivesAt.Format(time.Kitchen))
			break
		}

		s.log("en route - %.2f%% complete, %ds remaining", percentageComplete, flightplan.FlightTimeRemaining)
		time.Sleep(sleepDuration)
	}

	return nil
}
