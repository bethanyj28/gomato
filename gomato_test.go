package gomato

import (
	"log"
	"os"
	"testing"
	"time"

	gcache "github.com/patrickmn/go-cache"
)

func TestStart(t *testing.T) {
	var tests = []struct {
		name      string
		expectErr bool
		userID    string
		startTime time.Time
	}{
		{
			name:      "start success",
			expectErr: false,
			userID:    "testUser",
			startTime: time.Now(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, tk := setUpTest()
			var finished bool
			_, err := tk.Start(tt.userID, tt.startTime, 1, timerFinished(&finished))
			if !tt.expectErr && err != nil {
				t.Fatalf("an unexpected error occurred: %s", err.Error())
			}

			time.Sleep(65 * time.Second)

			if !finished {
				t.Fatal("final actions not run")
			}

			if _, ok := c.Get(tt.userID); ok {
				t.Fatal("cache record not deleted")
			}
		})
	}
}

func setUpTest() (*gcache.Cache, *TimeKeeper) {
	c := gcache.New(-1, -1)
	return c, NewTimeKeeper(log.New(os.Stdout, "GOMATO: ", log.Lshortfile), c)
}

func timerFinished(finished *bool) func() {
	return func() {
		*finished = true
	}
}
