package main

// blitz is a command line tool for searching logs.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/user"

	"github.com/blitzlog/app/client"
	"github.com/blitzlog/errors"
	"github.com/blitzlog/log"
	"gopkg.in/yaml.v2"
)

const (
	credentialsSubPath = ".blitz/credentials"
	apiAddress         = "https://test.blitzlog.com:8080"
)

func main() {

	// parse flags
	parseFlags()

	err := logs(fFilter, fStart, fEnd)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func logs(filter, start, end string) error {

	// get token
	accountId, token, err := getCredentials()
	if err != nil {
		return errors.Wrap(err, "getting token")
	}

	// create new api client
	apiClient := client.New(apiAddress)

	// use client to get logs
	resp, err := apiClient.GetLogs(accountId, token)
	if err != nil {
		return errors.Wrap(err, "getting response from api server")
	}

	for _, lg := range resp.GetLogs() {
		fmt.Println(log.Format(lg))
	}

	return nil
}

type credentials struct {
	AccountId string `yaml:"account_id"`
	Token     string `yaml:"token"`
}

// getCredentials from credentials file.
func getCredentials() (string, string, error) {

	// get user
	usr, err := user.Current()
	if err != nil {
		return "", "", errors.Wrap(err, "error getting current user")
	}

	// read credentials file
	credentialsPath := fmt.Sprintf("%s/%s", usr.HomeDir, credentialsSubPath)
	credentialsBytes, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to read credentials file")
	}

	// unmarshal credentials
	var creds credentials
	err = yaml.Unmarshal(credentialsBytes, &creds)

	return creds.AccountId, creds.Token, nil
}

// flags passed to blitz
var (
	fFilter string
	fStart  string
	fEnd    string
)

func parseFlags() {

	flag.StringVar(&fFilter, "filter", "", "apply filter when searching logs")
	flag.StringVar(&fStart, "start", "", "start of duration for searching logs")
	flag.StringVar(&fEnd, "end", "", "end of duration for searching logs")

	flag.Parse()
}
