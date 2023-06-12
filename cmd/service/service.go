package main

import (
	"db/apps/pkg/storage/postgres"
	"fmt"
	"log"
)

// var db storage.Interface
type Tasks []postgres.Task

func New() *Tasks {
	var arr Tasks
	return &arr
}

func (t *Tasks) Add(title, content string) {
	tnew := postgres.Task{
		Title:   title,
		Content: content,
	}
	*t = append(*t, tnew)
}
func main() {
	var err error

	connstr := "postgres://postgres:Ptds_17031993@localhost:5432/dataBase"
	db, err := postgres.New(connstr)
	if err != nil {
		log.Fatal(err)
	}
	task := postgres.Task{
		Title:   "New task",
		Content: "New Content",
	}
	t := New()
	for i := 0; i < 10; i++ {
		t.Add("TitleList", "ContentList")
	}
	id, err := db.NewTask(task)
	if err != nil {
		log.Fatal(err)
	}
	err = db.ChangeTask(id, "ChangeTrue Title", "Change Content")
	if err != nil {
		log.Fatal(err)
	}
	err = db.NewListTasks(*t)
	err = db.DeleteTask(3)
	viewTasks, err := db.Tasks(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(viewTasks)
}
