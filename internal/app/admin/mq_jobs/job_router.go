package mq_jobs

import "mogong/internal/app/admin/mq_jobs/job_functions/test"

var (
	jobsExecutor = map[string]interface{}{
		"user": map[string]func(interface{}) (interface{}, error){
			"hello": test.Hello,
		},
	}
)
