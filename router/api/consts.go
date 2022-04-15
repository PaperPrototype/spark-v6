package api

import "time"

// heroku times out at 30 seconds https://devcenter.heroku.com/articles/request-timeout
const MaxTimeoutSeconds float64 = 25
const SleepTime time.Duration = 3 // mulitply by seconds
