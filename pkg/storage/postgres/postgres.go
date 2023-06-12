package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Tasks возвращает список задач из БД.
func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}

// ChangeTask изменяет задачу по id
func (s *Storage) ChangeTask(id int, title string, content string) error {
	_, err := s.db.Exec(context.Background(), `
			UPDATE tasks SET title = $2, content = $3 WHERE id = $1;
		`, id,
		title,
		content,
	)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// DeleteTask удаление задачи по id
func (s *Storage) DeleteTask(id int) error {
	_, err := s.db.Exec(context.Background(), `
			DELETE FROM tasks WHERE id = $1;
	`,
		id,
	)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// NewListTasks создаёт список задач
func (s *Storage) NewListTasks(tasks []Task) error {
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.TODO())
	// пакетный запрос
	batch := new(pgx.Batch)
	for _, t := range tasks {
		batch.Queue(`INSERT INTO tasks(title, content) VALUES ($1, $2)`,
			t.Title, t.Content)
	}
	// отправка пакета в БД
	res := tx.SendBatch(context.TODO(), batch)
	err = res.Close()
	if err != nil {
		return err
	}
	err = tx.Commit(context.TODO())

	return err
}
