package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	l3 "github.com/go-joe/joe"
	"github.com/go-joe/redis-memory"
	"github.com/go-joe/slack-adapter"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

type L337Bot struct {
	*l3.Bot
}

type Jenkins struct {
	Job Details `json:"task"`
}
type Details struct {
	Url string `json:"url"`
}

func main() {
	// load .env file
	e := godotenv.Load()
	if e != nil {
		log.Fatal("Error loading .env file")
	}

	b := &L337Bot{
		Bot: l3.New("L3-Bot",
			redis.Memory("localhost:6379"),
			slack.Adapter(os.Getenv("SLACK_TOKEN")),
		),
	}

	b.Respond("run jenkins job", b.ExecuteJenkinsJob)
	b.Respond("remember (.+) is (.+)", b.Remember)
	b.Respond("what is (.+)", b.WhatIs)
	b.Respond("show keys", b.ShowKeys)

	err := b.Run()
	if err != nil {
		b.Logger.Fatal(err.Error())
	}
}

func (b *L337Bot) ExecuteJenkinsJob(msg l3.Message) error {
	url := fmt.Sprintf("%s/job/git-corp-connection-check/build", os.Getenv("JENKINS_URL"))
	res, err := b.sendJenkinsRequest("POST", url, "")
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 201 {
		errorStr := fmt.Sprintf("Did not receive successful response: %v", string(body))
		return errors.New(errorStr)
	}

	jenkinsBuildInfoUrl := res.Header.Get("Location") + "api/json"

	jobUrl, err := b.getJenkinsBuildInfo(jenkinsBuildInfoUrl)

	msg.Respond("Successfully triggered jenkins job!\nJob url: %s", jobUrl)

	return nil
}

func (b *L337Bot) sendJenkinsRequest(method string, url string, body string) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Basic "+b.Base64Encode(os.Getenv("JENKINS_USER"), os.Getenv("JENKINS_USER_TOKEN")))
	req.Header.Set("Jenkins-Crumb", os.Getenv("JENKINS_CRUMB"))
	client := &http.Client{}
	res, e := client.Do(req)
	if e != nil {
		errMsg := fmt.Sprintf("Error sending Jenkins API request: ", e)
		return nil, errors.New(errMsg)
	}
	return res, nil
}

func (b *L337Bot) getJenkinsBuildInfo(url string) (string, error) {
	res, err := b.sendJenkinsRequest("GET", url, "")
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		errorStr := fmt.Sprintf("Did not receive successful response: %v, %v", res.StatusCode, string(body))
		return "", errors.New(errorStr)
	}

	s := new(Jenkins)
	unmarshallError := json.Unmarshal(body, &s)
	if unmarshallError != nil {
		log.Error("Unmarshalling error: ", unmarshallError)
		return "", unmarshallError
	}
	return s.Job.Url, nil
}

func (b *L337Bot) Base64Encode(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (b *L337Bot) Remember(msg l3.Message) error {
	key, value := msg.Matches[0], msg.Matches[1]
	msg.Respond("OK, I'll remember %s is %s", key, value)
	return b.Store.Set(key, value)
}

func (b *L337Bot) WhatIs(msg l3.Message) error {
	key := msg.Matches[0]
	var value string
	ok, err := b.Store.Get(key, &value)
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve key %q from brain", key)
	}

	if ok {
		msg.Respond("%s is %s", key, value)
	} else {
		msg.Respond("I do not remember %q", key)
	}

	return nil
}

func (b *L337Bot) ShowKeys(msg l3.Message) error {
	keys, err := b.Store.Keys()
	if err != nil {
		return err
	}

	msg.Respond("I got %d keys:", len(keys))
	for i, k := range keys {
		msg.Respond("%d) %q", i+1, k)
	}
	return nil
}
