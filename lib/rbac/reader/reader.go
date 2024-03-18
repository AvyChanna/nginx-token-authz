package reader

import (
	"os"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/AvyChanna/nginx-token-authz/lib/app"
	"github.com/AvyChanna/nginx-token-authz/lib/rbac/auther"
	"github.com/AvyChanna/nginx-token-authz/lib/rbac/parser"
	"github.com/AvyChanna/nginx-token-authz/lib/set"
)

func ReadConfig(fileName string) (*auther.Auther, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var d parser.Config

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

func WatchConfig(filename string, pollDur time.Duration) *atomic.Pointer[auther.Auther] {
	autherPtr := &atomic.Pointer[auther.Auther]{}
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
