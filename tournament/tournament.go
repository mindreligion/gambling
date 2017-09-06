package tournament

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mindreligion/gambling/errors"
	"net/http"
)

func Announce(db *sql.DB, id int, deposit int) error{
	if deposit <= 0 || id <= 0 {
		return errors.New(http.StatusBadRequest, "Invalid parameters announce tournament")
	}
	res, err := db.Exec("INSERT IGNORE INTO Tournament(id, deposit) VALUES(?, ?)",
		id, deposit)
	if err != nil {
		return errors.New(http.StatusInternalServerError, err.Error())
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.New(http.StatusInternalServerError, err.Error())
	}
	if ra < 1 {
		return errors.New(http.StatusInternalServerError, "Can swear - this tournament id exists")
	}
	return nil
}