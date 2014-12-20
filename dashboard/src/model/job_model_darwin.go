package model

import (
	"api"
	"fmt"
)

type JobModel struct {
	jobs []jobContentVO
	jobsByID map[int]*jobContentVO
}

var _ = (api.Model)(&JobModel{})

func NewJobService() *api.Service {
	jobModel := &JobModel{
		jobs: make([]jobContentVO, 2),
		jobsByID: make(map[int]*jobContentVO),
	}

	jobModel.jobs[0] = jobContentVO{
		ID:    0,
		Title: "Cron Job 1",
		Desc:  "For testing",
		Lang:  "bash",
		Code: `for i in *; do
    echo $i
done

exit 1`,
		Cron:   "* * * * *",
		Enable: true,
	}
	jobModel.jobsByID[0] = &jobModel.jobs[0];

	jobModel.jobs[1] = jobContentVO{
		ID:    1,
		Title: "Cron Job 2",
		Desc:  "For testing",
		Lang:  "python",
		Code: `import sys
def main():
    print sys.argv
if __name__ == 'main':
    main()
`,
		Cron:   "1 * * * *",
		Enable: false,
	}
	jobModel.jobsByID[1] = &jobModel.jobs[1];

	return api.NewService(jobModel)
}

func (self *JobModel) Access(appKey string, token string) bool {
	return true
}

func (self *JobModel) GetApiJobList(res *jobContentRequest) error {

	res.Entries = self.jobs
	return nil
}

func (self *JobModel) PostApiJobContent(res *operationResponse, req *jobContentRequest) error {
	for _, job := range req.Entries {

		if addr, ok := self.jobsByID[job.ID]; ok {
			*addr = job
		} else {
			return fmt.Errorf("Can not find job id: %d", job.ID)
		}
	}

	res.Ok = true
	res.Msg = "ok"

	return nil
}
