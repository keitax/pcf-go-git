package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gopkg.in/src-d/go-billy.v2/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/textvid", func(w http.ResponseWriter, r *http.Request) {
		fs := memfs.New()
		ms := memory.NewStorage()
		s, err := git.Clone(ms, fs, &git.CloneOptions{
			URL: "https://github.com/keitax/textvid.git",
		})
		if err != nil {
			fatal(w, err)
			return
		}
		ci, err := s.CommitObjects()
		if err != nil {
			fatal(w, err)
			return
		}
		for {
			c, err := ci.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				fatal(w, err)
				return
			}
			w.Write([]byte(fmt.Sprintf("%s <%s> %s", c.Author.Name, c.Author.Email, c.Message)))
			w.Write([]byte("---\n"))
		}
	})

	logrus.Info("Start serve")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logrus.Fatal(err)
		return
	}
}

func fatal(w http.ResponseWriter, err error) {
	logrus.Fatal(err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
