package client

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	EnvHome = "SCRAPBOX_HOME"
)

func createQueryResultFile(host, project string, tags []string, skip, limit int) (*os.File, error) {

	directory := path.Join(getScrapboxHomeDir(), "query", trimPortFromHost(host), project, path.Join(tags...))
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "failed to make query cache directory")
	}
	filepath := path.Join(directory, EncodeFilename(fmt.Sprintf("%d-%d", skip, limit)))
	fout, err := os.Create(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create query cache file")
	}

	return fout, nil
}

func haveGoodQueryResultFile(host, project string, tags []string, skip, limit int, expiration time.Duration) bool {

	directory := path.Join(getScrapboxHomeDir(), "query", trimPortFromHost(host), project, path.Join(tags...))
	filepath := path.Join(directory, EncodeFilename(fmt.Sprintf("%d-%d", skip, limit)))
	finfo, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	if finfo.IsDir() {
		return false
	}
	mod := finfo.ModTime()
	now := time.Now()
	duration := now.Sub(mod)

	return duration <= expiration
}

func openQueryResultFile(host, project string, tags []string, skip, limit int) (*os.File, error) {

	directory := path.Join(getScrapboxHomeDir(), "query", trimPortFromHost(host), project, path.Join(tags...))
	filepath := path.Join(directory, EncodeFilename(fmt.Sprintf("%d-%d", skip, limit)))
	fin, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open query cache file")
	}

	return fin, nil
}

func createPageFile(host, project, page string) (*os.File, error) {

	directory := path.Join(getScrapboxHomeDir(), "page", trimPortFromHost(host), project)
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "failed to make page cache directory")
	}
	filepath := path.Join(directory, EncodeFilename(page))
	fout, err := os.Create(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create page cache file")
	}

	return fout, nil
}

func haveGoodPageFile(host, project, page string, expiration time.Duration) bool {

	directory := path.Join(getScrapboxHomeDir(), "page", trimPortFromHost(host), project)
	filepath := path.Join(directory, EncodeFilename(page))
	finfo, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	if finfo.IsDir() {
		return false
	}
	mod := finfo.ModTime()
	now := time.Now()
	duration := now.Sub(mod)

	return duration <= expiration
}

func openPageFile(host, project, page string) (*os.File, error) {

	directory := path.Join(getScrapboxHomeDir(), "page", trimPortFromHost(host), project)
	filepath := path.Join(directory, EncodeFilename(page))
	fin, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open page cache file")
	}

	return fin, nil
}

func getScrapboxHomeDir() string {
	value := os.Getenv(EnvHome)
	if len(value) == 0 {
		// TODO ohtomi: use go-homedir
		value = path.Join(os.Getenv("HOME"), ".scrapbox")
	}
	return value
}

func EncodeFilename(filename string) string {
	slashEscaped := strings.Replace(filename, "/", "%2F", -1)
	colonEscaped := strings.Replace(slashEscaped, ":", "%3A", -1)
	pipeEscaped := strings.Replace(colonEscaped, "|", "%7C", -1)
	return pipeEscaped
}

func trimPortFromHost(host string) string {
	if index := strings.Index(host, ":"); index == -1 {
		return host
	} else {
		return host[:index]
	}
}
