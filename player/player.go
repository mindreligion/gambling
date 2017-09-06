package player

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mindreligion/gambling/errors"
	"net/http"
)

func Fund(db *sql.DB, id int, amount int) error{
	if amount <= 0 || id <= 0 {
		return errors.New(http.StatusBadRequest, "Invalid parameters fund player")
	}
	res, err := db.Exec("INSERT INTO Player(id, amount) VALUES (?,?) ON DUPLICATE KEY UPDATE amount = amount + ?",
		id, amount, amount)
	if err != nil {
		return errors.New(http.StatusInternalServerError, err.Error())
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.New(http.StatusInternalServerError, err.Error())
	}
	if ra < 1 {
		return errors.New(http.StatusInternalServerError, "Update Failed")
	}
	return nil
}

func Take(db *sql.DB, id int, amount int) error {
	if amount <= 0 || id <= 0 {
		return errors.New(http.StatusBadRequest, "Invalid parameters take player")
	}
	res, err := db.Exec("UPDATE Player SET amount = amount - ? where id = ? and amount > ?",
		id, amount, id, amount)
	if err != nil {
		return errors.New(http.StatusInternalServerError, err.Error())
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.New(http.StatusInternalServerError, err.Error())
	}
	if ra < 1 {
		return errors.New(http.StatusBadRequest, "No player id or too low funds")
	}
	return nil
}

//func GetPlayers(db *sql.DB) []Player{
//	var players []Player
//	err := db.Ping()
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//	rows, err := db.Query("select id, amount from Player")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer rows.Close()
//	var (
//		id uint
//		amount uint
//	)
//	for rows.Next() {
//		err := rows.Scan(&id, &amount)
//		if err != nil {
//			log.Fatal(err)
//		}
//		players = append(players, Player{id, amount})
//	}
//	err = rows.Err()
//	if err != nil {
//		log.Fatal(err)
//	}
//	return players
//}
