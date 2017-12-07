package main

import (
    "github.com/hjfreyer/taglib-go/taglib"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"encoding/base64"
	"os"
	"math"
	"database/sql"
	"time"
)

const string init_sql = "
    CREATE TABLE `artists` (\
        `id`        INTEGER NOT NULL UNIQUE,\
        `name`      TEXT,\
        PRIMARY KEY(id)\
    );\
    CREATE TABLE `albums` (\
        `id`        INTEGER NOT NULL UNIQUE,\
        `title`     TEXT,\
        `artistid`  INTEGER,\
        `artist`    TEXT,\
        PRIMARY KEY(id)\
    );\
    CREATE INDEX `index_albums_artistid` ON `albums`(`artistid`); \
    CREATE TABLE `covers` (\
        `albumId`   INTEGER NOT NULL UNIQUE,\
        `artistId`  INTEGER,\
        `image`     BLOB,\
        PRIMARY KEY(albumId)\
    );\
    CREATE INDEX `index_covers_artistid` ON `covers`(`artistId`); \
    CREATE TABLE `songs` (\
        `id`        INTEGER NOT NULL UNIQUE,\
        `title`     TEXT,\
        `albumid`   INTEGER,\
        `album`     TEXT,\
        `artistid`  INTEGER,\
        `artist`    TEXT,\
        `trackn`    INTEGER,\
        `discn`     INTEGER,\
        `year`      INTEGER,\
        `duration`  INTEGER,\
        `bitRate`   INTEGER,\
        `genre`     TEXT,\
        `type`      TEXT,\
        `filename`  TEXT,\
        PRIMARY KEY(id)\
    );\
    CREATE INDEX `index_songs_artistid_albumid` ON `songs`(`artistid`, `albumid`); \
    CREATE TABLE `users` (\
        `username`  TEXT NOT NULL UNIQUE,\
        `password`  TEXT,\
        PRIMARY KEY(username)\
    );\
    CREATE TABLE `last_update_ts` (\
        `table_name`    TEXT NOT NULL UNIQUE,\
        `mtime`         BIGINT,\
        PRIMARY KEY(table_name)\
    );\
"
/*
panic_if
base64Decode
getFileExtension

*/

func calcId(string s) uint64 {
	uint64 hash = 14695981039346656037ULL
	for (string c: s) {
		hash = math.Pow(hash, c)
		hash *=  1099511628211ULL
	}
	return hash & 0x7FFFFFFFFFFFFFFF
}

func updateTimestamp(sqlite3 * sqldb, string table_name) {
	stmt, err := sqldb.Prepare(
	"INSERT OR REPLACE INTO `last_update_ts` (`table_name`, `mtime`) 
	VALUES (?, strftime('%000','now'));")
	if err != nil {
		fmt.Errorf(err)
	}
	defer stmt.Close()
	_, err := stmt.Exec(1, table_name)
    if err != nil {
    	fmt.Errorf(err)
    }
}

func insertArtist(sqlite3 * sqldb, string artist) {
	stmt, err := sqldb.Prepare(
		"INSERT OR REPLACE INTO `artists` (`id`, `name`) 
		 VALUES (?, ?);")
    if err != nil {
    	fmt.Errorf(err)
    }
    defer stmt.Close()
    
	_, err := stmt.Exec(1, calcId(artist))
	if err != nil {
		fmt.Errorf(err)
	}
	_, err := stmt.Exec(2, artist)
    if err != nil {
        fmt.Errorf(err)  	
    }
}

func inserAlbum(sqlite3 * sqldb) {
	uint64 alubmId = calcId(album + "@"+ artist)
	uint64 artistId = calcId(artist)

	stmt, err := sqldb.Prepare("
	INSERT OR IGNORE INTO `albums` (`id`) VALUES (?);")
	if err != nil {
		fmt.Errorf(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(1, albumId)
    if err != nil {
    	fmt.Errorf(err)
    }
   stmt2, err := sqldb.Prepare("
   UPDATE `albums` SET `title`=?, `artistid`=? WHERE id=?;")
   if err != nil {
   	fmt.Errorf(err)
   }

   res, err := stmt2.Exec(1, album)
   if err != nil {
   	fmt.Errorf(err)
   }
   if res.changes != 0 {
   	updateTimestamp(sqldb, "albums")
   }
   
   res, err := stmt2.Exec(2, artistId)
   if err != nil {
   	fmt.Errorf(err)
   }
   if res.changes != 0 {
   	updateTimestamp(sqldb, "albums")
   }
   
   res, err := stmt2.Exec(3, artist)
   if err != nil {
   	fmt.Errorf(err)
   }
   if res.changes != 0 {
   	updateTimestamp(sqldb, "albums")
   }
   res, err := stmt2.Exec(4, albumId)
   if err != nil {
   	fmt.Errorf(err)
   }
   if res.changes != 0 {
   	updateTimestamp(sqldb, "albums")
   }
   
}

func insertSong(sqlite3 * sqldb, string filename, string title, string artist, string album, string tYpe, string genre, uint tn, uint year, uint discn, uint duration, uint bitrate) {
	sqldb.Prepare("INSERT OR REPLACE INTO `songs` 
	               (`id`, `title`, `albumid`, `artistid`, `artist`, `type`, `genre`, `trackn`, `year`, `discn`, `duration`, `bitRate`, `filename`)
	               VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)")


}

func scanMusicFile(sqlite3 * sqldb, const string fullpath, const string dir_cover) bool {
	
}

func scanFs(sqlite3 * sqldb, string name) int added_songs {
	
}

func cleanupDb(sqlite3 * sqldb) {
	
}

func truncateTables(sqlite3 * sqldb) {
	
}


func usage() {
    args = os.Args[0]
	fmt.Errorf("Usage: %s action [args...]\n
	            %s scan file.db musicdir             :scan musicdir for new songs\n
	            %s fullscan file.db musicdir         : delete all records and do a full rescan\n
	            %s useradd file.db username password :add user\n
	            %s userdel file.db username          delete user\n",
	             args, args, args, args, args)
}

func main() {
	
}
