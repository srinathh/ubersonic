package main

import (
	"database/sql"
	_ "encoding/base64"
	"fmt"
	"github.com/dhowden/tag"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zhulik/go_mediainfo"
	_ "image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const initSQL = "CREATE TABLE `artists` (" +
	"`id`        INTEGER NOT NULL UNIQUE," +
	"`name`      TEXT," +
	"PRIMARY KEY(id)" +
	");" +
	"CREATE TABLE `albums` (" +
	"`id`        INTEGER NOT NULL UNIQUE," +
	"`title`     TEXT," +
	"`artistid`  INTEGER," +
	"`artist`    TEXT," +
	"PRIMARY KEY(id)" +
	");" +
	"CREATE INDEX `index_albums_artistid` ON `albums`(`artistid`); " +
	"CREATE TABLE `covers` (" +
	"`albumId`   INTEGER NOT NULL UNIQUE," +
	"`artistId`  INTEGER," +
	"`image`     BLOB," +
	"PRIMARY KEY(albumId)" +
	");" +
	"CREATE INDEX `index_covers_artistid` ON `covers`(`artistId`);" +
	"CREATE TABLE `songs` (" +
	"`id`        INTEGER NOT NULL UNIQUE," +
	"`title`     TEXT," +
	"`albumid`   INTEGER," +
	"`album`     TEXT," +
	"`artistid`  INTEGER," +
	"`artist`    TEXT," +
	"`trackn`    INTEGER," +
	"`discn`     INTEGER," +
	"`year`      INTEGER," +
	"`duration`  INTEGER," +
	"`bitRate`   INTEGER," +
	"`genre`     TEXT," +
	"`type`      TEXT," +
	"`filename`  TEXT," +
	"PRIMARY KEY(id)" +
	");" +
	"CREATE INDEX `index_songs_artistid_albumid` ON `songs`(`artistid`, `albumid`);" +
	"CREATE TABLE `users` (" +
	"`username`  TEXT NOT NULL UNIQUE," +
	"`password`  TEXT," +
	"PRIMARY KEY(username)" +
	");" +
	"CREATE TABLE `last_update_ts` (" +
	"`table_name`    TEXT NOT NULL UNIQUE," +
	"`mtime`         BIGINT," +
	"PRIMARY KEY(table_name)" +
	");"

/*
panic_if
base64Decode
getFileExtension

*/

func isExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func isDirectory(name string) (bool, error) {
	fInfo, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	return fInfo.IsDir(), nil
}

func dirs()

func searchCover(dir string) string {
	var cover string
	var coverNames []string
	coverNames = append(coverNames, "cover.jpg", "Cover.jpg", "cover.png", "Cover.png")

		for _, coverName := range coverNames {
			if isExists(coverName) == true {
				absPath, err := filepath.Abs(filepath.Join(dir + coverName))
				if err != nil {
					fmt.Errorf("%s", err)
					os.Exit(1)
				}
				cf, err := os.Open(absPath)
				if err != nil {
					fmt.Errorf("%s", err)
					os.Exit(1)
				}
				defer cf.Close()
				buf, err := ioutil.ReadAll(cf)
				if err != nil {
					fmt.Errorf("%s", err)
					os.Exit(1)
				}
				cover = string(buf)
				break
			}
		}
	return cover
}

func calcID(s string) uint64 {
	var hash uint64
	hash = 14695981039346656037
	for c := range s {
		hash = uint64(math.Pow(float64(hash), float64(c)))
		hash *= 1099511628211
	}
	hash &= 0x7FFFFFFFFFFFFFFF
	return hash
}

func updateTimestamp(sqldb *sql.DB, tableName string) {
	stmt, err := sqldb.Prepare(
		"INSERT OR REPLACE INTO `last_update_ts` (`table_name`, `mtime`)" +
			"VALUES (?, strftime('%000','now'));")
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	defer stmt.Close()
	_, err = stmt.Exec(tableName)
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
}

func insertArtist(sqldb *sql.DB, artist string) {

	stmt, err := sqldb.Prepare("INSERT OR REPLACE INTO `artists` (`id`, `name`) VALUES (?, ?);")
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	defer stmt.Close()

	_, err = stmt.Exec(calcID(artist), artist)
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	updateTimestamp(sqldb, "artists")
}

func insertAlbum(sqldb *sql.DB, album string, artist string, cover string) {
	var albumID uint64
	albumID = calcID(album + "@" + artist)
	var artistID uint64
	artistID = calcID(artist)

	stmtIs, err := sqldb.Prepare("INSERT OR IGNORE INTO `albums` (`id`) VALUES (?);")
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	defer stmtIs.Close()

	_, err = stmtIs.Exec(albumID)
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	stmtUp, err := sqldb.Prepare("UPDATE `albums` SET `title`=?, `artistid`=?, `artist`=? WHERE id=?;")
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}

	defer stmtUp.Close()

	_, err = stmtUp.Exec(album, artistID, artist, albumID)
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}

	if len(cover) > 0 {
		stmtCv, err := sqldb.Prepare("INSERT OR IGNORE INTO `covers` (`albumId`, `artistId`, `image`)VALUES (?, ?, ?)")
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}
		defer stmtCv.Close()

		_, err = stmtCv.Exec(albumID, artistID, cover)
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}
	}
	updateTimestamp(sqldb, "albums")
}

