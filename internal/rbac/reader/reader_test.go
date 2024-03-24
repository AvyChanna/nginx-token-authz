package reader

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AvyChanna/nginx-token-authz/internal/app"
)

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	pollTime     = 500 * time.Millisecond
	fileSyncTime = pollTime + (100 * time.Millisecond)
)

func RandFileName() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func copyFileData(srcFilePath string, filePath string) error {
	data, err := os.ReadFile(srcFilePath)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return f.Close()
}

// modified verison of os.CreateTemp
func createTemp(srcFilePath, tempFilePrefix string) (string, error) {
	dir := os.TempDir()

	try := 0
	for {
		name := fmt.Sprintf("%s_%s", tempFilePrefix, RandFileName())
		fullPath := path.Join(dir, name)
		f, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o600)
		if os.IsExist(err) {
			if try++; try < 10000 {
				continue
			}
			return "", &os.PathError{Op: "createtemp", Path: tempFilePrefix + "*", Err: os.ErrExist}
		}
		if err != nil {
			return "", err
		}

		err = f.Close()
		if err != nil {
			return "", err
		}

		err = copyFileData(srcFilePath, fullPath)
		return fullPath, err
	}
}

func TestWatchConfig(t *testing.T) {
	app.Init(context.Background(), true)

	t.Run("file reloads", func(t *testing.T) {
		// make temp file
		fName, err := createTemp("testdata/data1.yaml", "testfile")
		require.NoError(t, err)

		defer os.Remove(fName)

		// New reader
		r := New(fName)

		configPtr := r.WatchConfig(pollTime)
		require.NotNil(t, configPtr.Load())

		// perms loaded
		assert.True(t, configPtr.Load().IsValidClaim("uid1", "perm1"))
		assert.False(t, configPtr.Load().IsValidClaim("uid1", "perm2"))

		// truncate file
		err = os.Truncate(fName, 0)
		require.NoError(t, err)

		time.Sleep(fileSyncTime)

		// perms remain the same
		assert.True(t, configPtr.Load().IsValidClaim("uid1", "perm1"))
		assert.False(t, configPtr.Load().IsValidClaim("uid1", "perm2"))

		// new data
		err = copyFileData("testdata/data2.yaml", fName)
		require.NoError(t, err)

		time.Sleep(fileSyncTime)

		// perms changed
		assert.False(t, configPtr.Load().IsValidClaim("uid1", "perm1"))
		assert.True(t, configPtr.Load().IsValidClaim("uid1", "perm2"))
	})

	t.Run("bad file data", func(t *testing.T) {
		// make temp file
		fName, err := createTemp("testdata/invalid_data.yaml", "testfile")
		require.NoError(t, err)

		defer os.Remove(fName)

		// New reader
		r := New(fName)

		configPtr := r.WatchConfig(pollTime)
		require.Nil(t, configPtr.Load())

		// empty perms
		assert.False(t, configPtr.Load().IsValidClaim("uid1", "perm1"))

		// new data
		err = copyFileData("testdata/data1.yaml", fName)
		require.NoError(t, err)

		time.Sleep(fileSyncTime)

		// perms changed
		assert.True(t, configPtr.Load().IsValidClaim("uid1", "perm1"))
	})

	t.Run("do not refresh on errors", func(t *testing.T) {
		// make temp file
		fName, err := createTemp("testdata/data1.yaml", "testfile")
		require.NoError(t, err)

		defer os.Remove(fName)

		// New reader
		r := New(fName)

		configPtr := r.WatchConfig(pollTime)
		require.NotNil(t, configPtr.Load())

		// perms loaded
		assert.True(t, configPtr.Load().IsValidClaim("uid1", "perm1"))

		// new data
		err = copyFileData("testdata/invalid_data.yaml", fName)
		require.NoError(t, err)

		time.Sleep(fileSyncTime)

		// perms changed
		assert.True(t, configPtr.Load().IsValidClaim("uid1", "perm1"))
	})
}
