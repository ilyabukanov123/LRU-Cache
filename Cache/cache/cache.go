package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// Структура для одного элемента в кеше
type Item struct {
	Value      interface{}   // Значение
	Key        interface{}   // Ключ
	Created    time.Time     // Дата создания
	Expiration time.Duration // Время жизни кеша
}

// Интерфейс сервиса кэширования
type ICache interface {
	Cap() int                                                                                        // Получение capacity
	Clear()                                                                                          // Очистка кеша
	Add(key, value interface{})                                                                      // Добавление элемента в кэш
	AddWithTTL(key, value interface{}, ttl time.Duration)                                            // Добавление элемента в кэш с учетом времени жизни
	Get(key interface{}) (value interface{}, ok bool)                                                // Получение элемента из кеша
	Remove(key interface{})                                                                          // Удаление элемента из кеша
	CreatingTheCache(defaultExpiration time.Duration, cleanupInterval time.Duration, cap int) *Cache // Инициализация нового кеша
}

// Структура для кэша
type Cache struct {
	ICache
	sync.RWMutex // Мьютекс для безопасного доступа к данным во время чтения/записи
	items        map[interface{}]*list.Element
	cap          int // capacity
	list         *list.List
}

// Функция по созданию кэша; defaultExpiration - время жизни кэша по умолчанию, cleanupInterval - интервал между запуском Garbage Collector
func (c *Cache) CreatingTheCache(cap int) *Cache {
	// Создаем структуру кэша
	cache := Cache{
		cap:   cap,
		items: make(map[interface{}]*list.Element),
		list:  list.New(),
	}

	// Возвращаем кеш
	return &cache
}

// Метод с реализацией фунции интерфейса на добавление элемента в кэш
func (c *Cache) Add(key interface{}, value interface{}) {
	// Блокируем кэш чтобы не повредить память
	c.Lock()
	defer c.Unlock()

	// Если в кеше уже есть элемент с таким ключом помещаем его в начало списка и меняем ему значение
	if entry, exists := c.items[key]; exists {
		c.list.MoveToFront(entry)
		entry.Value = value
	}

	// Если длина двусвязного спискаравно capacity
	if c.list.Len() == c.cap {
		c.ClearItems()
	}

	if c.list.Len() == c.cap {
		c.ClearItems()
	}
	// Заполняем новый элемент в кеш
	item := Item{
		Key:     key,
		Value:   value,
		Created: time.Now(), // Время добавления элемента в кэш
	}

	// Данная конструкция создает новый элемент двусвязного списка и добавляет его в начало списка.
	//Переменная item содержит значение, которое будет храниться в этом элементе.
	//Переменная element будет содержать указатель на созданный элемент списка
	element := c.list.PushFront(item)
	c.items[key] = element
}

// Метод с реализацией функции интерфейса на добавление элемента в кэш с учетом время жизни элемента в кэше
func (c *Cache) AddWithTTL(key string, value interface{}, duration time.Duration) {

	// Блокируем кэш чтобы не повредить память
	c.Lock()
	defer c.Unlock()
	// Если в кеше уже есть элемент с таким ключом помещаем его в начало списка и меняем ему значение
	if entry, exists := c.items[key]; exists {
		c.list.MoveToFront(entry)
		entry.Value = value
	}

	// Если длина двусвязного списка равно capacity
	if c.list.Len() == c.cap {
		c.ClearItems()
	}
	// Заполняем новый элемент в кеш
	item := Item{
		Key:        key,
		Value:      value,
		Expiration: duration,
		Created:    time.Now(), // Время добавления элемента в кэш
	}

	// Данная конструкция создает новый элемент двусвязного списка и добавляет его в начало списка.
	//Переменная item содержит значение, которое будет храниться в этом элементе.
	//Переменная element будет содержать указатель на созданный элемент списка
	element := c.list.PushFront(item)
	c.items[key] = element

	// Запускаем горутину с таймером для элемента кеша
	go func() {
		//https://www.geeksforgeeks.org/time-newticker-function-in-golang-with-examples/
		ticker := time.NewTicker(duration)
		for {
			// Удаляем элемент как только заканчивается время
			<-ticker.C
			c.Remove(key)
		}
	}()
}

// Метод по очистке элемента из кеша
func (c *Cache) ClearItems() {
	// Получаем последний элемент из двусвязного списка
	element := c.list.Back()
	if element != nil {
		// Удаляем последний (к которому больше всего не было обращений) элемент из списка
		item := c.list.Remove(element).(Item)
		// Удаляем элемент из мэпы
		delete(c.items, item.Key)
	}
}

// Метод с реализацией фунции интерфейса на получение элемента из кэша по ключу
func (c *Cache) Get(key string) (interface{}, bool) {

	// Блокируем память
	c.RLock()

	defer c.RUnlock()

	// Получаем элемент из кеша
	item, flag := c.items[key]

	// Если ключ будет не найден - вернется false
	if !flag {
		return nil, false
	}

	// Во всех остальные случаях возвращаем значение и флаг, который сигнализирует о наличии кеша
	return item.Value, true
}

// Метод по получению всех объектов в кеше
func (c *Cache) GetAll() {
	for k, v := range c.items {
		fmt.Printf("Ключ: %v Значение: %v\n", k, v)
	}
}

// Метод с реализации функции интерфейса на очистку всего кэша
func (c *Cache) Remove(key interface{}) {
	c.Lock()
	defer c.Unlock()
	// Удаляем элемент из map по ключу
	delete(c.items, key.(string))
}

// // Метод с реализации функции интерфейса на очистку всего кэша
func (c *Cache) Clear() {
	c.Lock()
	defer c.Unlock()
	// Перебираем все элементы памы и удаляем элемень с каждой итерацией цикла
	for k := range c.items {
		delete(c.items, k)
	}
}

// Метод с реализацией функции интерфейса по получению capacity
func (c *Cache) Cap() int {
	return c.cap
}
