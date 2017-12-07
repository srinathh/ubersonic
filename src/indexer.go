package main

import (
    "github.com/dhowden/tag"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"encoding/base64"
	"os"
	"math"
	"database/sql"
	"time"
	"strings"
	"path"
	"path/filepath"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
)

init_sql const string

init_sql = "
    CREATE TABLE `artists` (
        `id`        INTEGER NOT NULL UNIQUE,
        `name`      TEXT,
        PRIMARY KEY(id)
    );
    CREATE TABLE `albums` (
        `id`        INTEGER NOT NULL UNIQUE,
        `title`     TEXT,
        `artistid`  INTEGER,
        `artist`    TEXT,
        PRIMARY KEY(id)
    );
    CREATE INDEX `index_albums_artistid` ON `albums`(`artistid`); 
    CREATE TABLE `covers` (
        `albumId`   INTEGER NOT NULL UNIQUE,
        `artistId`  INTEGER,
        `image`     BLOB,
        PRIMARY KEY(albumId)
    );
    CREATE INDEX `index_covers_artistid` ON `covers`(`artistId`); 
    CREATE TABLE `songs` (
        `id`        INTEGER NOT NULL UNIQUE,
        `title`     TEXT,
        `albumid`   INTEGER,
        `album`     TEXT,
        `artistid`  INTEGER,
        `artist`    TEXT,
        `trackn`    INTEGER,
        `discn`     INTEGER,
        `year`      INTEGER,
        `duration`  INTEGER,
        `bitRate`   INTEGER,
        `genre`     TEXT,
        `type`      TEXT,
        `filename`  TEXT,
        PRIMARY KEY(id)
    );
    CREATE INDEX `index_songs_artistid_albumid` ON `songs`(`artistid`, `albumid`); 
    CREATE TABLE `users` (
        `username`  TEXT NOT NULL UNIQUE,
        `password`  TEXT,
        PRIMARY KEY(username)
    );
    CREATE TABLE `last_update_ts` (
        `table_name`    TEXT NOT NULL UNIQUE,
        `mtime`         BIGINT,
        PRIMARY KEY(table_name)
    );
"
/*
panic_if
base64Decode
getFileExtension

*/

type AudioInfo struct {
	Filename string
	BitRate int
    SampleRate int
    Duration float64
    Length float64
} 

func getAudioInfo (audiopath string) AudioInfo {
	key string
	path = filepath.Dir(audiopath)
	multipartFormFile *multipart.Form.File

	
	if fhs := multipartFormFile["audio"]; len(fhs) > 0 {
		f, err := fhs[0].Open()
		if err != nil {
			fmt.Errorf(err)
		}
		defer f.Close()		
	}

	FileName string = fhs[0].FileName
	BitRate int = int(f.BitRate)
	SampleRate int = int(f.SampleRate)
	Duration float64 = float64(f.Data[0]/f.SampleRate)

	audio = AudioInfo{
		FileName: FileName,
		BitRate: BitRate,
		SampleRate: SampleRate,
		Duration: Duration,
	}
	return audio
}

func calcId(s string) uint64 {
	uint64 hash = 14695981039346656037ULL
	for (string c: s) {
		hash = math.Pow(hash, c)
		hash *=  1099511628211ULL
	}
	return hash & 0x7FFFFFFFFFFFFFFF
}

