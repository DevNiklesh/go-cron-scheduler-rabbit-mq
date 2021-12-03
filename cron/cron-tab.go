package crontab

import (
	"log"

	"github.com/robfig/cron/v3"
)

type CronTabService interface {
	AddJob(string, func()) error
	Entries() []cron.EntryID
	StartCron()
}

type crontab struct {
	cron     *cron.Cron
	entryIds []cron.EntryID
}

func New() CronTabService {
	log.Println("Creating a Cron Scheduler")
	return &crontab{
		cron: cron.New(),
	}
}

// Methods
func (c *crontab) StartCron() {
	log.Println("Start Cron")
	c.cron.Start()
}

func (c *crontab) StopCron() {
	log.Println("Stop Cron")
	c.cron.Stop()
}

func (c *crontab) AddJob(exp string, job func()) error {
	entryId, err := c.cron.AddFunc(exp, job)
	if err != nil {
		return err
	}
	c.entryIds = append(c.entryIds, entryId)
	return nil
}

func (c *crontab) Schedule(scheduleTime cron.Schedule, job cron.Job) {
	c.cron.Schedule(scheduleTime, job)
}

func (c *crontab) Entries() []cron.EntryID {
	return c.entryIds
}
