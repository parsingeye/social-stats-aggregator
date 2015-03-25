package main

import (
	"errors"
	"github.com/codegangsta/cli"
	"github.com/guilex/social-stats-aggregator/service"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
)

func getConfig(c *cli.Context) (service.Config, error) {

	yamlPath := c.GlobalString("config")

	config := service.Config{}

	if _, err := os.Stat(yamlPath); err != nil {
		return config, errors.New("config path not valid")
	}

	ymlData, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(ymlData), &config)
	return config, err
}

func main() {

	app := cli.NewApp()
	app.Name = "SSA"
	app.Usage = "work with the social stat aggregator microservice"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			"config, c",
			"config.yaml",
			"config file to use",
			"APP_CONFIG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run the http server",
			Action: func(c *cli.Context) {
				cfg, err := getConfig(c)
				if err != nil {
					log.Fatal(err)
					return
				}

				svc := service.StatService{}

				if err = svc.Run(cfg); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "migratedb",
			Usage: "Perform database migrations",
			Action: func(c *cli.Context) {
				cfg, err := getConfig(c)
				if err != nil {
					log.Fatal(err)
					return
				}

				svc := service.StatService{}

				if err = svc.Migrate(cfg); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "update",
			Usage: "Perform stats update",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "interval",
					Usage: "set update interval",
					Value: 1,
				},
			},
			Action: func(c *cli.Context) {
				cfg, err := getConfig(c)
				if err != nil {
					log.Fatal(err)
					return
				}

				interval := c.Int("interval")

				svc := service.StatService{}

				if err = svc.Update(cfg, interval); err != nil {
					log.Fatal(err)
				}
			},
		},
	}

	app.Run(os.Args)

}
