package api

import "time"

// heroku times out at 30 seconds https://devcenter.heroku.com/articles/request-timeout
const MaxTimeoutSeconds float64 = 20
const SleepTime time.Duration = 3
