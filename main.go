package main

import (
	"fmt"
	"os"
	"log"

	"github.com/joho/godotenv"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"

	"strconv"
)

type VideoGame struct {
	id int64 `json:"id"`
	Name string `json:"name"`
	Genre string `json:"genre"`
	Year int64 `json:"year"`
}

var _ = godotenv.Load(".env")


var (
	connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("user"),
		os.Getenv("pass"),
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("db_name"),
		)
)

func getDb()(*sql.DB, error){
	return sql.Open("mysql", connectionString)
}

func stringToInt64(s string)(int64, error){
	num, err := strconv.ParseInt(s, 0, 64)
	if err != nil{
		return 0, err
	}
	return num, err
}

//CURD operation perform here 

func createVideoGame(db *sql.DB, videoGame VideoGame)error{
	_, err := db.Exec("insert into video_games (name, genre, year) values (?, ?, ?)", videoGame.Name, videoGame.Genre, videoGame.Year)
	return err
}

func deleteVideoGame(db *sql.DB, id int64)error{
	_, err := db.Exec("delete from vedio_games where id = ?", id)
	return err
}

func updateVideoGame(db *sql.DB, videoGame VideoGame)error{
	_, err := db.Exec("update vedio_games set name = ?, genre = ?, year = ? where id = ?", videoGame.Name, videoGame.Genre, videoGame.Year, videoGame.id)
	return err
}

func getVideoGamebyId(db *sql.DB, id int64)(*VideoGame, error){
	var videoGame VideoGame
	row := db.QueryRow("select id, name, genre, year from vedio_games where id = ?", id)
	err := row.Scan(&videoGame.id, &videoGame.Name, &videoGame.Genre, &videoGame.Year)
	return &videoGame, err
}

func getVideoGame(db *sql.DB)([]VideoGame, error){
	videoGames := []VideoGame{}

	rows, err := db.Query("select * from video_games")
	if err != nil{
		return videoGames, err
	}
	for rows.Next(){
		var videoGame VideoGame
		err := rows.Scan(&videoGame.id, &videoGame.Name, &videoGame.Genre, &videoGame.Year)
		if err != nil{
			return videoGames, err
		}
		videoGames = append(videoGames, videoGame)
	}
	return videoGames, nil
}
// CURD Function end here



func main(){
	db, err := getDb()
	if err != nil {
		log.Fatal("Cannot connect to db")
	}else{
		if db.Ping() != nil{
			log.Fatal("cannot ping to db")
		}
	}
	fmt.Println("Starting Game API service...")

	fmt.Println(`
  ____  ___       _    ____ ___ 
 / ___|/ _ \     / \  |  _ \_ _|
| |  _| | | |   / _ \ | |_) | | 
| |_| | |_| |  / ___ \|  __/| | 
 \____|\___/  /_/   \_\_|  |___|
	`)
	// router and making router function here 
	router := mux.NewRouter()
	HandlingRoute(router, db)

	http.ListenAndServe(":8080", router)
}

func HandlingRoute(router *mux.Router, db *sql.DB){
	router.HandleFunc("/videogames", func(w http.ResponseWriter, r *http.Request){
		getVideoGame(db)
		}).Methods("GET")
	router.HandleFunc("/videogames/{id}", func(w http.ResponseWriter, r *http.Request){
		idstr := mux.Vars(r)["id"]
		id, err := stringToInt64(idstr)
		if err != nil{
			http.Error(w, "Invalid ID: must enter the number", http.StatusBadRequest)
		}
		vedioGame, err := getVideoGamebyId(db, id)
		json.NewEncoder(w).Encode(vedioGame)
	}).Methods("GET")
	router.HandleFunc("/videogames/{id}", func(w http.ResponseWriter, r *http.Request){
		idstr := mux.Vars(r)["id"]
		id, err := stringToInt64(idstr)
		if err != nil{
			http.Error(w, "Invalid ID: must enter the number", http.StatusBadRequest)
		}
		deleteVideoGame(db, id)
		}).Methods("DELETE")
	router.HandleFunc("/videogames/")
}