func insertSong(sqldb *sql.DB, fileName string, title string, artist string, album string, ext string, genre string, tn uint, year uint, discn uint, duration uint, bitrate uint) {
	stmt, err := sqldb.Prepare("INSERT OR REPLACE INTO `songs` " +
		"(`id`, `title`, `albumid`, `album`,`artistid`, `artist`, `type`, `genre`, `trackn`, `year`, `discn`, `duration`, `bitRate`, `filename`)" +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	defer stmt.Close()

	_, err = stmt.Exec(calcID(string(tn)+"@"+string(discn)+"@"+title+"@"+album+"@"+artist),
		title,
		calcID(album+"@"+artist),
		album,
		calcID(artist),
		artist,
		ext,
		genre,
		tn,
		year,
		discn,
		duration,
		bitrate,
		fileName)
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	updateTimestamp(sqldb, "songs")

}

func scanMusicFile(sqldb *sql.DB, fullPath string, dirCover string) bool {
	var artist string
	var albumartist string
	var album string
	var title string
	var genre string
	var track int
	var year int
	var discn int
	var bitrate uint
	var length uint

	f, err := os.Open(filepath.Dir(fullPath))
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		return false
	}
	var ext tag.FileType
	ext = m.FileType()

	if ext != tag.MP3 && ext != tag.OGG && ext != tag.FLAC && ext != tag.M4A && ext != tag.M4B && ext != tag.M4P && ext != tag.ALAC {
		return false
	}

	var coverImage *tag.Picture
	var cover string
	coverImage = m.Picture()
	if ext == tag.UnknownFileType {
		fmt.Errorf("This file can't index with me.")
		os.Exit(1)
	}

	cover = string(coverImage.Data)
	if len(cover) == 0 && len(dirCover) > 0 && len(dirCover) < 400*1024 {
		cover = dirCover
	}

	albumartist = strings.Trim(m.AlbumArtist(), " ")
	if len(albumartist) > 0 {
		artist = albumartist
	} else {
		artist = strings.Trim(m.Artist(), " ")
	}
	album = strings.Trim(m.Album(), " ")
	title = strings.Trim(m.Title(), " ")
	discn, _ = m.Disc()

	genre = strings.Trim(m.Genre(), " ")
	year = m.Year()
	track, _ = m.Track()

	mediabytes, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	mi := mediainfo.NewMediaInfo()
	err = mi.OpenMemory(mediabytes)

	length = uint(mi.Duration())
	i, _ := strconv.Atoi(mi.Get("BitRate"))
	bitrate = uint(i)
	tx, err := sqldb.Begin()
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	defer sqldb.Close()

	insertSong(sqldb, fullPath, title, artist, album, string(ext), genre, uint(track), uint(year), uint(discn), length, bitrate)

	insertAlbum(sqldb, album, artist, cover)

	insertArtist(sqldb, artist)

	fmt.Println("added")
	tx.Commit()
	return true
}

func scanFs(sqldb *sql.DB, name string) int {
	var addedSongs int
	addedSongs = 0
	var coverData string

	var isDir, _ = isDirectory(name)
	if isDir != true {
		return 0
	}

	fileInfos, err := ioutil.ReadDir(name)
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	for _, fileInfo := range fileInfos {
		err = filepath.Walk(name, walkFunk)
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		} 
	}	
}

