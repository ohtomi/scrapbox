package client

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

const (
	EnvHome = "SCRAPBOX_HOME"
)

func createQueryResultFile(host, project string, tags []string, skip, limit int) (*os.File, error) {

	baseDir := path.Join(getScrapboxHomeDir(), "query", trimPortFromHost(host), project, path.Join(tags...))
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "failed to make query cache directory")
	}
	queryResultFilePath := path.Join(baseDir, EncodeFilename(fmt.Sprintf("%d-%d", skip, limit)))
	queryResultFile, err := os.Create(queryResultFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create query cache file")
	}

	return queryResultFile, nil
}

func haveGoodQueryResultFile(host, project string, tags []string, skip, limit int, expiration time.Duration) bool {

	baseDir := path.Join(getScrapboxHomeDir(), "query", trimPortFromHost(host), project, path.Join(tags...))
	queryResultFilePath := path.Join(baseDir, EncodeFilename(fmt.Sprintf("%d-%d", skip, limit)))
	fs, err := os.Stat(queryResultFilePath)
	if err != nil {
		return false
	}
	if fs.IsDir() {
		return false
	}
	mod := fs.ModTime()
	now := time.Now()
	duration := now.Sub(mod)

	return duration <= expiration
}

func openQueryResultFile(host, project string, tags []string, skip, limit int) (*os.File, error) {

	baseDir := path.Join(getScrapboxHomeDir(), "query", trimPortFromHost(host), project, path.Join(tags...))
	queryResultFilePath := path.Join(baseDir, EncodeFilename(fmt.Sprintf("%d-%d", skip, limit)))
	queryResultFile, err := os.Open(queryResultFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open query cache file")
	}

	return queryResultFile, nil
}

func createPageFile(host, project, page string) (*os.File, error) {

	baseDir := path.Join(getScrapboxHomeDir(), "page", trimPortFromHost(host), project)
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "failed to make page cache directory")
	}
	pageFilePath := path.Join(baseDir, EncodeFilename(page))
	pageFile, err := os.Create(pageFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create page cache file")
	}

	return pageFile, nil
}

func haveGoodPageFile(host, project, page string, expiration time.Duration) bool {

	baseDir := path.Join(getScrapboxHomeDir(), "page", trimPortFromHost(host), project)
	pageFilePath := path.Join(baseDir, EncodeFilename(page))
	fs, err := os.Stat(pageFilePath)
	if err != nil {
		return false
	}
	if fs.IsDir() {
		return false
	}
	mod := fs.ModTime()
	now := time.Now()
	duration := now.Sub(mod)

	return duration <= expiration
}

func openPageFile(host, project, page string) (*os.File, error) {

	baseDir := path.Join(getScrapboxHomeDir(), "page", trimPortFromHost(host), project)
	pageFilePath := path.Join(baseDir, EncodeFilename(page))
	pageFile, err := os.Open(pageFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open page cache file")
	}

	return pageFile, nil
}

func getScrapboxHomeDir() string {
	value := os.Getenv(EnvHome)
	if len(value) == 0 {
		if userHomeDir, err := homedir.Dir(); err != nil {
			value = path.Join(userHomeDir, ".scrapbox")
		} else {
			value = "./.scrapbox"
		}
	} else {
		if expanded, err := homedir.Expand(value); err != nil {
			value = expanded
		}
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
