package main

import (
	"errors"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"logger/service"
	"os"
)

func getConfig(c *cli.Context) (service.Config, error) {
	yamlPath := c.String("config")
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
	app := cli.App{}
	app.Name = "logger"
	app.Usage = "work with the `logger` microservice"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "config, c", Value: "config.yaml", Usage: "config file to use", EnvVars: []string{"APP_CONFIG"}},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "server",
			Usage: "Run the http server",
			Action: func(c *cli.Context) error {
				cfg, err := getConfig(c)
				if err != nil {
					log.Fatal(err)
					return err
				}

				svc := service.LoggerService{}

				if err = svc.Run(cfg); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
	app.Run(os.Args)

}
