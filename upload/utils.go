package upload

import (
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func getChildrenFolders(path string) ([]fs.FileInfo, error) {
	fsInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return fsInfos, err
	}

	// slice (go list) holding all folders
	var folders []fs.FileInfo

	// foreach item in directory
	for _, fsInfo := range fsInfos {
		if fsInfo.IsDir() && !isFolderNameToIgnore(fsInfo.Name()) {
			folders = append(folders, fsInfo)
		}
	}

	// give back folders
	return folders, nil
}

var folderNamesToIgnore []string = []string{
	".git",
	"ignore",
	"archive",
}

func isFolderNameToIgnore(name string) bool {
	for _, nameIgnore := range folderNamesToIgnore {
		if strings.EqualFold(nameIgnore, name) || strings.Contains(strings.ToLower(name), "ignore") {
			return true
		}
	}

	return false
}

var mdFileNamesToIgnore []string = []string{
	"license.md",
	"readme.md",
	"license",
	"readme",
}

func mdNameToIgnore(name string) bool {
	for _, nameToIgnore := range mdFileNamesToIgnore {
		// if match, then ignore = true
		if strings.EqualFold(name, nameToIgnore) {
			return true
		}
	}
	return false
}

// get all children files in a particular folder
func getChildrenMDFiles(path string) ([]fs.FileInfo, error) {
	log.Println("getting children .md files for", path)

	// read the directory's contents
	fsInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return fsInfos, err
	}

	// slice (go list) variable for all files
	var files []fs.FileInfo

	// foreach item in direcotry
	for _, fsInfo := range fsInfos {
		// if it is a file
		if !fsInfo.IsDir() && filepath.Ext(fsInfo.Name()) == ".md" && !mdNameToIgnore(fsInfo.Name()) {
			// add current folder to folders slice
			files = append(files, fsInfo)
		}
	}

	// give back only files in direcotry
	return files, nil
}

var allowedMediaFiles []string = []string{
	".gif",
	".png",
	".jpg",
	".zip",
}

func isAllowedMediaFile(fileName string) bool {
	ext := filepath.Ext(fileName)

	for _, name := range allowedMediaFiles {
		if name == ext {
			return true
		}
	}

	return false
}

// get all children files in a particular folder
func getMediaFiles(path string) ([]fs.FileInfo, error) {
	log.Println("worker: getting media files for", path)

	// read the directory's contents
	fsInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return fsInfos, err
	}

	// slice (go list) variable for all files
	var files []fs.FileInfo

	// foreach item in direcotry
	for _, fsInfo := range fsInfos {
		// if it is a file
		if !fsInfo.IsDir() && isAllowedMediaFile(fsInfo.Name()) {
			// add current folder to folders slice
			files = append(files, fsInfo)
		}
	}

	// give back only files in direcotry
	return files, nil
}
