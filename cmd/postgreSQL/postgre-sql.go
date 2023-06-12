package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// структуры обычно используются для
// описания модели данных сущностей,
// хранящихся в БД
type user struct {
	id   int
	name string
}

func main() {
	// пустой контекст
	var ctx context.Context = context.Background()

	//pwd := os.Getenv("Ptds_1703")
	// Подключение к БД. Функция возвращает объект БД.
	db, err := pgxpool.Connect(ctx, "postgres://postgres:Ptds_17031993@localhost:5432/dataBase")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	// Не забываем очищать ресурсы.
	defer db.Close()
	// Проверка соединения с БД. На случай, если sql.Open этого не делает.
	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	data := []user{
		{name: "Rob Pike"},
		{name: "Ken Thompson"},
		{name: "Robert Grismer"},
	}
	addUsersTx(ctx, db, data)
	users, _ := users(ctx, db)

	fmt.Println(users)
}

// addUsersTx добавляет пользователей в БД
// Использует транзакцию
func addUsersTx(ctx context.Context, db *pgxpool.Pool, users []user) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	// отмена транзакции в случае ошибки
	defer tx.Rollback(ctx)
	// пакетный запрос
	batch := new(pgx.Batch)
	for _, u := range users {
		batch.Queue(`INSERT INTO users(name) VALUES ($1)`, u.name)
	}
	// отправка пакета в БД (может выполняться для транзакции или соединения)
	res := tx.SendBatch(ctx, batch)
	// обязательная операция закрытия соединения
	err = res.Close()
	if err != nil {
		return err
	}
	// подтверждение транзакции
	err = tx.Commit(ctx)
	return err
}

// users возвращает всех пользователей.
func users(ctx context.Context, db *pgxpool.Pool) ([]user, error) {
	// запрос на выборку данных
	rows, err := db.Query(ctx, `
		SELECT * FROM users ORDER BY id;
	`)
	if err != nil {
		return nil, err
	}
	var users []user
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var u user
		err = rows.Scan(
			&u.id,
			&u.name,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		users = append(users, u)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return users, rows.Err()
}
