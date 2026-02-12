package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func GetTask(id string) (*Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat 
				FROM scheduler 
			    WHERE id=:id`
	err := storage.db.QueryRow(query, sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// The AddTask function inserts a task into a scheduler table and returns the generated ID along with
// any errors encountered.
func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date,:title,:comment,:repeat)`
	res, err := storage.db.Exec(query, sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler 
			SET date = :date, 
                title = :title, 
                comment = :comment, 
                repeat = :repeat
            WHERE id = :id`
	result, err := storage.db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	return nil
}
func UpdateDate(next, id string) error {
	query := `UPDATE scheduler 
			SET date = :date
            WHERE id = :id`
	result, err := storage.db.Exec(query,
		sql.Named("date", next),
		sql.Named("id", id))

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %s not found", id)
	}

	return nil
}
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler 
				WHERE id=:id`
	result, err := storage.db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %s not found", id)
	}
	return nil

}

// This Go function retrieves tasks from a database based on a specified limit and returns them as a
// slice of Task pointers along with any potential errors.
func Tasks(limit int) ([]*Task, error) {
	var res []*Task

	query := `SELECT id, date, title, comment, repeat 
				FROM scheduler 
				ORDER BY date 
				LIMIT :limit;`
	rows, err := storage.db.Query(query, sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		res = append(res, &task)
	}

	if err := rows.Err(); err != nil {

		return nil, err
	}
	return res, nil
}
func TasksSearch(limit int, search string) ([]*Task, error) {
	var res []*Task

	query := `SELECT id, date, title, comment, repeat
				FROM scheduler WHERE title LIKE :search OR comment LIKE :search 
				ORDER BY date LIMIT :limit`
	rows, err := storage.db.Query(query, sql.Named("search", "%"+search+"%"),
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		res = append(res, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
func TasksDate(limit int, date string) ([]*Task, error) {
	var res []*Task

	query := `SELECT  id, date, title, comment, repeat
				FROM scheduler
				WHERE date = :date LIMIT :limit`
	rows, err := storage.db.Query(query, sql.Named("date", date),
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		res = append(res, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
