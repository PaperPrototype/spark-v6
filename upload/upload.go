package upload

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v4"
)

const MediaChunkSize int = 16384

func UploadCourse(conn *pgx.Conn, path string, versionID uint64) error {
	folders, err := getChildrenFolders(path)
	if err != nil {
		return err
	}

	for _, folder := range folders {
		if strings.ToLower(folder.Name()) == "assets" || strings.ToLower(folder.Name()) == "resources" {
			log.Println("found assets/resources folder!")
			err2 := uploadMedia(conn, path+"/"+folder.Name(), versionID)
			if err2 != nil {
				return err2
			}
		} else {
			err1 := recursiveUploadSections(conn, folder.Name(), path+"/"+folder.Name(), versionID)
			if err1 != nil {
				return err1
			}
		}
	}

	return nil
}

func recursiveUploadSections(conn *pgx.Conn, parentFolderName string, path string, versionID uint64) error {
	// save section
	row := conn.QueryRow(context.Background(), "INSERT INTO sections (name, version_id) VALUES ($1, $2) RETURNING id", parentFolderName, versionID)
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

			_, err3 := conn.Exec(context.Background(), "INSERT INTO contents(language, section_id, markdown) VALUES($1, $2, $3)", name, sectionID, string(data))
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

func uploadMedia(conn *pgx.Conn, path string, versionID uint64) error {
	mediaFiles, err := getMediaFiles(path)
	if err != nil {
		return err
	}

	for _, media := range mediaFiles {
		file, err1 := os.Open(path + "/" + media.Name())
		if err1 != nil {
			return err1
		}
		defer file.Close()

		row := conn.QueryRow(context.Background(), "INSERT INTO media (name, version_id, length, type) VALUES($1, $2, $3, $4) RETURNING id",
			media.Name(), versionID, 0, filepath.Ext(media.Name()))

		var mediaID uint64
		err2 := row.Scan(&mediaID)
		if err2 != nil {
			log.Println("upload ERROR scanning to get mediaID:", err2)
			return err2
		}

		buffer := make([]byte, MediaChunkSize)
		totalMediaLength, err3 := uploadChunkedMediaRecursive(buffer, conn, mediaID, bufio.NewReader(file), 1)
		if err3 != nil {
			log.Println("upload ERROR uploadingChunkedMedia:", err3)
			return err3
		}

		_, err4 := conn.Exec(context.Background(), "UPDATE media SET length = $1 WHERE id = $2", totalMediaLength, mediaID)
		if err4 != nil {
			log.Println("upload ERROR updating media length:", err4)
			return err4
		}

		log.Println("uploaded media of name:", media.Name())
	}

	return nil
}

// current var keeps track of a position number for the chunks
// current var is recursively increased each recursion
func uploadChunkedMediaRecursive(buffer []byte, conn *pgx.Conn, mediaID uint64, reader *bufio.Reader, current uint32) (int, error) {
	numBytesRead, eofErr := reader.Read(buffer)

	// if no error and some bytes were read
	if eofErr == nil || numBytesRead > 0 {
		// save bytes to db
		_, dbErr := conn.Exec(context.Background(), `INSERT INTO media_chunks (media_id, data, position) VALUES($1, $2, $3)`, fmt.Sprint(mediaID), buffer[:numBytesRead], current)
		if dbErr != nil {
			return numBytesRead, dbErr
		}

		numBytesRead1, err1 := uploadChunkedMediaRecursive(buffer, conn, mediaID, reader, current+1)
		// possibly a dbErr
		if err1 != nil {
			return numBytesRead + numBytesRead1, err1
		}
		return numBytesRead + numBytesRead1, nil
	}

	return numBytesRead, nil
}
