package memdb

import "db/apps/pkg/storage/postgres"

// ПОльзовательский тип данных - реализация БД в памяти.
// "заглушка"
type DB []postgres.Task

// Выполнений контракта интерфейса storage.Interface
func (db DB) Tasks(int, int) ([]postgres.Task, error) {
	return db, nil
}
func (db DB) NewTask(postgres.Task) (int, error) {
	return 0, nil
}
