package routes

import (
	"social-network/controllers"
	"social-network/pkg/db/sqlite"
	"sync"
	"time"
)

type RateLimiter struct {
	LastTime    time.Time
	Counter     int
	TimeToUnban time.Time
}

type Handler struct {
	Repo    *sqlite.Repo
	Cntrlrs *controllers.Controller
	Mu      sync.RWMutex

	// ratelimiter for http
	LastRL map[string]*RateLimiter
}

func NewHandler(repo *sqlite.Repo, cntrlrs *controllers.Controller) *Handler {
	return &Handler{
		Repo:    repo,
		Cntrlrs: cntrlrs,
		LastRL:  make(map[string]*RateLimiter),
	}
}
