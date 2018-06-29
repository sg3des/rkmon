package main

import (
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/Masterminds/vcs"
)

type Repository struct {
	repo vcs.Repo
}

//NewRepository initialize vcs object to repository
func NewRepository(remote string) (*Repository, error) {
	u, err := url.Parse(remote)
	if err != nil {
		return nil, err
	}

	local := filepath.Base(u.Path)

	repo, err := vcs.NewRepo(remote, local)
	if err != nil {
		return nil, err
	}

	return &Repository{repo: repo}, nil
}

//Get method download or update repository and return last commit time
func (rep *Repository) Get() (date time.Time, err error) {
	if rep.repo.CheckLocal() {
		err = rep.repo.Update()
	} else {
		err = rep.repo.Get()
	}
	if err != nil {
		return
	}

	return rep.repo.Date()
}

//OpenFile from downloaded repository
func (rep *Repository) OpenFile(filename string) (*os.File, error) {
	filename = filepath.Join(rep.repo.LocalPath(), filename)

	return os.Open(filename)
}
