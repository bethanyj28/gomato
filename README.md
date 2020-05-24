# gomato
An open source pomodoro library written in Go. 

## Getting Started

You can initialize a timer in one of two ways:

1. Using default settings

Default settings use the standard library log package and a cache with no expiration or cleanup.

```
pomodoro := gomato.NewDefaultTimeKeeper()
```

2. Using custom settings

You can use your own standard log implementation and cache with custom cleanup. If you leave the cache `nil`, then it will create a cache with no expiration or cleanup. If you leave the logger `nil`, it will use a no-op logger.

```
import (
        "log"
        "time"

        gcache "github.com/patrickmn/go-cache"
        "github.com/bethanyj28/gomato"
)

logger := log.New(os.Stdout, "GOMATO: ", log.Lshortfile)
cache := gcache.New(5 * time.Minute, 5 * time.Minute)

pomodoro := gomato.NewTimeKeeper(logger, cache)
```

Starting a timer is simple:
```
pomodoro := gomato.NewDefaultTimeKeeper()
uID, err := pomodoro.Start("userID", 25 * time.Minute, action1(), action2())
```
The `userID` is optional. If you opt to not use your own, then one will be generated and returned. With `StartWithTime` you can optionally provide a starting time if it is not whatever `time.Now()` returns. If you enter a zero duration, it will be set to the typical 25 minutes that pomodoros are set. Actions are variadic, so you can enter as many or as few as you'd like, just make sure that there's nothing passed into them and they do not return anything. Feel free to look at the tests for an example on how to still pass in variables to the functions.

Pausing, Resuming, and Stopping simply require a userID.

## Testing

For simple testing, this package abides by the following interface: 

```
type PomodoroManager interface {
	StartWithTime(uID string, start time.Time, duration time.Duration, actions ...func()) (string, error)
	Start(uID string, duration time.Duration, actions ...func()) (string, error)
	Pause(uID string) error
	Resume(uID string) error
	Stop(uID string) error
}
```

## Planned use cases
- [ ] [Slack App](https://github.com/bethanyj28/gomato-slack)
- [ ] CLI
