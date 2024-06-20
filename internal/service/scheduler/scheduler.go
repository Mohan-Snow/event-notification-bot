package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

func Start(schedulerFunc func()) error {
	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(10).Second().Do(schedulerFunc)
	if err != nil {
		return err
	}

	scheduler.StartAsync()

	return nil
}
