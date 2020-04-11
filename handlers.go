package gomato

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
)

type timerDetail struct {
	timeStart    time.Time
	timeDuration time.Duration
	timer        *time.Timer
}

func (s *Server) startTimer(c *gin.Context) {
	sc, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing slash command": err.Error()})
		return
	}

	if !sc.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ud := sc.Text
	if ud == "" {
		ud = "1"
	}

	ud = fmt.Sprintf("%sm", ud)

	d, err := time.ParseDuration(ud)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	t := timerDetail{
		timeStart:    time.Now(),
		timeDuration: d,
		timer:        time.AfterFunc(d, s.setTimer(sc.UserID)),
	}

	s.Cache.SetDefault(sc.UserID, &t)

	c.Status(http.StatusOK)
}

func (s *Server) pauseTimer(c *gin.Context) {
	sc, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing slash command": err.Error()})
		return
	}

	if !sc.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	td, ok := s.Cache.Get(sc.UserID)
	if !ok {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	tData, ok := td.(*timerDetail)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_ = tData.timer.Stop()
	tData.timeDuration = tData.timeDuration - time.Since(tData.timeStart)

	c.Status(http.StatusOK)
}

func (s *Server) resumeTimer(c *gin.Context) {
	sc, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing slash command": err.Error()})
		return
	}

	if !sc.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	td, ok := s.Cache.Get(sc.UserID)
	if !ok {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	tData, ok := td.(*timerDetail)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_ = tData.timer.Reset(tData.timeDuration)

	c.Status(http.StatusOK)
}

func (s *Server) stopTimer(c *gin.Context) {
	sc, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing slash command": err.Error()})
		return
	}

	if !sc.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	td, ok := s.Cache.Get(sc.UserID)
	if !ok {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	tData, ok := td.(*timerDetail)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_ = tData.timer.Stop()

	s.Cache.Delete(sc.UserID)

	c.Status(http.StatusOK)
}

func (s *Server) notifyUser(userID string) {
	fmt.Printf("Time's up!")
	s.Cache.Delete(userID)
}

func (s *Server) setTimer(userID string) func() {
	return func() {
		s.notifyUser(userID)
	}
}
