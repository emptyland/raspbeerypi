package model

import (
	"fmt"
	"time"
	"io"
	"os"
	"os/exec"
	"bufio"
	"log"
	"encoding/json"

	"api"
)

type JobModel struct {
	jobs     []jobContentVO
	jobsByID map[int]int
	jobEnv   *JobEnvDef
	uniqueID int
}

var _ = (api.Model)(&JobModel{})
var _ = time.Sleep
var _ = os.Open
var _ = log.Println

func NewJobService(jobEnv *JobEnvDef) *api.Service {
	jobModel := &JobModel{
		jobs:     loadJobs(jobEnv.Metadata),
		jobEnv:   jobEnv,
	}

	jobModel.syncIndex()
	return api.NewService(jobModel)
}

func (self *JobModel) Access(appKey string, token string) bool {
	return true
}

func (self *JobModel) GetApiJobList(res *jobContentRequest) error {
	res.Entries = self.jobs
	return nil
}

func (self *JobModel) PostApiJobContent(res *jobOperationResponse, req *jobContentRequest) error {
	for _, job := range req.Entries {

		if index, ok := self.jobsByID[job.ID]; ok {
			self.jobs[index] = job
		} else {
			self.uniqueID++
			job.ID = self.uniqueID

			self.jobs = append(self.jobs, job)
			self.syncIndex()
		}
	}

	storeJobs(self.jobEnv.Metadata, self.jobs)
	res.Ok = true

	return self.generate()
}

func (self *JobModel) PostApiJobRun(res *jobOperationResponse, req *jobContentRequest) error {

	res.Output = make([]string, 0)
	for _, job := range req.Entries {
		return self.runJob(res, &job)
	}

	return nil
}

func (self *JobModel) PostApiJobDelete(res *jobOperationResponse, req *jobContentRequest) error {
	newNum := 0
	for _, job := range req.Entries {
		if _, ok := self.jobsByID[job.ID]; ok {
			self.jobsByID[job.ID] = -1
			newNum++
		} else {
			return fmt.Errorf("can not find jobID: %d", job.ID)
		}
	}
	newNum = len(self.jobs) - newNum

	newJobs := make([]jobContentVO, newNum)
	i := 0
	for _, job := range self.jobs {
		if id := self.jobsByID[job.ID]; id != -1 {
			newJobs[i] = job
			i++
		}
	}
	self.jobs = newJobs
	self.syncIndex()

	return self.generate()
}

func (self *JobModel) generate() error {
	crontabFile, err := os.Create(self.jobEnv.Crontab)
	if err != nil {
		return err
	}
	defer crontabFile.Close()

	for _, job := range self.jobs {
		if !job.Enable {
			continue
		}
		stubFileName := fmt.Sprintf("%s/%03d.script", self.jobEnv.Pwd, job.ID)
		var stubFile *os.File
		if stubFile, err = os.Create(stubFileName); err != nil {
			return err
		}
		defer stubFile.Close()

		if _, err = stubFile.Write([]byte(job.Code)); err != nil {
			return err;
		}

		// cron user exec args
		fmt.Fprintf(crontabFile, "%s %s ", job.Cron, job.User)

		execv := self.jobEnv.Lang[job.Lang]
		for _, segment := range execv {
			fmt.Fprintf(crontabFile, "%s ", segment)
		}
		crontabFile.Write([]byte(stubFileName + "\n"))
	}

	return nil
}

func (self *JobModel) runJob(res *jobOperationResponse, job *jobContentVO) error {
	lang, ok := self.jobEnv.Lang[job.Lang]
	if !ok {
		return fmt.Errorf("Unsupport langugage: %s", job.Lang)
	}

	cmd := exec.Command(lang[0], lang[1:]...)
	var (
		stderr io.ReadCloser
		stdout io.ReadCloser
		stdin  io.WriteCloser
		err    error
	)
	if stderr, err = cmd.StderrPipe(); err != nil {
		log.Println("get stderr pipe fail:", err)
		return err
	}
	if stdout, err = cmd.StdoutPipe(); err != nil {
		log.Println("get stdout pipe fail:", err)
		return err
	}
	if stdin,  err = cmd.StdinPipe(); err != nil {
		log.Println("get stdin pipe fail:", err)
		return err
	}

	if err = cmd.Start(); err != nil {
		log.Println("command start:", lang, "fail:", err)
		return err
	}
	if _, err = stdin.Write([]byte(job.Code)); err != nil {
		return err
	}
	stdin.Close()

	res.Output = readLines(stdout, res.Output)
	res.Output = readLines(stderr, res.Output)

	if err = cmd.Wait(); err != nil {
		log.Println("command wait fail:", err)
	}
	res.Ok = cmd.ProcessState.Success()

	return nil
}

func readLines(reader io.Reader, output []string) []string {
	rd := bufio.NewReader(reader)
	for {
		if line, _, err := rd.ReadLine(); err != nil {
			break
		} else {
			output = append(output, string(line))
		}
	}
	return output
}

func (self *JobModel) syncIndex() {
	self.jobsByID = make(map[int]int)

	for i, job := range self.jobs {
		self.jobsByID[job.ID] = i

		if job.ID > self.uniqueID {
			self.uniqueID = job.ID
		}
	}
}

func loadJobs(fileName string) (rv []jobContentVO) {
	jobs := &jobContentPersistented {}

	rd, err := os.Open(fileName)
	if err != nil {
		log.Fatal("can not open:", fileName)
	}
	defer rd.Close()

	dd := json.NewDecoder(rd)
	if err = dd.Decode(jobs); err != nil {
		log.Fatal("can not decode:", fileName)
	}

	rv = jobs.Entries
	return
}

func storeJobs(fileName string, vos []jobContentVO) {
	jobs := &jobContentPersistented {
		Entries: vos,
	}

	wt, err := os.Create(fileName)
	if err != nil {
		log.Fatal("can not open:", fileName)
	}
	defer wt.Close()

	ed := json.NewEncoder(wt)
	if err = ed.Encode(jobs); err != nil {
		log.Fatal("can not encode:", err)
	}
}
