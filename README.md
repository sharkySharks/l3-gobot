# L3-Gobot

A Slack bot inspired by [Hubot](https://hubot.github.com/) written in Golang using the [go-joe/joe](https://github.com/go-joe/joe) library. 

Named after the piloting droid [L3-37](https://starwars.fandom.com/wiki/L3-37), aka L3.

This bot will trigger a remote Jenkins job when you interact with it in Slack. It can also save information to a redis store.

## Quickstart

1. Clone the repo to any directory and navigate into the new cloned directory.
2. Run `go get && go install` in the `l3-gobot` directory.
3. Start a redis server, running `redis-server`, in a separate terminal window. (Install [redis-server](https://redis.io/topics/quickstart) if you do not already have it on your system)
4. Copy `.env.example` to `.env` and fill in secret information. 
5. Run `go run main.go` to start the bot.
6. Test out the bot by running a command in Slack: `@bot-name ping` and if it is configured correctly the bot should respond in Slack with `PONG`

### Secrets needed for .env

| Key                 | Description |
|--                   |--           |
| `SLACK_TOKEN`       | Follow [Slack docs](https://slack.com/help/articles/202035138-Add-an-app-to-your-workspace) to add an app to your Slack workspace. After creation the app will be given a token. |
| `JENKINS_CRUMB`     | Follow [this article](https://beginnersforum.net/blog/2019/11/28/jenkin-paramerized-job-api-json/) on how to retrieve a Jenkins crumb in order to authenticate a remote Jenkins API call to your Jenkins server |
| `JENKINS_USER`      | Service account user for Jenkins server |
| `JENKINS_USER_TOKEN`| Log in to Jenkins server UI with the `JENKINS_USER` and [configure an API token](https://stackoverflow.com/a/45466184/4541964) |
| `JENKINS_URL`       | Full URL (including `https://`) for Jenkins server |

## Future Features

- Containerize bot and redis setup
- Deploy bot somewhere
- Add authNZ to limit who can run certain commands
- Add tests
- Add Jenkins build info output and/or command to get information on the triggered job
