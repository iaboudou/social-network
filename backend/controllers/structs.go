package controllers

import (
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
	"time"
)

type Controller struct {
	DB *sqlite.Repo
}

func NewController(db *sqlite.Repo) *Controller {
	return &Controller{DB: db}
}

// rate limiter struct and methods for websocket messages
type RateLimiter struct {
	Last        time.Time
	Count       int
	Blocked     bool
	Deleteblock time.Time
}

func (rl *RateLimiter) Check() bool {
	now := time.Now()
	if rl.Blocked {
		if now.After(rl.Deleteblock) {
			rl.Blocked = false
			rl.Count = 0
		} else {
			return false
		}
	}

	if helpers.MessageRLExceeded(rl.Count, rl.Last) {
		rl.Blocked = true
		rl.Deleteblock = now.Add(10 * time.Second)
		return false
	}

	rl.Count++
	rl.Last = now
	return true
}
