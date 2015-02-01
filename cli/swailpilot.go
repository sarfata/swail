package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/sarfata/swailpilot"
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

func (bi *BoatInfos) Position() swail.Position {
	if len(bi.Coord) != 2 {
		log.Panic("Invalid GPS coordinates on this boat!!", bi.Coord)
	}
	return swail.NewPosition(bi.Coord[1], bi.Coord[0])
}

func (boat *Boat) GetBoatInfos() (infos *BoatInfos, err error) {
	// http://swail.io/api/v1/boats/54c35ac521021e05f886ee01/user/378f0e0c-8fa6-4914-9dfe-98c2f1a7b56a

	infos = &BoatInfos{}
	infos.Coord = make([]float64, 2)

	err = boat.swailGet("/boats/boatId/user/userToken", infos)
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

func (boat *Boat) Goto(target swail.Position) {
	log.Printf("Starting autopilot with target %v", target)
	log.Printf("Press Ctrl-C to stop")
	for {
		boatInfos, err := boat.GetBoatInfos()
		if err != nil {
			log.Printf("Error fetching boat infos: %s", err)
			time.Sleep(10 * time.Second)
			continue
		}
		currentPosition := boatInfos.Position()

		newCourse := boatInfos.Dir
		if currentPosition.DistanceTo(target) < 5 {
			newCourse = int(currentPosition.BearingTo(target))
		} else {
			// TODO: Find the course with the best VMG
			newCourse = int(currentPosition.BearingTo(target))
		}
		log.Printf("DTW: %.1fnm BRG: %.1fº Course: %v", currentPosition.DistanceTo(target), currentPosition.BearingTo(target), newCourse)
		if newCourse != boatInfos.Dir {
			if err := boat.SetCourse(newCourse); err != nil {
				log.Printf("Error setting course: %s", err)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "SwailPilot"
	app.Usage = "A helping head for the swail.io sailing simulation"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "userToken, u", Usage: "Your user unique token (get it after you identify on the web)"},
		cli.StringFlag{Name: "boatId, b", Usage: "The identifier of the boat you want to control"},
	}

	var boat Boat

	app.Before = func(c *cli.Context) error {
		boat = Boat{c.GlobalString("userToken"), c.GlobalString("boatId")}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "infos",
			Usage: "Returns infos on the current status of your boat",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) {
				_, err := boat.GetBoatInfos()
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "polars",
			Usage: "Returns the current polars for the boat",
			Action: func(c *cli.Context) {
				_, err := boat.GetBoatPolars()
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "goto",
			Usage: "Go to given waypoint with max VMG",
			Action: func(c *cli.Context) {
				lat, err := strconv.ParseFloat(c.Args()[0], 64)
				if err != nil {
					log.Fatal("Invalid latitude %s", c.Args()[0])
				}
				lon, err := strconv.ParseFloat(c.Args()[1], 64)
				if err != nil {
					log.Fatal("Invalid longitude %s", c.Args()[1])
				}

				target := swail.NewPosition(lat, lon)
				boat.Goto(target)
			},
		},
	}

	app.RunAndExitOnError()
}