func cleanupDB(sqldb *sql.DB) {
	stmtAl, err := sqldb.Prepare("DELETE FROM `albums` WHERE ( SELECT count(id) FROM `songs` WHERE albumId = albums.id) = 0;")
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	defer stmtAl.Close()

	stmtAr, err := sqldb.Prepare("DELETE FROM `artists` WHERE ( SELECT count(id) FROM `albums` WHERE artistId = artists.id) = 0;")
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	defer stmtAr.Close()

	stmtCv, err := sqldb.Prepare("DELETE FROM `covers` WHERE ( SELECT count(id) FROM `albums` WHERE albumId = covers.albumId) = 0;")
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	defer stmtCv.Close()

}

func truncateTables(sqldb *sql.DB) {
	tx, err := sqldb.Begin()
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	defer sqldb.Close()

	stmtSn, err := sqldb.Prepare("DELETE FROM `songs`;")
	defer stmtSn.Close()
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	} else {
		updateTimestamp(sqldb, "songs")
	}

	stmtAl, err := sqldb.Prepare("DELETE FROM `albums`;")
	defer stmtAl.Close()
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	} else {
		stmtCv, err := sqldb.Prepare("DELETE FROM `covers`;")
		defer stmtCv.Close()
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		} else {
			updateTimestamp(sqldb, "albums")
		}
	}
	stmtAr, err := sqldb.Prepare("DELETE FROM `artists`;")
	defer stmtAr.Close()
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	} else {
		updateTimestamp(sqldb, "artists")
	}
	tx.Commit()
}

func usage(args string) {
	fmt.Printf("Usage: %s action [args...]\n"+
		"%s scan file.db musicdir: scan musicdir for new songs\n"+
		"%s fullscan file.db musicdir: delete all records and do a full rescan\n"+
		"%s useradd file.db username password: add user\n"+
		"%s userdel file.db username: delete user\n",
		args, args, args, args, args)
}

func main() {
	if len(os.Args) < 3 {
		usage(os.Args[0])
		os.Exit(0)
	}

	var sqldb *sql.DB

	var action string
	action = os.Args[1]
	var dbpath string
	fmt.Println(dbpath)
	dbpath, err := filepath.Abs(os.Args[2])
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	sqldb, err = sql.Open("sqlite3", dbpath)
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	defer sqldb.Close()
	if action != "scan" && action != "fullscan" && action != "useradd" && action != "userdel" {
		fmt.Println("Invalid action")
		os.Exit(1)
	}

	if isExists(dbpath) != true {
		_, err = sqldb.Exec(initSQL)
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}
	}
	switch action {
	case "scan":
		var musicDir string
		musicDir = os.Args[3]
		var added int
		added = 0

		fmt.Println("scanning " + musicDir + "...")

		added = scanFs(sqldb, musicDir)
		fmt.Println("added" + strconv.Itoa(added) + "files")

		fmt.Println("cleaning up database...")
		cleanupDB(sqldb)

	case "fullscan":
		var musicDir string
		musicDir = os.Args[3]
		var added int
		added = 0

		fmt.Println("flushing database...")
		truncateTables(sqldb)

		fmt.Println("scanning" + musicDir)
		added = scanFs(sqldb, musicDir)
		fmt.Println("added" + strconv.Itoa(added) + "files")

	case "useradd":
		var user string
		var pass string
		user = os.Args[3]
		pass = os.Args[4]
		stmtUa, err := sqldb.Prepare("INSERT INTO `users` (`username`, `password`) VALUES (?, ?);")
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}
		defer stmtUa.Close()

		_, err = stmtUa.Exec(user, pass)
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}
		fmt.Println("Succeeded to add user: " + user)

	case "userdel":
		var user string
		user = os.Args[3]

		stmtUd, err := sqldb.Prepare("DELETE FROM `users` WHERE `username` = ?;")
		stmtUd.Exec(user)

		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}
		defer stmtUd.Close()
		fmt.Println("Succeeded to delete user: " + user)

	default:
		usage(os.Args[0])
		os.Exit(0)
	}

}
