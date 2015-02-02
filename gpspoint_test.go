package swail

import (
	"math"
	"testing"
)

type geoTest struct {
	a        Position
	b        Position
	distance float64
	course   float64
}

var geoTests = []geoTest{
	// Same point
	{Position{43.1, 12.2}, Position{43.1, 12.2}, 0, 0},
	// Going south for 1 nm
	{Position{43.0 + 1.0/60, 12.2}, Position{43.0, 12.2}, 1, 180},
	// Going east for 1 nm (on the equator)
	{Position{0, 12.0 + 2.0/60}, Position{0, 12.0 + 3.0/60}, 1, 90},
	// Going west for 60 nm (on the equator)
	{Position{0, 12}, Position{0, 11}, 60, 270},
	// Going east over the date separation line
	{Position{0, -180.0 - 1.0/60}, Position{0, 180.0 + 1.0/60}, 2, 90},
	// Going west over the date separation line
	{Position{0, 180.0 + 1.0/60}, Position{0, -180.0 - 1.0/60}, 2, 270},
	// Sailing SW in Transquadra
	{Position{29.6, -41.23}, Position{14.5, -60.8}, 1412.8, 234},
}

func TestDistances(t *testing.T) {
	for _, tc := range geoTests {
		d := tc.a.DistanceTo(tc.b)
		delta := math.Abs(d - tc.distance)
		if delta > 0.1 {
			t.Logf("Incorrect distance between %v and %v - Expected %f, Got %f (∂=%f)", tc.a, tc.b, tc.distance, d, delta)
			t.Fail()
		}
	}
}

func TestCourses(t *testing.T) {
	for _, tc := range geoTests {
		c := tc.a.BearingTo(tc.b)
		delta := math.Abs(c - tc.course)
		if delta > 0.5 {
			t.Logf("Incorrect course between %v and %v - Expected %f, Got %f (∆=%f)", tc.a, tc.b, tc.course, c, delta)
			t.Fail()
		}
	}
}
