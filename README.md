# gomato
An open source slack pomodoro timer written in Go. This project is currently in development and is about halfway to a usable product. Please see [upcoming tasks](#upcoming-tasks)

## Getting Started

On slack, create a bot user. Under **Settings** -> **Basic Information**, scroll to **App Credentials** and copy the **Verification Token** (TODO: update authentication)

Clone the repo, create a file called `environment.env`. To this file add the following:

```
SLACK_VERIFICATION_TOKEN=<Your Slack verification token>
```

Run the server by running 

```
docker build -t bethanyj28/gomato .
docker run --rm -p 8080:8080 bethanyj28/gomato
```

Install ngrok and run 

```
ngrok http 8080
```

Copy the forwarding URL (the http one). On the Slack API homepage for your app, under **Features** -> **Slash Commands**, add the following commands mapped to the following endpoints (command names are suggestions):

```
/gomato_start -> <ngrok url>/timer/start (optional set duration)
/gomato_pause -> <ngrok url>/timer/pause
/gomato_resume -> <ngrok url>/timer/resume
/gomato_stop -> <ngrok url>/timer/stop
```

From here, you should be able to run those commands from your workspace! The commands are straightforward except for start, which has the option to set a duration. The default duration is 20 minutes

```
/gomato_start
/gomato_start 40
```

## Upcoming Tasks
- [ ] Respond to user via Slack when timer is up
- [ ] Option to set do not disturb during timer
- [ ] Set a timer for a break
