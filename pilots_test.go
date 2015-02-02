package swail

import "testing"

type windTest struct {
	targetTwa      int
	currentTwa     int
	currentBearing int
	newBearing     int
}

var windTests = []windTest{
	{-90, 0, 0, 270},
	{90, 0, 0, 90},
	{90, 0, 180, 270},
	{-39, -11, 234, 206},
}

type TestBoat struct {
	infos BoatInfos
}

func (b TestBoat) GetBoatInfos() (BoatInfos, error) {
	return b.infos, nil
}

func TestWindVane(t *testing.T) {
	for _, tc := range windTests {

		bi := BoatInfos{TWA: tc.currentTwa, Dir: tc.currentBearing}
		b := TestBoat{infos: bi}

		pilot := WindVane{TargetTWA: tc.targetTwa}
		newBearing, err := pilot.Bearing(b)

		if err != nil {
			t.Logf("Pilot returned error: %v", err)
			t.Fail()
		}

		if newBearing != tc.newBearing {
			t.Logf("Windvane returned wrong new bearing. Env: %v Returned: %v", tc, newBearing)
			t.Fail()
		}
	}

}
