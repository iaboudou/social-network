package config

import "time"

const Port = ":3000"

// message rate limiting
const (
	Max_Messages = 100
	FiveSec      = 5 * time.Second
)