package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type Album struct {
	Id     int
	Title  string
	Artist string
	Price  float32
}

func PrintError(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func AddAlbum(ctx context.Context, conn *pgx.Conn, album Album) {
	result, err := conn.Exec(ctx, `insert into album(title, artist, price) values($1, $2, $3)`, album.Title, album.Artist, album.Price)
	PrintError("error executing insert statement", err)
	if result.Insert() == true {
		log.Printf("added an album: %s to the database\n", album.Title)
	}
}

func DeleteAlbum(ctx context.Context, conn *pgx.Conn, id int) {
	result, err := conn.Exec(ctx, `delete from album where id = $1`, id)
	PrintError("error deleting from database: ", err)
	if result.Delete() == true {
		log.Printf("delete album with id: %d from the database.\n", id)
	}
}

func UpdateAlbum(ctx context.Context, conn *pgx.Conn, album Album, id int) {
	result, err := conn.Exec(ctx, `update album set title = $1, artist = $2, price= $3 where id = $4`, album.Title, album.Artist, album.Price, id)
	PrintError("error updating an album:", err)
	if result.Update() == true {
		log.Printf("updated album with id: %d\n", id)
	}
}

func GetAlbum(ctx context.Context, conn *pgx.Conn, id int) {
	var album Album
	row := conn.QueryRow(ctx, `select * from album where id = $1`, id)
	err := row.Scan(&album.Id, &album.Title, &album.Artist, &album.Price)
	PrintError("error scanning row query:", err)
	log.Printf("%v", album)

}

func GetAllAlbums(ctx context.Context, conn *pgx.Conn) {
	var albums []Album
	rows, err := conn.Query(ctx, `select * from album`)
	PrintError("error querying rows:", err)
	for rows.Next() {
		var album Album
		err := rows.Scan(&album.Id, &album.Title, &album.Artist, &album.Price)
		PrintError("error scanning rows:", err)
		albums = append(albums, album)
	}

	if err = rows.Err(); err != nil {
		PrintError("error collected during scanning rows:", err)
	}

	log.Printf("%v", albums)
}

//Note: crud is done, add api routes

func main() {
	connStr := "postgres://postgres:password@localhost:5432/recordings"
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	PrintError("couldn't connect to database:", err)

	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS album(
        id SERIAL PRIMARY KEY NOT NULL,
        title VARCHAR(128) NOT NULL,
        artist VARCHAR(255) NOT NULL,
        price NUMERIC(5,2) NOT NULL
    )`)
	PrintError("can't execute sql query:", err)

}
