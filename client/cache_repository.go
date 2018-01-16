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

func ScrapboxHomeFromEnv() string {
	value := os.Getenv(EnvHome)
	if len(value) == 0 {
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

func createQueryResultFile(homeDir, host, project string, tags []string, skip, limit int) (*os.File, error) {

	directory := path.Join(homeDir, "query", trimPortFromHost(host), project, path.Join(tags...))
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

func haveGoodQueryResultFile(homeDir, host, project string, tags []string, skip, limit int, expiration time.Duration) bool {

	directory := path.Join(homeDir, "query", trimPortFromHost(host), project, path.Join(tags...))
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

func openQueryResultFile(homeDir, host, project string, tags []string, skip, limit int) (*os.File, error) {

	directory := path.Join(homeDir, "query", trimPortFromHost(host), project, path.Join(tags...))
	filepath := path.Join(directory, EncodeFilename(fmt.Sprintf("%d-%d", skip, limit)))
	fin, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open query cache file")
	}

	return fin, nil
}

func createPageFile(homeDir, host, project, page string) (*os.File, error) {

	directory := path.Join(homeDir, "page", trimPortFromHost(host), project)
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

func haveGoodPageFile(homeDir, host, project, page string, expiration time.Duration) bool {

	directory := path.Join(homeDir, "page", trimPortFromHost(host), project)
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

func openPageFile(homeDir, host, project, page string) (*os.File, error) {

	directory := path.Join(homeDir, "page", trimPortFromHost(host), project)
	filepath := path.Join(directory, EncodeFilename(page))
	fin, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open page cache file")
	}

	return fin, nil
}
