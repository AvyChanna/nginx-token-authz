package rbac

import (
	"os"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/AvyChanna/nginx-token-authz/lib/app"
	"github.com/AvyChanna/nginx-token-authz/lib/set"
)

func ReadConfig(fileName string) (*Auther, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var d Config

	err = yaml.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}

	for _, user := range d.Users {
		if user.Groups == nil {
			user.Groups = &set.StrSet{}
		}
	}

	if d.AllowAllSet == nil {
		d.AllowAllSet = &set.StrSet{}
	}

	return d.Done()
}

func WatchConfig(filename string, pollDur time.Duration) *atomic.Pointer[Auther] {
	autherPtr := &atomic.Pointer[Auther]{}
	newAuther, _ := ReadConfig(filename)
	autherPtr.Store(newAuther)

	go func() {
		ticker := time.NewTicker(pollDur)
		defer ticker.Stop()

		for {
			select {
			case <-app.Get().Ctx().Done():
				return
			case <-ticker.C:
				newAuther, err := ReadConfig(filename)
				if err != nil {
					app.Get().Log().Errorf("error parsing file for changes = %s", err)
					continue
				}
				autherPtr.Store(newAuther)
			}
		}
	}()

	return autherPtr
}
