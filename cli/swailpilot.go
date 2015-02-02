package main

import (
	"log"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/sarfata/swail"
)

func main() {
	app := cli.NewApp()
	app.Name = "SwailPilot"
	app.Usage = "A helping head for the swail.io sailing simulation"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "userToken, u", Usage: "Your user unique token (get it after you identify on the web)"},
		cli.StringFlag{Name: "boatId, b", Usage: "The identifier of the boat you want to control"},
	}

	var boat swail.Boat

	app.Before = func(c *cli.Context) error {
		boat = swail.Boat{c.GlobalString("userToken"), c.GlobalString("boatId")}
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
			Name:  "wind",
			Usage: "In Wind mode, the pilot will act as a wind-vane and keep a constant true wind angle",
			Action: func(c *cli.Context) {
				twa, err := strconv.ParseInt(c.Args()[0], 10, 0)
				if err != nil {
					log.Fatal("Please provide a valid target True Wind Angle")
				}
				pilot := swail.WindVane{int(twa)}
				boat.Autopilot(pilot)
			},
		},
		{
			Name:  "motor",
			Usage: "Go to given waypoint in direct line",
			Action: func(c *cli.Context) {
				lat, err := strconv.ParseFloat(c.Args()[0], 64)
				if err != nil {
					log.Fatal("Invalid latitude %s", c.Args()[0])
				}
				lon, err := strconv.ParseFloat(c.Args()[1], 64)
				if err != nil {
					log.Fatal("Invalid longitude %s", c.Args()[1])
				}

				pilot := swail.MotorPilot{swail.NewPosition(lat, lon)}
				boat.Autopilot(pilot)
			},
		},
	}

	app.RunAndExitOnError()
}
