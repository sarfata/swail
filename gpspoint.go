package swail

import "math"

type GeoAngle float64

func (a GeoAngle) toRadians() float64 {
	return float64(a * math.Pi / 180)
}

type Position struct {
	Latitude  GeoAngle
	Longitude GeoAngle
}

func NewPosition(lat, lon float64) Position {
	return Position{GeoAngle(lat), GeoAngle(lon)}
}

/* Returns the distance in nautical miles to the destination */
func (p *Position) DistanceTo(p2 Position) float64 {
	// Using the Haversine formula to calculate distance along the great-circles
	// Thanks to: http://www.movable-type.co.uk/scripts/latlong.html

	var R float64 = 6371 // km
	var φ1 = p.Latitude.toRadians()
	var φ2 = p2.Latitude.toRadians()

	var Δφ = ((GeoAngle)(p2.Latitude - p.Latitude)).toRadians()
	var Δλ = ((GeoAngle)(p2.Longitude - p.Longitude)).toRadians()

	var a = math.Sin(Δφ/2)*math.Sin(Δφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	var d = R * c

	return d / 1.852
}

/* Returns the course in true degrees to the destination */
func (p *Position) BearingTo(p2 Position) float64 {
	// Again, thanks to: http://www.movable-type.co.uk/scripts/latlong.html

	var φ1 = float64(p.Latitude.toRadians())
	var φ2 = float64(p2.Latitude.toRadians())

	var Δλ = float64(((GeoAngle)(p2.Longitude - p.Longitude)).toRadians())

	var y = math.Sin(Δλ) * math.Cos(φ2)
	var x = math.Cos(φ1)*math.Sin(φ2) - math.Sin(φ1)*math.Cos(φ2)*math.Cos(Δλ)
	var brng = math.Atan2(y, x) * 180 / math.Pi

	// Convert from -180 / 180 to 0 / 360
	// Cheap-ass modulo for float
	brng += 360
	if brng >= 360 {
		brng -= 360
	}
	return brng
}
