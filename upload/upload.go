package upload

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
)

func UploadCourse(conn *sql.Conn, path string, versionID uint64) error {
	folders, err := getChildrenFolders(path)
	if err != nil {
		return err
	}

	for _, folder := range folders {
		if strings.ToLower(folder.Name()) == "assets" || strings.ToLower(folder.Name()) == "resources" {
			log.Println("found assets/resources folder!")
			uploadMedia(conn, path+"/"+folder.Name(), versionID)
		} else {
			err1 := recursiveUploadSections(conn, folder.Name(), path+"/"+folder.Name(), versionID)
			if err1 != nil {
				return err1
			}
		}
	}

	return nil
}

func recursiveUploadSections(conn *sql.Conn, parentFolderName string, path string, versionID uint64) error {

	// save section
	row := conn.QueryRowContext(context.Background(), "INSERT INTO sections (name, version_id) VALUES ($1, $2) RETURNING id", parentFolderName, versionID)
	var sectionID uint64
	err := row.Scan(&sectionID)
	if err != nil {
		return err
	}

	contents, err1 := getChildrenMDFiles(path)
	if err1 != nil {
		return err1
	}

	if len(contents) != 0 {
		log.Println("yes contents for section", parentFolderName)

		for _, file := range contents {
			data, err2 := os.ReadFile(path + "/" + file.Name())
			if err2 != nil {
				log.Println("upload ERROR reading", path+"/"+file.Name(), "file:", err2)
				return err2
			}

			// strip .md off of name
			name := file.Name()[:len(file.Name())-3]

			_, err3 := conn.ExecContext(context.Background(), "INSERT INTO contents(language, section_id, markdown) VALUES($1, $2, $3)", name, sectionID, string(data))
			if err3 != nil {
				log.Println("upload ERROR inserting into contents for section", parentFolderName+":", err3)
				return err3
			}
		}

	} else {
		log.Println("no contents for section", parentFolderName)
	}

	childrenFolders, err4 := getChildrenFolders(path)
	if err4 != nil {
		log.Println("upload ERROR getting children folders:", err4)
		return err4
	}

	// freach children folder
	for _, folder := range childrenFolders {
		err5 := recursiveUploadSections(conn, folder.Name(), path+"/"+folder.Name(), versionID)
		if err5 != nil {
			log.Println("upload ERROR uploading sections for", folder.Name()+":", err5)
			return err5
		}
	}

	return nil
}

func uploadMedia(conn *sql.Conn, path string, versionID uint64) {

}
