package cronjob

import "mogong/internal/pkg/transports/cron"

func CreateInitServersFn(cronJob *DefaultCronJobService) cron.InitServers {
	return map[string]func(){
		"hello": func() {
			go cronJob.Hello()
		},
	}
}