func updateTimestamp(sqldb * sqlite3, table_name string) {
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

func insertArtist(sqldb *sqlite3, artist string) {

    sqldb.RegisterUpdateHook(updateTimestamp(sqldb, "artists"))
        
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

func inserAlbum(sqldb *sqlite3, arubum string, artist string) {
	uint64 alubmId = calcId(album + "@"+ artist)
	uint64 artistId = calcId(artist)

    sqldb.RegisterUpdateHook(updateTimestamp(sqldb, "albums"))
       
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
   
   defer stmt2.Close()
   
   _, err := stmt2.Exec(1, album)
   if err != nil {
   	fmt.Errorf(err)
   }
   
   _, err := stmt2.Exec(2, artistId)
   if err != nil {
   	fmt.Errorf(err)
   }
   
   _, err := stmt2.Exec(3, artist)
   if err != nil {
   	fmt.Errorf(err)
   }
   
   _, err := stmt2.Exec(4, albumId)
   if err != nil {
   	fmt.Errorf(err)
   }
   
}

func insertSong(sqldb *sqlite3, filename string, title string, artist string, album string, tYpe string, genre string, tn uint, year uint, discn uint, duration uint, bitrate uint) {
    sqldb.RegisterUpdateHook(updateTimestamp(sqldb, "songs"))
	stmt, err := sqldb.Prepare("
	INSERT OR REPLACE INTO `songs` 
	 (`id`, `title`, `albumid`, `artistid`, `artist`, `type`, `genre`, `trackn`, `year`, `discn`, `duration`, `bitRate`, `filename`)
	 VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	 ")
	if err != nil {
		fmt.Errorf(err)
	}
	defer stmt.Close()

    _, err := stmt.Exec(1, calcId(string(tn) + "@" + string(discn) + "@" + title + "@" + album + "@" + artist))
    if err != nil {
    	fmt.Errorf(err)
    }
    _, err != stmt.Exec(2, title)
    if err != nil {
    	fmt.Errorf(err)
    }
    
    _, err != stmt.Exec(3, calcId(album + "@" + artist))
    if err != nil {
        fmt.Errorf(err)
    }
    _, err != stmt.Exec(4, album)
    if err != nil {
    	fmt.Errorf(err)
    }
    
    _, err != stmt.Exec(5, calcId(artist))
    if err != nil {
    	fmt.Errorf(err)
    }
    _, err != stmt.Exec(6, artist)
    if err != nil {
    	fmt.Errorf(err)
    }
    
    _, err != stmt.Exec(7, tYpe)
    if err != nil {
    	fmt.Errorf(err)
    }
    
    _, err != stmt.Exec(8, genre)
    if err != nil {
    	fmt.Errorf(err)
    }
    
    _, err != stmt.Exec(9, tn)
    if err != nil {
    	fmt.Errorf(err)
    }
    
    _, err != stmt.Exec(10, year)
    if err != nil {
    	fmt.Errorf(err)
    }

    _, err != stmt.Exec(11, discn)
    if err != nil {
    	fmt.Errorf(err)
    }

    _, err != stmt.Exec(12, duration)
    if err != nil {
    	fmt.Errorf(err)
    }

    _, err != stmt.Exec(13, bitrate)
    if err != nil {
    	fmt.Errorf(err)
    }

    _, err != stmt.Exec(14, filename)
    if err != nil {
    	fmt.Errorf(err)
    }

}
    
func scanMusicFile(sqldb *sqlite3, fullpath const string, dir_cover const string) bool {
	artist string
	albumartist string
	album string
	title string
	genre string
    track int
    year int
    discn int
    bitrate uint
    length uint
    
	if ext != "mp3" && ext != "ogg" && ext != "flac" && ext != "m4a" && ext != "mp4 "{
		return false
	}
	f, err := os.Open(filepath.Dir(fullpath))
	if err != nil {
		return false
	}
	defer f.Close()
	
	m, err := tag.ReadFrom(f)
	if err != nil {
		return false
	}
    
    *Picture coverImage
    string cover
    coverImage = m.Picture()
    switch ext {
    case: "mp3"
     cover = base64.stdEncoding.EncodeToString(coverImage.)
    case: "ogg"
    case: "flac"
    case: "mp4"
    case: "m4a"
    default:
    fmt.Errorf("This file can't index with me.")
    }
    albumartist = strings.Trim(m.Albumartist(), " ")
    
    if len(albumartist) != 0 {
    	srtist = albumartist
    } else {
    	artist = strings.Trim(m.Artist, " ")
    }
    album = strings.Trim(m.Album(), " ")
    title = strings.Trim(m.Title(), " ")

    discn = 0
    if 
    
    genre = strings.Trim(m.Genre())
    year = m.Year()
    track = m.Track()
    
   _, err := sqldb.Begin()
   if err != nil {
   	fmt.Errorf(err)
   }
   defer sqldb.Close()

   insertSong(sqldb, fullpath, title, artist, album, ext, genre, track, year, discn, length, bitrate)

   insertAlbum(sqldb, album, artist, cover)

   inserArtist(sqldb, artist)

   fmt.Println("added")

   return true
}

func scanFs(sqldb *sqlite3, name string) added_songs int{
	dirent entry struct
	dir DIR
	statbuf struct stat
	added_songs int
	added_songs = 0
	coverfile ifstream
	cover_data string


	return added_sings
	
}

func cleanupDb(sqldb *sqlite3) {
	stmt, err := sqldb.Prepare("
	DELETE FROM `albums` WHERE ( SELECT count(id) FROM `songs` WHERE albumId = albums.id) = 0;
	")
	if err != nil {
		fmt.Errorf(err)
	}
	defer stmt.Close()

    
    stmt2, err := sqldb.Prepare("
	DELETE FROM `artists` WHERE ( SELECT count(id) FROM `albums` WHERE artistId = artists.id) = 0;
    ")
	if err != nil {
		fmt.Errorf(err)
	}
	defer stmt2.Close()

	stmt3, err := sqldb.Prepare("
	DELETE FROM `covers` WHERE ( SELECT count(id) FROM `albums` WHERE albumId = covers.albumId) = 0;
	")
	if err != nil {
		fmt.Errorf(err)
	}
	defer stmt3.Close()

}

func truncateTables(sqldb *sqlite3) {
	_, err := sqldb.Begin()
	if err != nil {
		fmt.Errorf(err)
	}
	defer sqldb.Close()
	
    
	stmt, err := sqldb.Prepare("
	DELETE FROM `songs`;
	")
	if err != nil {
		fmt.Errorf(err)
	} else {
		updateTimestamp(sqldb, "songs")
	}

	stmt2, err := sqldb.Prepare("
	DELETE FROM `albums`;
	")
	defer stmt2.Close()
	if err != nil {
		fmt.Errorf(err)
	} else {
	

	stmt3, err := sqldb.Prepare("
	DELETE FROM `covers`;
	")
	defer stmt3.Close()
	if err != nil {
		fmt.Errorf(err)
	} else {
		updateTimestamp(sqldb, "albums")
	} 
	}
    stmt4, err := sqldb.Prepare("
	DELETE FROM `artists`;
    ")
    if err != nil {
    	fmt.Errorf(err)
    } else {
    	updateTimestamp(sqldb, "artists")
    }
    defer stmt4.Close()
    sqldb.AutoCommit()
    
}


func usage(srgs string) {
	fmt.Errorf("Usage: %s action [args...]\n
	            %s scan file.db musicdir             :scan musicdir for new songs\n
	            %s fullscan file.db musicdir         : delete all records and do a full rescan\n
	            %s useradd file.db username password :add user\n
	            %s userdel file.db username          delete user\n",
	             args, args, args, args, args)
}

func main() {
	if len(os.Args) < 3 {
		usage(os.Args[0])
		exit(1)
	}

	sqldb *sqlite3
	
    action string
    action = os.Args[1]
    dbpath string
    dbpath = filepath.Dir(os.Args[2])


	sqldb, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		fmt.Errorf(err)
	}
	defer sqldb.Close()

	_, err := sqldb.Exec(init_sql)
	if err != nil {
		fmt.Errorf(err)
	}

    switch action {
    	case "scan":
    		musicdir string
    		musicdir = os.Args[3]
    		added int
    		added = 0
    		
    		fmt.Println("scanning " + musicdir + "...")

    		added = scanFs(sqldb, musicdir)
    		fmt.Println("added" + added + "files")

    		fmt.Println("cleaning up database...")
    		cleanupDb(sqldb)
    		
    	case "fullscan":
    		musicdir string
    		musicdir = os.Args[3]
    		added int
    		added = 0

    		fmt.Println("flushing database...")
    		truncateTables(sqldb)
    		
    		fmt.Println("scanning" + musicdir)
			added = scanFs(sqldb, musicdir)
			fmt.Println("added" + added + "files")
			 		
    	case "useradd":
    		user string
    		user = os.Args[3]
    		pass = os.Args[4]
    		stmt_ua, err := sqldb.Prepare("
			INSERT INTO `users` (`username`, `password`) VALUES (?, ?);
    		")
  	   		if err != nil {
   	    			fmt.Errorf(err)
  	   		}
    	    defer stmt_ua.Close()
    		    		
    		_, err := stmt_ua.Exec(1, user)
    		if err != nil {
    			fmt.Errorf(err)
    		}
    		_, err := stmt_ua.Exec(2, pass)
			if err != nil {
				fmt.Errorf(err)
			}
			
    	case "userdel":
    		user string
    		user = os.Args[3]
    		
    		stmt_ud, err := sqldb.Prepare("
			DELETE FROM `users` WHERE `username` = ?;
    		")
    		stmt_ud.Exec(1, user)
    		
    		if err != nil {
    			fmt.Errorf(err)
    		}
    		defer stmt_ud.Close()
    	
    	default:
    		usage(os.Args[0])
    		exit(1)
    }
    
	
}
