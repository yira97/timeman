package store

import (
	"context"
	"database/sql"
	"fmt"
	bosql "github.com/yrfg/boast/sql"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/yrfg/boast/validate"
)

// TaskStat describe different status of task in db
type TaskStat int

const (
	// Pending task is just add, to do.
	Pending TaskStat = iota
	// Withdraw task is delete from the pending list, should not shows on timeline
	Withdraw
	// Doing task is the task on going
	Doing
	// Cancel task is abandon from doing
	Cancel
	// Timeout task is refuse by setting
	Timeout
	// Done task is finished on schedule
	Done
)

// TableSetup is to init table if not exist
func TableConstruct(ctx context.Context, conn *pgx.Conn) (err error) {
	sql := `
	CREATE TABLE IF NOT EXISTS timeman_task (
	  id           BIGSERIAL                          NOT NULL                         ,
	  name         VARCHAR(63)                        NOT NULL                         ,
	  note         TEXT                                                                ,
	  -- weight > 0 order by created_at DESC
	  -- (1, 1001) > (1, 1000) > (-1, 999) > (-1, 1002)
	  -- weight < 0 order by created_at ASC
	  weight       SMALLINT                           NOT NULL     DEFAULT -1          ,
	  avatar       VARCHAR(255)                       NOT NULL     DEFAULT ''          ,
	  -- stat: 0: PENDING | 1: WITHDRAW |  2: DOING | 3: CANCEL | 4: TIMEOUT | 5: DONE
	  stat         SMALLINT                           NOT NULL     DEFAULT 0           ,
	  created_at   TIMESTAMP WITHOUT TIME ZONE        NOT NULL     DEFAULT NOW()       ,
	  updated_at   TIMESTAMP WITHOUT TIME ZONE        NOT NULL     DEFAULT NOW()       ,
	  withdraw_at  TIMESTAMP WITHOUT TIME ZONE                                         ,
	  doing_at     TIMESTAMP WITHOUT TIME ZONE                                         ,
	  cancel_at    TIMESTAMP WITHOUT TIME ZONE                                         ,
	  dead_time    TIMESTAMP WITHOUT TIME ZONE                                         ,
	  done_at      TIMESTAMP WITHOUT TIME ZONE                                         ,
	  map_id       BIGINT                             NOT NULL                         ,
	  PRIMARY KEY(id)
	);
	
	CREATE INDEX IF NOT EXISTS task_list_rank_up          ON    timeman_task       (weight DESC, created_at DESC      );
	CREATE INDEX IF NOT EXISTS task_list_rank_bottom      ON    timeman_task       (weight DESC, created_at ASC       );
	
	CREATE TABLE IF NOT EXISTS timeman_task_map (
	  id           BIGSERIAL                          NOT NULL                         ,
	  name         VARCHAR(63)                        NOT NULL                         ,
	  PRIMARY KEY(id)
	);
`
	_, err = conn.Exec(ctx, sql)
	return
}

// TableDestroy is a func to drop every table that service in used
func TableDestroy(ctx context.Context, conn *pgx.Conn) (err error) {
	sql := `
		drop table if exists timeman_task;

		drop table if exists timeman_task_map;
	`
	_, err = conn.Exec(ctx, sql)
	return
}

func (td TaskStat) String() string {
	return [...]string{"pending", "withdraw", "doing", "cancel", "timeout", "done"}[td]
}

var (
	// FilterToDo is a view of unpland task
	FilterToDo = []TaskStat{Pending}
	// FilterHistory is a view shows main stat
	FilterHistory = []TaskStat{Done, Timeout, Cancel, Doing}
	// FilterFinished is a view show only happend on main stat
	FilterFinished = []TaskStat{Done, Doing}
	// FilterAll is a view of all task
	FilterAll = []TaskStat{Done, Timeout, Cancel, Doing, Withdraw, Pending}
)

const (
	// TableNameTimeManTask is the schema name of task table
	TableNameTimeManTask string = `timeman_task`
	// TableNameTimeManTaskMap is the schema name of map table
	TableNameTimeManTaskMap string = `timeman_task_map`
)

var (
	// TaskNameChecker is a validator for task name
	TaskNameChecker validate.Checker = &validate.StringCheck{
		MaxLen: 63,
		MinLen: 1,
		Mode:   validate.MaxLen | validate.MinLen,
	}
	// TaskMapNameChecker is a validator for map name
	TaskMapNameChecker validate.Checker = &validate.StringCheck{
		MaxLen: 63,
		MinLen: 1,
		Mode:   validate.MaxLen | validate.MinLen,
	}
)

// TimeManTaskMapDisplay is for output
type TimeManTaskMapDisplay struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TimeManTaskDisplay is for output
type TimeManTaskDisplay struct {
	ID     int64  `json:"id"`
	Stat   string `json:"stat"`
	Avatar string `json:"avatar"`
	// CreatedAt unit is ms
	CreatedAt int64 `json:"created_at"`
	// ChangeAt unit is ms
	StableAt int64  `json:"change_at"`
	Note     string `json:"note"`
	Name     string `json:"name"`
}

// CreateTimeManTaskMap is a func to create new task map
func CreateTimeManTaskMap(ctx context.Context, conn *pgx.Conn, name string) (taskMapID int64, err error) {
	if err = TaskMapNameChecker.Check(name); err != nil {
		return
	}
	err = conn.QueryRow(ctx, "INSERT INTO "+TableNameTimeManTaskMap+" (name) VALUES ($1) RETURNING id", name).Scan(&taskMapID)
	return taskMapID, err
}

