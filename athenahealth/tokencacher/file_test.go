package tokencacher

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFile_Get(t *testing.T) {
	assert := assert.New(t)

	file, err := ioutil.TempFile("", "go-athenahealth_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	c := &fileCache{
		Token:     "foo",
		ExpiresAt: time.Now().Add(time.Minute * 1),
	}

	b, _ := json.Marshal(c)
	ioutil.WriteFile(file.Name(), b, 0644)

	cacher := NewFile(file.Name())

	token, err := cacher.Get()

	assert.Equal(c.Token, token)
	assert.NoError(err)
}

func TestFile_Get_ErrTokenNotExist(t *testing.T) {
	assert := assert.New(t)

	file, err := ioutil.TempFile("", "go-athenahealth_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	cacher := NewFile(file.Name())

	token, err := cacher.Get()

	assert.Empty(token)
	assert.Error(err)
	assert.True(errors.Is(err, ErrTokenNotExist))
}

func TestFile_Get_ErrTokenExpired(t *testing.T) {
	assert := assert.New(t)

	file, err := ioutil.TempFile("", "go-athenahealth_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	c := &fileCache{
		Token:     "foo",
		ExpiresAt: time.Now().Add(-time.Minute * 1),
	}

	b, _ := json.Marshal(c)
	ioutil.WriteFile(file.Name(), b, 0644)

	cacher := NewFile(file.Name())

	token, err := cacher.Get()

	assert.Empty(token)
	assert.Error(err)
	assert.True(errors.Is(err, ErrTokenExpired))
}

func TestFile_Set(t *testing.T) {
	assert := assert.New(t)

	file, err := ioutil.TempFile("", "go-athenahealth_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	cacher := NewFile(file.Name())

	token := "foo"
	expiresAt := time.Now().Add(time.Minute * 1)

	err = cacher.Set(token, expiresAt)

	b, _ := ioutil.ReadFile(file.Name())
	c := &fileCache{}
	json.Unmarshal(b, c)

	assert.NoError(err)
	assert.Equal(token, c.Token)
	assert.True(expiresAt.Equal(c.ExpiresAt))
}
