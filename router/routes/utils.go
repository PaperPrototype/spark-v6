package routes

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func notFound(c *gin.Context) {
	c.Redirect(http.StatusFound, "/lost")
}

func WriteMediaChunks(conn *pgxpool.Pool, writer io.Writer, mediaID uint64) {
	row := conn.QueryRow(context.Background(), "SELECT data FROM media_chunks WHERE media_id = $1 ORDER BY position", mediaID)

	buffer := []byte{}
	err := row.Scan(&buffer)
	if err != nil {
		log.Println("ERROR scanning data from media_chunks:", err)
		return
	}

	num, err1 := writer.Write(buffer)
	if num <= 0 {
		return
	} else if err1 != nil {
		return
	}

	writeMediaChunk(conn, writer, mediaID, 1)
}

// this is a recursive function
func writeMediaChunk(conn *pgxpool.Pool, writer io.Writer, mediaID uint64, current int) {
	row := conn.QueryRow(context.Background(), "SELECT data FROM media_chunks WHERE media_id = $1 ORDER BY position OFFSET $2", mediaID, current)

	buffer := []byte{}
	err := row.Scan(&buffer)
	if err != nil {
		log.Println("ERROR scanning data from media_chunks:", err)
		return
	}

	num, err1 := writer.Write(buffer)
	if num <= 0 {
		log.Println("database: Ignore the above error \n The above query error just means that there are no more chunks for the image in the db. ")
		return
	} else if err1 != nil {
		return
	}

	writeMediaChunk(conn, writer, mediaID, current+1)
}
