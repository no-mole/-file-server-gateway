package bootstrap

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"syscall"
	"time"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/cron"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/utils"
)

type job struct {
	Name string
	Cron string
}

func (j *job) Key() string {
	return fmt.Sprintf("task-%d", j.Name)
}

func (j *job) Spec() string {
	return j.Cron
}

type FileAtime struct {
	bucket   string
	fileName string
	atime    int64
}

func (j *job) Run() {
	fileAtime := make([]*FileAtime, 0)

	paths := path.Join(utils.GetCurrentAbPath(), "data")
	buckets, _ := ioutil.ReadDir(paths)
	for _, bucket := range buckets {
		if bucket.IsDir() {
			files, _ := ioutil.ReadDir(path.Join(paths, bucket.Name()))
			for _, file := range files {
				fi, err := os.Stat(path.Join(paths, bucket.Name(), file.Name()))
				if err != nil {
					return
				}
				stat := fi.Sys().(*syscall.Stat_t)
				atime := time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec).Unix()
				fileAtime = append(fileAtime, &FileAtime{
					bucket:   bucket.Name(),
					fileName: file.Name(),
					atime:    atime,
				})
			}
		}
	}

	if len(fileAtime) < 1025 {
		return
	}

	sort.Slice(fileAtime, func(i, j int) bool {
		return fileAtime[i].atime < fileAtime[j].atime
	})

	remain := len(fileAtime) - 1024

	for i := 0; i < remain; i++ {
		_ = os.RemoveAll(path.Join(paths, fileAtime[i].bucket, fileAtime[i].fileName))
	}
}

func InitCleanFile(_ context.Context) error {
	errCh := make(chan error, 1)
	c := cron.New(cron.WithErrChan(errCh))
	j := &job{
		Name: "clear_file",
		Cron: "* * * * *",
	}
	go c.Add(j)
	c.Start()
	return nil
}
