package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/yrfg/timeman/pkg/tmerr"
)

const (
	TableNameTimeManToDo    string = "timeman_todo"
	TableNameTimeManMap     string = "timeman_map"
	TableTimeManToDoNameLen int    = 63
	TableTimeManMapNameLen  int    = 63
)

type TimeManMapDisplay struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TimeManTodoDisplay struct {
	ID        int64  `json:"id"`
	Stat      string `json:"stat"`
	Avatar    string `json:"avatar"`
	CreatedAt int64  `json:"created_at"`
	StableAt  int64  `json:"change_at"`
	Note      string `json:"note"`
	Name      string `json:"name"`
}

func CreateTimeManMap(conn *pgx.Conn, name string) (timeManMapID int64, err error) {
	if name == "" {
		return timeManMapID, tmerr.Btw(tmerr.NullStringParamsError, "name is empty")
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO "+TableNameTimeManMap+" (name) VALUES ($1) RETURNING id", name).Scan(&timeManMapID)
	return timeManMapID, err
}

func UpdateTimeManMapName(conn *pgx.Conn, timelineID int64, newName string) (err error) {
	if newName == "" {
		return tmerr.Btw(tmerr.NullStringParamsError, "newName is empty")
	} else if len(newName) > TableTimeManToDoNameLen {
		return tmerr.Btw(tmerr.TooLongStringParamsError, "newName is too long")
	}
	_, err = conn.Exec(context.Background(), "UPDATE "+TableNameTimeManMap+" SET name = $1 WHERE id = $2", newName, timelineID)
	return err
}

func ListTimeManMap(conn *pgx.Conn, offset int, count int) (mapList []TimeManMapDisplay) {
	mapList = make([]TimeManMapDisplay, 0)
	rows, err := conn.Query(context.Background(), "SELECT (id, name) FROM "+TableNameTimeManMap+" LIMIT $1 OFFSET $2", count, offset)
	if err != nil {
		fmt.Println("error list timemanmap query failed")
		return mapList
	}
	defer rows.Close()
	for rows.Next() {
		dis := TimeManMapDisplay{}
		rows.Scan(&dis.ID, &dis.Name)
		mapList = append(mapList, dis)
	}
	return mapList
}

func CreateToDo(conn *pgx.Conn, timeMapID int64, todoName string) (newToDoID int64, err error) {
	if todoName == "" {
		return newToDoID, tmerr.Btw(tmerr.NullStringParamsError, "todoName is empty")
	} else if len(todoName) > TableTimeManToDoNameLen {
		return newToDoID, tmerr.Btw(tmerr.TooLongStringParamsError, "todoName is too long")
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO "+TableNameTimeManToDo+" (map_id, name) VALUES ($1, $2) RETURNING id", timeMapID, todoName).Scan(&newToDoID)
	return newToDoID, err
}

func ListTimeManToDo(conn *pgx.Conn, timeMapID int64, offset int, limit int) (todoList []TimeManTodoDisplay) {
	todoList = make([]TimeManTodoDisplay, 0)
	return todoList
}
