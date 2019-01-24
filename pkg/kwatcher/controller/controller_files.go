package controller

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/juju/errors"
)

func getCurrentFiles(dir string) (map[string][]byte, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Annotatef(err, "incorrect dir: %s", dir)
	}

	out := make(map[string][]byte, len(files))
	for _, file := range files {
		if name := file.Name(); filepath.Ext(name) == ".vcl" {
			contents, err := ioutil.ReadFile(filepath.Join(dir, name))
			if err != nil {
				return nil, errors.Annotatef(err, "problem reading file %s", name)
			}
			out[name] = contents
		}
	}
	return out, nil
}

func (r *ReconcileVarnish) reconcileFiles(dir string, currFiles map[string][]byte, newFiles map[string][]byte) (bool, error) {
	diffFiles := make(map[string]int, len(newFiles))
	for k := range newFiles {
		diffFiles[k] = 1
	}
	for k := range currFiles {
		diffFiles[k] = diffFiles[k] - 1
	}

	filesTouched := false
	for fileName, status := range diffFiles {
		fullpath := filepath.Join(dir, fileName)
		if status == -1 {
			filesTouched = true
			r.logger.Infow("Removing file", "path", fullpath)
			if err := os.Remove(fullpath); err != nil {
				return filesTouched, errors.Annotatef(err, "could not delete file %s", fullpath)
			}
		} else if status == 0 && bytes.Compare(currFiles[fileName], newFiles[fileName]) != 0 {
			filesTouched = true
			if err := ioutil.WriteFile(fullpath, newFiles[fileName], 0644); err != nil {
				return filesTouched, errors.Annotatef(err, "could not write file %s", fullpath)
			}
			r.logger.Infow("Rewriting file", "path", fullpath)
		} else if status == 1 {
			filesTouched = true
			if err := ioutil.WriteFile(fullpath, newFiles[fileName], 0644); err != nil {
				return filesTouched, errors.Annotatef(err, "could not write file %s", fullpath)
			}
			r.logger.Infow("Writing new file", "path", fullpath)
		}
	}
	return filesTouched, nil
}
