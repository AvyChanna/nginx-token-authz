package reader

import (
	"os"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/AvyChanna/nginx-token-authz/internal/app"
	"github.com/AvyChanna/nginx-token-authz/internal/rbac/auther"
	"github.com/AvyChanna/nginx-token-authz/internal/rbac/parser"
)

type Reader struct {
	fileName string
	config   *atomic.Pointer[auther.Auther]
	lastData string // string is immutable representation of bytes
}

func New(filename string) Reader {
	return Reader{
		fileName: filename,
		config:   &atomic.Pointer[auther.Auther]{},
		lastData: "",
	}
}

func (r *Reader) readConfig() error {
	data, err := os.ReadFile(r.fileName)
	if err != nil {
		return err
	}

	if r.lastData == string(data) {
		app.Get().Log().Debug("No changes in config file")
		return nil
	}

	app.Get().Log().Debug("Reading config file")

	var d parser.Config

	err = yaml.Unmarshal(data, &d)
	if err != nil {
		return err
	}

	autherObj, err := d.Done()
	if err == nil {
		r.config.Store(autherObj)
	}

	return err
}

func (r *Reader) ReadConfig() (*auther.Auther, error) {
	err := r.readConfig()
	return r.config.Load(), err
}

func (r *Reader) WatchConfig(pollDur time.Duration) *atomic.Pointer[auther.Auther] {
	err := r.readConfig()
	if err != nil {
		app.Get().Log().Error(err)
	}

	go func() {
		ticker := time.NewTicker(pollDur)
		defer ticker.Stop()

		for {
			select {
			case <-app.Get().Ctx().Done():
				return
			case <-ticker.C:
				err := r.readConfig()
				if err != nil {
					app.Get().Log().Errorf("error parsing file for changes = %s", err)
					continue
				}
			}
		}
	}()

	return r.config
}
