package controller

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/ibm/varnish-operator/pkg/logger"

	"github.com/pkg/errors"
)

func getCurrentFiles(dir string) (map[string]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "incorrect dir: %s", dir)
	}

	out := make(map[string]string, len(files))
	for _, file := range files {
		if name := file.Name(); filepath.Ext(name) == ".vcl" {
			contents, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				return nil, errors.Wrapf(err, "problem reading file %s", name)
			}
			out[name] = string(contents)
		}
	}
	return out, nil
}

func (r *ReconcileVarnish) reconcileFiles(ctx context.Context, dir string, currFiles map[string]string, newFiles map[string]string) (bool, error) {
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
		logr := logger.FromContext(ctx).With(logger.FieldFilePath, fullpath)
		if status == -1 {
			filesTouched = true
			logr.Infow("Removing file")
			if err := os.Remove(fullpath); err != nil {
				return filesTouched, errors.Wrapf(err, "could not delete file %s", fullpath)
			}
		} else if status == 0 && strings.Compare(currFiles[fileName], newFiles[fileName]) != 0 {
			filesTouched = true
			if err := os.WriteFile(fullpath, []byte(newFiles[fileName]), 0644); err != nil {
				return filesTouched, errors.Wrapf(err, "could not write file %s", fullpath)
			}
			logr.Infow("Rewriting file")
		} else if status == 1 {
			filesTouched = true
			if err := os.WriteFile(fullpath, []byte(newFiles[fileName]), 0644); err != nil {
				return filesTouched, errors.Wrapf(err, "could not write file %s", fullpath)
			}
			logr.Infow("Writing new file")
		}
	}
	return filesTouched, nil
}
