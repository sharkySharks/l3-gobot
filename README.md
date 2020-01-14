# L3-Gobot

A Slack bot based on [Hubot](https://hubot.github.com/) written in Golang using the [go-joe/joe](https://github.com/go-joe/joe) library. Named after the piloting droid [L3-37](https://starwars.fandom.com/wiki/L3-37), aka L3.

This bot will trigger a remote Jenkins job when you interact with it in Slack. It can also save information to a redis store.

## Quickstart

1. Clone the repo to any directory and navigate into the new cloned directory.
2. Run `go get && go install` in the `l3-gobot` directory.
3. Start a redis server, running `redis-server`, in a separate terminal window.
4. Run `go run main.go` to start the bot.

You also need to have a `.env` file that has a number of configuration/secret values set. These values can be found in the ComOps Vault account under the `l3-gobot` secret. `.env.example` lists all of the values needed for the .env file.