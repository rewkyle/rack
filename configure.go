package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/gopkg.in/ini.v1"
	"github.com/rackspace/rack/util"
)

func configure(c *cli.Context) {
	intro := []string{"\nThis interactive session will walk you through creating",
		"a profile in your configuration file. You may fill in all or none of the",
		"values.\n"}
	fmt.Println(strings.Join(intro, "\n"))
	reader := bufio.NewReader(os.Stdin)
	m := map[string]string{
		"username": "",
		"api-key":  "",
		"region":   "",
	}
	fmt.Print("Rackspace Username: ")
	username, _ := reader.ReadString('\n')
	m["username"] = strings.TrimSpace(username)

	fmt.Print("Rackspace API key: ")
	apiKey, _ := reader.ReadString('\n')
	m["api-key"] = strings.TrimSpace(apiKey)

	fmt.Print("Rackspace Region: ")
	region, _ := reader.ReadString('\n')
	m["region"] = strings.ToUpper(strings.TrimSpace(region))

	fmt.Print("Profile Name (leave blank to create a default profile): ")
	profile, _ := reader.ReadString('\n')
	profile = strings.TrimSpace(profile)

	configFileLoc, err := util.ConfigFileLocation()
	var cfg *ini.File
	cfg, err = ini.Load(configFileLoc)
	if err != nil {
		// fmt.Printf("Error loading config file: %s\n", err)
		cfg = ini.Empty()
	}

	if strings.TrimSpace(profile) == "" || strings.ToLower(profile) == "default" {
		profile = "DEFAULT"
	}

	for {
		if section, err := cfg.GetSection(profile); err == nil && len(section.Keys()) != 0 {
			fmt.Printf("\nA profile named %s already exists. Overwrite? (y/n): ", profile)
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)
			switch strings.ToLower(choice) {
			case "y", "yes":
				break
			case "n", "no":
				fmt.Print("Profile Name: ")
				profile, _ = reader.ReadString('\n')
				profile = strings.TrimSpace(profile)
				continue
			default:
				continue
			}
			break
		}
		break
	}

	section, err := cfg.NewSection(profile)
	if err != nil {
		//fmt.Printf("Error creating new section [%s] in config file: %s\n", profile, err)
		return
	}

	for key, val := range m {
		section.NewKey(key, val)
	}

	err = cfg.SaveTo(configFileLoc)
	if err != nil {
		//fmt.Printf("Error saving config file: %s\n", err)
		return
	}

	if profile == "DEFAULT" {
		fmt.Printf("\nCreated new default profile for username %s", username)
	} else {
		fmt.Printf("\nCreated profile %s with username %s", profile, username)
	}

}
