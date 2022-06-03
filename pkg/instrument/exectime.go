package instrument

import (
	"time"

	"github.com/rs/zerolog/log"
)

func Track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func Duration(msg string, start time.Time) {
	log.Debug().Msgf("%v: %v\n", msg, time.Since(start))
}
