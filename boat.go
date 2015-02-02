package swail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Boat struct {
	UserToken string
	BoatId    string
}

// Returned by GetBoatInfos webservice
type BoatInfos struct {
	BoatType string    `json:"Bti"`
	Colour   string    `json:"Colour"`
	Dir      int       `json:"Dir"`
	Coord    []float64 `json:"Coord"`
	Name     string    `json:"Name"`
	Speed    float32   `json:"Speed"`
	TWA      int       `json:"TWA"`
	TWS      float32   `json:"TWS"`
}

type BoatPolar struct {
	Cog   int     `json:"Dir"`
	Speed float32 `json:"Speed"`
}

func (boat *Boat) swailGet(webservice string, result interface{}) (err error) {
	url := "http://swail.io/api/v1" + webservice
	url = strings.Replace(url, "boatId", boat.BoatId, 1)
	url = strings.Replace(url, "userToken", boat.UserToken, 1)

	log.Printf("Swail Request: %s", url)

	req, err := http.Get(url)
	if err != nil {
		return
	}
	if req.StatusCode != 200 {
		err = Error(fmt.Sprintf("Swail returned invalid status code %i", req.StatusCode))
		return
	}

	bodyBuffer := &bytes.Buffer{}
	_, err = bodyBuffer.ReadFrom(req.Body)
	if err != nil {
		return
	}
	req.Body.Close()

	body := bodyBuffer.Bytes()
	if err = json.Unmarshal(body, result); err != nil {
		return
	}
	return
}

func (bi *BoatInfos) Position() Position {
	if len(bi.Coord) != 2 {
		log.Panic("Invalid GPS coordinates on this boat!!", bi.Coord)
	}
	return NewPosition(bi.Coord[1], bi.Coord[0])
}

func (boat *Boat) GetBoatInfos() (infos BoatInfos, err error) {
	// http://swail.io/api/v1/boats/54c35ac521021e05f886ee01/user/378f0e0c-8fa6-4914-9dfe-98c2f1a7b56a

	infos.Coord = make([]float64, 2)

	err = boat.swailGet("/boats/boatId/user/userToken", &infos)
	if err != nil {
		return
	}
	if len(infos.Coord) != 2 {
		err = Error(fmt.Sprintf("Invalid coordinates provided for this boat (%v)", infos.Coord))
	}
	log.Printf("Successful http request: %v", infos)
	return
}

func (boat *Boat) GetBoatPolars() (polars []BoatPolar, err error) {
	err = boat.swailGet("/boats/boatId/polars", &polars)
	if err != nil {
		return
	}
	log.Printf("Polars: %v", polars)
	return
}

func (boat *Boat) SetCourse(course int) error {
	var infos BoatInfos
	return boat.swailGet(fmt.Sprintf("/boats/boatId/heading/%d/user/userToken", course), &infos)
}

func (boat *Boat) Autopilot(strategy PilotStrategy) {
	log.Printf("Starting autopilot with strategy %v", strategy)
	log.Printf("Press Ctrl-C to stop")
	for {
		bearing, err := strategy.Bearing(boat)

		if err != nil {
			log.Printf("Pilot error: %s", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if err := boat.SetCourse(bearing); err != nil {
			log.Printf("Error setting course: %s", err)
		}

		time.Sleep(10 * time.Second)
	}
}
