package swail

import "log"

type PilotStrategy interface {
	Bearing(b Pilotable) (int, error)
}

type Pilotable interface {
	GetBoatInfos() (BoatInfos, error)
}

type WindVane struct {
	TargetTWA int
}

func (w WindVane) Bearing(b Pilotable) (newBearing int, err error) {
	bi, err := b.GetBoatInfos()
	if err != nil {
		return
	}

	// TWA = Bearing - WindAngle
	// Bearing = WindAngle + TWA
	// WindAngle = Bearing - TWA
	wind := bi.Dir - bi.TWA
	// -11 + 234 = 223 --- SHOULD BE 245

	// 223 - (-39) = 262

	newBearing = (wind + w.TargetTWA) % 360
	if newBearing < 0 {
		newBearing += 360
	}

	return newBearing, nil
}

/* Stupid pilot who goes straight to target */
type MotorPilot struct {
	Target Position
}

func (m MotorPilot) Bearing(boat Pilotable) (newCourse int, err error) {
	boatInfos, err := boat.GetBoatInfos()
	if err != nil {
		return 0, err
	}

	currentPosition := boatInfos.Position()

	newCourse = int(currentPosition.BearingTo(m.Target))
	log.Printf("DTW: %.1fnm BRG: %.1fº Course: %v", currentPosition.DistanceTo(m.Target), currentPosition.BearingTo(m.Target), newCourse)
	return
}
