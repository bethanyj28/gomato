package gomato

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	gcache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

// TimeKeeper represents the necessary components to manage pomodoros
type TimeKeeper struct {
	cache  *gcache.Cache
	logger *log.Logger
}

// NewDefaultTimeKeeper instantiates a TimeKeeper with default logging and cache options
func NewDefaultTimeKeeper() *TimeKeeper {
	return NewTimeKeeper(log.New(os.Stdout, "GOMATO: ", log.Lshortfile), gcache.New(-1, -1))
}

// NewTimeKeeper instantiates a TimeKeeper object
func NewTimeKeeper(l *log.Logger, c *gcache.Cache) *TimeKeeper {
	if c == nil {
		c = gcache.New(-1, -1) // cache with no expiration, no cleanup
	}

	if l == nil { // assume no logging is wanted
		l = log.New(ioutil.Discard, "", 0)
	}
	return &TimeKeeper{cache: c, logger: l}
}

type pomodoro struct {
	startTime       time.Time
	currentDuration time.Duration
	timer           *time.Timer
}

// Start begins a new pomodoro
// A user identifier should be passed through, but if it is not then it will be generated and returned
func (t *TimeKeeper) Start(uID string, start time.Time, duration int32, actions ...func()) (string, error) {
	if uID == "" {
		t.logger.Print("[INFO] User ID not provided, setting ID")
		uID = xid.New().String()
	}

	if start.IsZero() {
		t.logger.Print("[INFO] Time is zero, setting to current time")
		start = time.Now()
	}

	if duration == 0 {
		t.logger.Print("[INFO] Duration not set, setting to 20 minutes")
		duration = 20
	}

	fDur, err := time.ParseDuration(fmt.Sprintf("%dm", duration))
	if err != nil {
		t.logger.Printf("[ERROR] %s", err.Error())
		return "", errors.Wrap(err, "failed to parse duration")
	}

	p := pomodoro{
		startTime:       start,
		currentDuration: fDur,
		timer:           time.AfterFunc(fDur, t.runActions(uID, actions...)),
	}

	t.cache.SetDefault(uID, &p)

	return uID, nil
}

func (t *TimeKeeper) runActions(userID string, actions ...func()) func() {
	return func() {
		t.logger.Print("[INFO] Running finish timer actions")
		for _, action := range actions {
			action()
		}

		t.cache.Delete(userID)
	}
}