// UpdateTimeManTaskMapName is a func to change the task map name
func UpdateTimeManTaskMapName(ctx context.Context, conn *pgx.Conn, taskMapID int64, newName string) (err error) {
	if err = TaskMapNameChecker.Check(newName); err != nil {
		return
	}
	_, err = conn.Exec(ctx, "UPDATE "+TableNameTimeManTaskMap+" SET name = $1 WHERE id = $2", newName, taskMapID)
	return err
}

// ListTimeManTaskMap is a func to shows the task map
func ListTimeManTaskMap(ctx context.Context, conn *pgx.Conn, offset int, count int) (mapList []TimeManTaskMapDisplay) {
	mapList = make([]TimeManTaskMapDisplay, 0)
	rows, err := conn.Query(ctx, "SELECT id, name FROM "+TableNameTimeManTaskMap+" LIMIT $1 OFFSET $2", count, offset)
	if err != nil {
		fmt.Println("error ListTimeManTaskMap query failed")
		return mapList
	}
	defer rows.Close()
	for rows.Next() {
		dis := TimeManTaskMapDisplay{}
		err := rows.Scan(&dis.ID, &dis.Name)
		if err != nil {
			fmt.Println(err)
		}
		mapList = append(mapList, dis)
	}
	return mapList
}

// CreateTask is a func to create new task under a task map
func CreateTask(ctx context.Context, conn *pgx.Conn, taskMapID int64, taskName string) (newTaskID int64, err error) {
	if err = TaskNameChecker.Check(taskName); err != nil {
		return
	}
	err = conn.QueryRow(ctx, "INSERT INTO "+TableNameTimeManTask+" (map_id, name) VALUES ($1, $2) RETURNING id", taskMapID, taskName).Scan(&newTaskID)
	return newTaskID, err
}

// UpdateTaskName is a func to change the task name of a task
func UpdateTaskName(ctx context.Context, conn *pgx.Conn, taskID int64, newTaskName string) (err error) {
	if err = TaskNameChecker.Check(newTaskName); err != nil {
		return
	}
	_, err = conn.Exec(ctx, "UPDATE "+TableNameTimeManTask+" SET name = $1 WHERE id = $2", newTaskName, taskID)
	return err
}

func ListTimeManTaskByFilterTag(ctx context.Context, conn *pgx.Conn, timeMapID int64, filterTag string, offset int, count int) (taskList []TimeManTaskDisplay) {
	var filter = FilterAll
	switch filterTag {
	case "todo":
		filter = FilterToDo
	case "history":
		filter = FilterFinished
	case "finished":
		filter = FilterFinished
	case "all":
		// default
	}
	return ListTimeManTask(ctx,conn,timeMapID,filter,offset, count)
}

// ListTimeManTask is a func to shows task
func ListTimeManTask(ctx context.Context, conn *pgx.Conn, timeMapID int64, filter []TaskStat, offset int, count int) (taskList []TimeManTaskDisplay) {
	taskList = make([]TimeManTaskDisplay, 0)
	elected := make([]int64,len(filter))
	for _, st := range filter {
		elected = append(elected, int64(st))
	}
	inCons := &bosql.InConstraint{
		ColumnName: "stat",
		Elected:    elected,
		Mode:       bosql.InAsEqual | bosql.SOFT,
		Dialect:    bosql.PG,
	}
	rows, err := conn.Query(ctx, " SELECT id, stat, avatar, created_at, withdraw_at, doing_at, cancel_at, dead_time, done_at, note, name "+
		" FROM "+TableNameTimeManTask+
		"WHERE map_id = $1 "+
		inCons.String()+
		" LIMIT $2 OFFSET $3 ",
		timeMapID,
		count,
		offset,
	)
	if err != nil {
		fmt.Println("error ListTimeManTask query failed")
		return taskList
	}
	defer rows.Close()
	for rows.Next() {
		task := TimeManTaskDisplay{}
		var rStat int64
		var rCreatedAt time.Time
		var rWithDrawAt sql.NullTime
		var rDoingAt sql.NullTime
		var rCancelAt sql.NullTime
		var rDeadTime sql.NullTime
		var rDoneAt sql.NullTime
		var rNote sql.NullString
		if err = rows.Scan(&task.ID, &rStat, &task.Avatar, &rCreatedAt, &rWithDrawAt, &rDoingAt, &rCancelAt, &rDeadTime, &rDoneAt, &task.Name); err != nil {
			fmt.Println(err)
		}
		task.CreatedAt = rCreatedAt.Unix() * 1000
		switch TaskStat(rStat) {
		case Pending:
			task.StableAt = task.CreatedAt
			task.Stat = Pending.String()
		case Withdraw:
			task.Stat = Withdraw.String()
			if rWithDrawAt.Valid {
				task.StableAt = rWithDrawAt.Time.Unix() * 1000
			}
		case Doing:
			if rDoingAt.Valid {
				if rDeadTime.Valid && rDeadTime.Time.Before(time.Now()) {
					task.Stat = Timeout.String()
					task.StableAt = rDeadTime.Time.Unix() * 1000
					break
				}
				task.Stat = Doing.String()
				task.StableAt = rDoingAt.Time.Unix() * 1000
			}
		case Cancel:
			task.Stat = Cancel.String()
			if rCancelAt.Valid {
				task.StableAt = rCancelAt.Time.Unix() * 1000
			}
		case Done:
			task.Stat = Done.String()
			if rDoneAt.Valid {
				task.StableAt = rDoneAt.Time.Unix() * 1000
			}
		}
		if rNote.Valid {
			task.Note = rNote.String
		}
		taskList = append(taskList, task)
	}
	return taskList
}
