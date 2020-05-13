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
		{
			name:      "start 0 time success",
			expectErr: false,
			userID:    "testUser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, tk := setUpTest()
			var finished bool
			_, err := tk.Start(tt.userID, tt.startTime, 2*time.Second, timerFinished(&finished))
			if !tt.expectErr && err != nil {
				t.Fatalf("an unexpected error occurred: %s", err.Error())
			}

			time.Sleep(4 * time.Second)

			if !finished {
				t.Fatal("final actions not run")
			}

			if _, ok := c.Get(tt.userID); ok {
				t.Fatal("cache record not deleted")
			}
		})
	}
}

func TestPause(t *testing.T) {
	var tests = []struct {
		name      string
		expectErr bool
		userID    string
		setUp     func(id string, finished *bool, tk *TimeKeeper)
	}{
		{
			name:      "pause success",
			expectErr: false,
			userID:    "testUser",
			setUp: func(id string, finished *bool, tk *TimeKeeper) {
				if _, err := tk.Start(id, time.Now(), 2*time.Second, timerFinished(finished)); err != nil {
					t.Fatalf("an unexpected error occurred: %s", err.Error())
				}
			},
		},
		{
			name:      "pause fail - no ID",
			expectErr: true,
			userID:    "",
			setUp:     func(id string, finished *bool, tk *TimeKeeper) {},
		},
		{
			name:      "pause fail - ID DNE",
			expectErr: true,
			userID:    "ID_DNE",
			setUp: func(id string, finished *bool, tk *TimeKeeper) {
				if _, err := tk.Start("testID", time.Now(), 2*time.Second, timerFinished(finished)); err != nil {
					t.Fatalf("an unexpected error occurred: %s", err.Error())
				}

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, tk := setUpTest()
			var finished bool
			tt.setUp(tt.userID, &finished, tk)

			err := tk.Pause(tt.userID)
			if !tt.expectErr && err != nil {
				t.Fatalf("an unexpected error occurred: %s", err.Error())
			}

			if tt.expectErr && err == nil {
				t.Fatalf("expected error")
			}

			if tt.expectErr {
				return
			}

			time.Sleep(4 * time.Second)

			if finished {
				t.Fatalf("the timer finished when it should be paused")
			}

			if _, ok := c.Get(tt.userID); !ok {
				t.Fatal("cache record deleted")
			}

		})
	}
}

func TestResume(t *testing.T) {
	var tests = []struct {
		name      string
		expectErr bool
		userID    string
		setUp     func(id string, finished *bool, tk *TimeKeeper)
	}{
		{
			name:      "resume success",
			expectErr: false,
			userID:    "testUser",
			setUp: func(id string, finished *bool, tk *TimeKeeper) {
				if _, err := tk.Start(id, time.Now(), 2*time.Second, timerFinished(finished)); err != nil {
					t.Fatalf("an unexpected error occurred: %s", err.Error())
				}

				if err := tk.Pause(id); err != nil {
					t.Fatalf("an unexpected error occurred: %s", err.Error())
				}
			},
		},
		{
			name:      "resume fail - no ID",
			expectErr: true,
			userID:    "",
			setUp:     func(id string, finished *bool, tk *TimeKeeper) {},
		},
		{
			name:      "resume fail - ID DNE",
			expectErr: true,
			userID:    "ID_DNE",
			setUp: func(id string, finished *bool, tk *TimeKeeper) {
				if _, err := tk.Start("testID", time.Now(), 2*time.Second, timerFinished(finished)); err != nil {
					t.Fatalf("an unexpected error occurred: %s", err.Error())
				}

				if err := tk.Pause("testID"); err != nil {
					t.Fatalf("an unexpected error occurred: %s", err.Error())
				}

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, tk := setUpTest()
			var finished bool
			tt.setUp(tt.userID, &finished, tk)

			err := tk.Resume(tt.userID)
			if !tt.expectErr && err != nil {
				t.Fatalf("an unexpected error occurred: %s", err.Error())
			}

			if tt.expectErr && err == nil {
				t.Fatalf("expected error")
			}

			if tt.expectErr {
				return
			}

			time.Sleep(4 * time.Second)

			if !finished {
				t.Fatalf("the timer did not resume")
			}

			if _, ok := c.Get(tt.userID); ok {
				t.Fatal("cache record not deleted")
			}

		})
	}
}

func TestStop(t *testing.T) {
	var tests = []struct {
		name      string
		expectErr bool
		userID    string
		setUp     func(id string, finished *bool, tk *TimeKeeper)
	}{
		{
			name:      "stop success",
			expectErr: false,
			userID:    "testUser",
			setUp: func(id string, finished *bool, tk *TimeKeeper) {
				if _, err := tk.Start(id, time.Now(), 2*time.Second, timerFinished(finished)); err != nil {
					t.Fatalf("an unexpected error occurred: %s", err.Error())
				}
			},
		},
		{
			name:      "stop fail - no ID",
			expectErr: true,
			userID:    "",
			setUp:     func(id string, finished *bool, tk *TimeKeeper) {},
		},
		{
			name:      "stop fail - ID DNE",
			expectErr: true,
			userID:    "ID_DNE",
			setUp: func(id string, finished *bool, tk *TimeKeeper) {
				if _, err := tk.Start("testID", time.Now(), 2*time.Second, timerFinished(finished)); err != nil {
					t.Fatalf("an unexpected error occurred: %s", err.Error())
				}

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, tk := setUpTest()
			var finished bool
			tt.setUp(tt.userID, &finished, tk)

			err := tk.Stop(tt.userID)
			if !tt.expectErr && err != nil {
				t.Fatalf("an unexpected error occurred: %s", err.Error())
			}

			if tt.expectErr && err == nil {
				t.Fatalf("expected error")
			}

			if tt.expectErr {
				return
			}

			time.Sleep(4 * time.Second)

			if finished {
				t.Fatalf("the timer finished when it should be paused")
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
