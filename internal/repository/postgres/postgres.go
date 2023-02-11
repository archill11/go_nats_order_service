package postgres

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"sync"

	"myapp/internal/config"
	"myapp/internal/entity"
	"myapp/internal/repository"

	_ "github.com/jackc/pgx/v4/stdlib"
)

//go:embed schema.sql
var InitTableQuery string

// явная имплементация интерфейса
var _ repository.OrderRepo = &Database{}

// Database - хранилище заказов.
type Database struct {
	db      *sql.DB
	storeCh chan entity.OrderEntity
	wg      *sync.WaitGroup
}

// NewStorage подключается к базе данных postgres и запускает воркер, в котором хранятся все входящие заказы.
func New(conf config.Config) (*Database, error) {
	dbc := conf.Database

	databaseURI := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		dbc.Username, dbc.Password, dbc.Host, dbc.Port, dbc.DBname,
	)

	db, err := sql.Open(dbc.DriverName, databaseURI)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil { // проверка что есть подключеие к БД
		return nil, err
	}
	if _, err := db.Exec(InitTableQuery); err != nil { //создаем таблицу в БД
		return nil, err
	}

	storage := &Database{
		db:      db,
		storeCh: make(chan entity.OrderEntity, 10),
		wg:      &sync.WaitGroup{},
	}
	storage.wg.Add(1)
	go storage.InsertWorker()
	return storage, nil
}

// Get Метод возвращает заказ по uid из БД
func (s *Database) Get(orderUID string) (entity.OrderEntity, error) {
	var order entity.OrderEntity

	q := `SELECT json_order FROM orders WHERE uid = $1;`
	err := s.db.QueryRow(q, orderUID).Scan(&order)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.OrderEntity{}, repository.ErrNotFound
		}
		return entity.OrderEntity{}, err
	}
	return order, nil
}

func (s *Database) GetAll() ([]entity.OrderEntity, error) {
	orders := make([]entity.OrderEntity, 0)

	q := `SELECT uid, json_order FROM orders;`
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order entity.OrderEntity
		if err := rows.Scan(order.GetUIDPtr(), order.GetJsonDataPtr()); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// Метод отправляет заказ в канал для созранения в БД
func (s *Database) AddOrder(orderUID, order string) error {
	newEntity := entity.NewOrderEntity(orderUID, order) // создаем новую сущность
	s.storeCh <- newEntity // пихаем новую сущность в канал
	return nil
}

// InsertWorker Метод слушает storeCh канал и добавляет в БД все поступающие заказы
func (s *Database) InsertWorker() {
	log.Println("storage: postgres: storer started")

	q := `INSERT INTO orders (uid, json_order) VALUES ($1, $2)`

	for newEntity := range s.storeCh { // достаем новую сущность из канала
		_, err := s.db.Exec(q, newEntity.GetUID(), newEntity.GetJsonData())
		if err != nil {
			log.Printf("postgresDB: ERR: could not save the order %s: %s", newEntity.GetUID(), err)
		} else {
			log.Printf("postgresDB: save %s order", newEntity.GetUID())
		}
	}
	// произойдет если закрыть канал
	log.Println("postgresDB: DB stopped")
	s.wg.Done()
}

// CloseDb Метод закрывает соединение с БД
func (s *Database) CloseDb() error {
	close(s.storeCh)
	s.wg.Wait()
	return s.db.Close()
}
