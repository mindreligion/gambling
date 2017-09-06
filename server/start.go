package main

import (
	"regexp"
	"net/http"
	"net/url"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"github.com/mindreligion/gambling/player"
	"github.com/mindreligion/gambling/tournament"
	"github.com/mindreligion/gambling/errors"
)

var validPath = regexp.MustCompile("^/(fund|take|announceTournament|joinTournament|resultTournament|balance|reset)$")
var positiveInt = regexp.MustCompile("^([1-9]\\d*)$")


func checkHTTPMethod(m string, r *http.Request) error{
	if r.Method != m {
		return errors.New(http.StatusNotFound, "Invalid HTTP method")
	}
	return nil
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, db *sql.DB) error, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(url.ParseQuery(r.URL.RawQuery))
		err := fn(w, r, db)
		if err != nil {
			gamblingError, ok := err.(errors.Code)
			if ok {
				w.WriteHeader(gamblingError.Code())
			} else {
				log.Println("Unexpected error type")
				w.WriteHeader(http.StatusInternalServerError)
			}
			log.Println(err.Error())
		}
	}
}

func takeHandler(w http.ResponseWriter, r *http.Request,db *sql.DB) error{
	if err := checkHTTPMethod("GET", r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	log.Println(r.Form)
	playerID, _ := strconv.Atoi(r.Form.Get("playerId"))
	amount, _ := strconv.Atoi(r.Form.Get("points"))
	if err := player.Take(db, playerID, amount); err != nil {
		return err
	}
	return nil
}

func fundHandler(w http.ResponseWriter, r *http.Request,db *sql.DB) error{
	if err := checkHTTPMethod("GET", r); err != nil {
		return err
	}
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return errors.New(http.StatusBadRequest, "Invalid query")
	}
	if len(params) != 2 {
		return errors.New(http.StatusBadRequest, "Invalid query")
	}
	playerIDs, ok := params["playerId"]
	if !ok {
		return errors.New(http.StatusBadRequest, "Player id is missed")
	}
	if len(playerIDs) != 1 {
		return errors.New(http.StatusBadRequest, "Player id duplicate")
	}
	if !positiveInt.MatchString(playerIDs[0]){
		return errors.New(http.StatusBadRequest, "Invalid Player id")
	}
	playerID, err := strconv.Atoi(playerIDs[0])
	if err != nil {
		return errors.New(http.StatusBadRequest, "Invalid Player id atoi")
	}
	if playerID <= 0 {
		return errors.New(http.StatusBadRequest, "Player id must be positive int")
	}
	pointAmounts, ok := params["points"]
	if !ok {
		return errors.New(http.StatusBadRequest, "Points amount is missed")
	}
	if len(pointAmounts) != 1 {
		return errors.New(http.StatusBadRequest, "Points amount duplicate")
	}
	if !positiveInt.MatchString(pointAmounts[0]){
		return errors.New(http.StatusBadRequest, "Invalid Points amount")
	}
	pointAmount, err := strconv.Atoi(pointAmounts[0])
	if err != nil {
		return errors.New(http.StatusBadRequest, "Invalid Points amount atoi")
	}
	if pointAmount <= 0 {
		return errors.New(http.StatusBadRequest, "Points amount must be positive int")
	}
	if err := player.Fund(db, playerID, pointAmount); err != nil {
		return err
	}
	return nil
}

func announceTournamentHandler(w http.ResponseWriter, r *http.Request,db *sql.DB) error {
	if err := checkHTTPMethod("GET", r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	log.Println(r.Form)
	tournamentID, _ := strconv.Atoi(r.Form.Get("tournamentId"))
	deposit, _ := strconv.Atoi(r.Form.Get("deposit"))
	if err := tournament.Announce(db, tournamentID, deposit); err != nil {
		return err
	}
	return nil
}

func joinTournamentHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if checkHTTPMethod("GET", r) != nil {
		return
	}

}

func resultTournamentHandler(w http.ResponseWriter, r *http.Request,db *sql.DB) {
	if checkHTTPMethod("POST", r) != nil {
		return
	}

}

func balanceHandler(w http.ResponseWriter, r *http.Request,db *sql.DB) {
	if checkHTTPMethod("GET", r) != nil {
		return
	}

}

func resetHandler(w http.ResponseWriter, r *http.Request,db *sql.DB) {
	if checkHTTPMethod("GET", r) != nil {
		return
	}
}

func errorHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/gambling_test")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	http.HandleFunc("/take", makeHandler(takeHandler, db))
	http.HandleFunc("/fund", makeHandler(fundHandler, db))
	http.HandleFunc("/announceTournament", makeHandler(announceTournamentHandler, db))
	//http.HandleFunc("/joinTournament", makeHandler(joinTournamentHandler, db))
	//http.HandleFunc("/resultTournament", makeHandler(resultTournamentHandler, db))
	//http.HandleFunc("/balance", makeHandler(balanceHandler, db))
	//http.HandleFunc("/reset", makeHandler(resetHandler, db))
	http.HandleFunc("/", errorHandler)
	http.ListenAndServe(":8080", nil)
}
