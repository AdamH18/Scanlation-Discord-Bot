package database

import (
	"database/sql"
	"log"
)

/*
These are wrappers around the db.Exec function for the SQLite repository. Every table should have
one, and an Exec function should never be called directly in the dbrepo.go file. These are used
to check what information was changed in the database and let things that use that information
know that they should update (mostly billboards). Updates are sent using preregistered channels
*/

func (r *SQLiteRepository) SerielExec(query string, args ...any) (sql.Result, error) {
	//Send execution over to single threaded execution handler
	results := make(chan ExecOut)
	DBWriterCh <- ExecIn{
		quer: query,
		vals: args,
		ch:   results,
	}
	out := <-results
	res := out.res
	err := out.err
	if err == nil {
		log.Printf("Query %s with args %v succeeded\n", query, args)
		return res, nil
	}

	log.Printf("Error on query %s with args %v - %s\n", query, args, err.Error())
	return nil, err
}

func (r *SQLiteRepository) RemindersExec(query string, args ...any) (sql.Result, error) {
	return r.SerielExec(query, args...)
}

// Changing the series table can affect series billboards (repo_link or name_full)
func (r *SQLiteRepository) SeriesExec(guild string, series string, query string, args ...any) (sql.Result, error) {
	res, err := r.SerielExec(query, args...)
	//If query errored, return error
	if err != nil {
		return res, err
	}
	num, err := res.RowsAffected()
	if err != nil {
		//No clue why this would error out, but do nothing out of caution if it does
		log.Println("Error getting number of affected rows: " + err.Error())
		return res, nil
	}
	if num == 0 {
		//If nothing changed, just leave
		return res, nil
	}
	//Otherwise update billboard based on passed values
	SeriesCh <- func() (string, string) { return guild, series }
	return res, nil
}

func (r *SQLiteRepository) ChannelsExec(query string, args ...any) (sql.Result, error) {
	return r.SerielExec(query, args...)
}

// Changing the users table can affect colors and series billboards (color)
func (r *SQLiteRepository) UsersExec(guild string, query string, args ...any) (sql.Result, error) {
	res, err := r.SerielExec(query, args...)
	if err != nil {
		return res, err
	}
	//Bit complex to find affected user and all relevant series to update, so just update everything
	SeriesCh <- func() (string, string) { return guild, "" }
	ColorsCh <- guild
	return res, err
}

// Changing the assignments table can affect assignments and series billboards
func (r *SQLiteRepository) SeriesAssignmentsExec(guild string, series string, query string, args ...any) (sql.Result, error) {
	res, err := r.SerielExec(query, args...)
	if err != nil {
		return res, err
	}

	num, err := res.RowsAffected()
	if err != nil {
		//No clue why this would error out, but just update everything
		log.Println("Error getting number of affected rows: " + err.Error())
		AssignmentsCh <- guild
		SeriesCh <- func() (string, string) { return guild, "" }
		return res, nil
	}
	if num == 0 {
		//If nothing changed, just leave
		return res, nil
	}

	//If something changed, update billboards
	AssignmentsCh <- guild
	SeriesCh <- func() (string, string) { return guild, series }
	return res, nil
}

// Changing the series notes table can affect series billboards
func (r *SQLiteRepository) SeriesNotesExec(guild string, series string, query string, args ...any) (sql.Result, error) {
	res, err := r.SerielExec(query, args...)
	if err != nil {
		return res, err
	}
	SeriesCh <- func() (string, string) { return guild, series }
	return res, err
}

func (r *SQLiteRepository) JobsExec(query string, args ...any) (sql.Result, error) {
	return r.SerielExec(query, args...)
}

func (r *SQLiteRepository) MemberRoleExec(query string, args ...any) (sql.Result, error) {
	return r.SerielExec(query, args...)
}

func (r *SQLiteRepository) SeriesChannelsExec(query string, args ...any) (sql.Result, error) {
	return r.SerielExec(query, args...)
}

func (r *SQLiteRepository) SeriesBillboardsExec(query string, args ...any) (sql.Result, error) {
	return r.SerielExec(query, args...)
}

func (r *SQLiteRepository) RolesBillboardsExec(query string, args ...any) (sql.Result, error) {
	return r.SerielExec(query, args...)
}

func (r *SQLiteRepository) ColorsBillboardsExec(query string, args ...any) (sql.Result, error) {
	return r.SerielExec(query, args...)
}

func (r *SQLiteRepository) NotificationChannelsExec(query string, args ...any) (sql.Result, error) {
	return r.SerielExec(query, args...)
}
