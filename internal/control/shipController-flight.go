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

	s.log("departing... flightplan ID: %s", flightplan.ID)

	sleepDuration := time.Minute
	totalFlightDuration := flightplan.ArrivesAt.Sub(*flightplan.CreatedAt)
	for {

		var ferryString string
		if fp.cargo == nil {
			ferryString = "(FERRY) "
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

			time.Sleep(time.Second)

			flightplan, err = s.core.user.GetFlightplan(flightplan.ID)
			if err != nil {
				return err
			}

			if flightplan.FlightTimeRemaining == 0 {
				s.log("%sarrived at %s", ferryString, flightplan.ArrivesAt.Format(time.Kitchen))
				break
			}

		}

		s.log("%sen route - %.2f%% complete, %.0fs remaining", ferryString, percentageComplete, time.Until(*flightplan.ArrivesAt).Seconds())
		time.Sleep(sleepDuration)
	}

	return nil
}
