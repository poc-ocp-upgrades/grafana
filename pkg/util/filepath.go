package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var WalkSkipDir = errors.New("skip this directory")

type WalkFunc func(resolvedPath string, info os.FileInfo, err error) error

func Walk(path string, followSymlinks bool, detectSymlinkInfiniteLoop bool, walkFn WalkFunc) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	var symlinkPathsFollowed map[string]bool
	var resolvedPath string
	if followSymlinks {
		resolvedPath = path
		if detectSymlinkInfiniteLoop {
			symlinkPathsFollowed = make(map[string]bool, 8)
		}
	}
	return walk(path, info, resolvedPath, symlinkPathsFollowed, walkFn)
}
func walk(path string, info os.FileInfo, resolvedPath string, symlinkPathsFollowed map[string]bool, walkFn WalkFunc) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if info == nil {
		return errors.New("Walk: Nil FileInfo passed")
	}
	err := walkFn(resolvedPath, info, nil)
	if err != nil {
		if info.IsDir() && err == WalkSkipDir {
			err = nil
		}
		return err
	}
	if resolvedPath != "" && info.Mode()&os.ModeSymlink == os.ModeSymlink {
		path2, err := os.Readlink(resolvedPath)
		if err != nil {
			return err
		}
		if symlinkPathsFollowed != nil {
			if _, ok := symlinkPathsFollowed[path2]; ok {
				errMsg := "Potential SymLink Infinite Loop. Path: %v, Link To: %v"
				return fmt.Errorf(errMsg, resolvedPath, path2)
			}
			symlinkPathsFollowed[path2] = true
		}
		info2, err := os.Lstat(path2)
		if err != nil {
			return err
		}
		return walk(path, info2, path2, symlinkPathsFollowed, walkFn)
	}
	if info.IsDir() {
		list, err := ioutil.ReadDir(path)
		if err != nil {
			return walkFn(resolvedPath, info, err)
		}
		var subFiles = make([]subFile, 0)
		for _, fileInfo := range list {
			path2 := filepath.Join(path, fileInfo.Name())
			var resolvedPath2 string
			if resolvedPath != "" {
				resolvedPath2 = filepath.Join(resolvedPath, fileInfo.Name())
			}
			subFiles = append(subFiles, subFile{path: path2, resolvedPath: resolvedPath2, fileInfo: fileInfo})
		}
		if containsDistFolder(subFiles) {
			err := walk(filepath.Join(path, "dist"), info, filepath.Join(resolvedPath, "dist"), symlinkPathsFollowed, walkFn)
			if err != nil {
				return err
			}
		} else {
			for _, p := range subFiles {
				err = walk(p.path, p.fileInfo, p.resolvedPath, symlinkPathsFollowed, walkFn)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	return nil
}

type subFile struct {
	path, resolvedPath	string
	fileInfo			os.FileInfo
}

func containsDistFolder(subFiles []subFile) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, p := range subFiles {
		if p.fileInfo.IsDir() && p.fileInfo.Name() == "dist" {
			return true
		}
	}
	return false
}
