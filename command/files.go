package command

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

func canonicalFilepath(directory, filename string) string {
	escapedFilename := strings.Replace(filename, "/", "%2F", -1)
	return path.Join(directory, escapedFilename)
}

func readLocalCache(host, project, page string) ([]string, error) {

	var lines []string

	directory := path.Join(scrapboxHome, "page", host, project)
	filepath := canonicalFilepath(directory, page)

	fin, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fin.Close()

	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func writeLocalCache(host, project, page string, lines []string) error {

	directory := path.Join(scrapboxHome, "page", host, project)
	filepath := canonicalFilepath(directory, page)

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return err
	}
	fout, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fout.Close()

	for _, line := range lines {
		fout.WriteString(fmt.Sprintf("%s\n", line))
	}

	return nil
}
