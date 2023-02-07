package in_memory_cache

// in-memory cache

import (
	"fmt"
	"log"
	"sync"
	"time"

	"myapp/internal/entity"
	"myapp/internal/repository"

	"github.com/zekroTJA/timedmap"
)

// явная имплементация интерфейса
var _ repository.OrderRepo = &Cache{}

// Cache - хранилище заказов.
type Cache struct {
	timedRepo        *timedmap.TimedMap
	repository       sync.Map
	repositoryLength int
	generalStorage   repository.OrderRepo
}

// NewCache Метод принимает хранилище БД в параметре, создает на его основе in-memory кэш хранилице
// и возвращает его
func New(dbStore repository.OrderRepo) (*Cache, error) {
	store := &Cache{
		timedRepo: timedmap.New(time.Minute),
	}

	if err := store.FillCacheFromDB(store, dbStore); err != nil {
		return nil, fmt.Errorf("in_memory_cache: could not apply option: %w", err)
	}

	return store, nil
}

// FillCacheFromDB Метод заполяет in-memory кэш данными из БД
func (s *Cache) FillCacheFromDB(cacheStore *Cache, dbStore repository.OrderRepo) error {

	outCh := make(chan []entity.OrderEntity, 1)
	errsCh := make(chan error, 1)

	go func() {
		orders, err := dbStore.GetAll() // получаем все заказы из БД
		if err != nil {
			errsCh <- err
		} else {
			outCh <- orders
		}
		close(outCh)
		close(errsCh)
	}()

	cacheStore.generalStorage = dbStore

	select {
	case orders := <-outCh:
		for _, order := range orders { // добавляем все заказы в in-memory кэш
			go cacheStore.PutInCache(order.GetUID(), order.GetJsonData())
		}
		if len(orders) > 0 {
			log.Printf("in_memory_cache: add %d orders from the database", len(orders))
		}
		return nil

	case err := <-errsCh:
		return err
	}
}

// AddOrder Метод добавляет заказ в in-memory кэш и в БД
func (s *Cache) AddOrder(orderUID, jsonOrder string) error {
	err := s.PutInCache(orderUID, jsonOrder)
	if err != nil {
		return err
	}

	if s.generalStorage != nil {
		s.generalStorage.AddOrder(orderUID, jsonOrder)
	}
	return nil
}

// PutInCache Метод добавляет заказ в in-memory кэш
func (s *Cache) PutInCache(orderUID, jsonOrder string) error {

	_, ok := s.repository.Load(orderUID)
	if ok {
		return repository.ErrAlreadyExists
	}

	s.repository.Store(orderUID, jsonOrder) // добавляем заказ в in-memory кэш
	s.repositoryLength = 0                  // удаляем длинну in-memory кэш репозитория
	return nil
}

// Get Метод возвращает заказ по uid из in-memory кэш
func (s *Cache) Get(orderUID string) (entity.OrderEntity, error) {

	orderJsonb, ok := s.repository.Load(orderUID)
	if !ok { // если в in-memory кэш нет заказа по такому ключу
		ordEntity, ok := s.searchOrderIntimedRepo(orderUID) // ищу заказ в timedRepo
		if !ok {
			return entity.OrderEntity{}, repository.ErrNotFound
		} else {
			orderJsonb = ordEntity.GetJsonData()
			s.PutInCache(orderUID, orderJsonb.(string))
		}
	}

	order := entity.NewOrderEntity(orderUID, orderJsonb.(string))

	return order, nil
}

// GetAll Метод возвращает все заказы из in-memory кэш
func (s *Cache) GetAll() ([]entity.OrderEntity, error) {

	var cacheLen int
	if s.repositoryLength > 0 { // если значение длинны in-memory кэш есть берем его
		cacheLen = s.repositoryLength
	} else {
		cl := s.CacheLen(&s.repository) // иначе считаем и кэшируем
		cacheLen = *cl
	}

	orders := make([]entity.OrderEntity, 0, cacheLen)

	s.repository.Range(func(uid, jsonb interface{}) bool {
		orderEntity := entity.NewOrderEntity(uid.(string), jsonb.(string))

		orders = append(orders, orderEntity)
		return true
	})

	return orders, nil
}

// CacheLen Метод возвращает длинну in-memory кэш репозитория
func (s *Cache) CacheLen(sMap *sync.Map) *int {
	var cacheMapLen int
	s.repository.Range(func(k, v interface{}) bool {
		cacheMapLen++
		return true
	})
	// кэшируем высчитанное значение, оно будет удалено при следующей оперции добавления
	s.repositoryLength = cacheMapLen

	return &cacheMapLen
}

// Метод проверяет был ли поиск этого заказа менее чем 5 минут назад,
// если этот заказ уже пытались найти менее чем 5 минут назад, он возвращает false,
// иначе он сделает запрос в БД и вернет результат.
// Если в БД тоже нет заказа, он добавит значение в timedRepo на 5 минут,
// что бы не ходить в БД чаще чем раз в 5 минут
func (s *Cache) searchOrderIntimedRepo(orderUID string) (entity.OrderEntity, bool) {

	_, ok := s.timedRepo.GetValue(orderUID).(bool) // проверяю uid заказа в timedRepo
	if !ok {                                       // если там нет такого ключа
		ordEntity, err := s.generalStorage.Get(orderUID) // делаю запрос в БД
		if err != nil {                                  // если в БД нет такого заказа
			s.timedRepo.Set(orderUID, true, time.Minute*5) // кладу ключ orderUID и устанавливаю значение время
			return entity.OrderEntity{}, false             // когда он удалится из timedRepo
		} else {
			return ordEntity, true
		}
	} else {
		return entity.OrderEntity{}, false
	}
}
