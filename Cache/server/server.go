package server

import (
	"cache/cache"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Кеш
var work_cache *cache.Cache

// Логгер для логирования запросов
func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s\n", r.Method, r.URL)
		next(w, r)
	}
}

func RunServer() error {
	// Инициализация кеша
	work_cache = work_cache.CreatingTheCache(5 - 1)
	//GET
	http.HandleFunc("/add", loggerMiddleware(handleEvent))
	http.HandleFunc("/get", loggerMiddleware(handleEvent))
	http.HandleFunc("/remove", loggerMiddleware(handleEvent))
	http.HandleFunc("/clear", loggerMiddleware(handleEvent))
	http.HandleFunc("/cap", loggerMiddleware(handleEvent))
	http.HandleFunc("/add_with_ttl", loggerMiddleware(handleEvent))
	http.HandleFunc("/get_all", loggerMiddleware(handleEvent))
	return http.ListenAndServe("localhost:8080", nil)
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	url_resuest := r.URL.Path
	switch url_resuest {
	case "/add":
		Add(w, r)
	case "/get":
		Get(w, r)
	case "/remove":
		Remove(w, r)
	case "/clear":
		Clear()
	case "/cap":
		Cap(w, r)
	case "/add_with_ttl":
		AddWithTTL(w, r)
	case "/get_all":
		work_cache.GetAll()
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		fmt.Fprintf(w, "Не передан ключ для получения элемента из кеша")
	}
	value, flag := work_cache.Get(key)
	fmt.Fprintf(w, "Ключ: %v Значение: %v", value, flag)
}

func Add(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		fmt.Fprintf(w, "Ключ для добавление элемента в кеш не передан")
	}
	value := r.URL.Query().Get("value")
	if value == "" {
		fmt.Fprintf(w, "Значение для добавление элемента в кеш не передано")
	}
	work_cache.Add(key, value)
}

func Remove(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		fmt.Fprintf(w, "Не передан ключ для удаления элемента из кеша")
	}
	work_cache.Remove(key)
}

func Clear() {
	work_cache.Clear()
}

func Cap(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Capacity: %v", work_cache.Cap())
	response := fmt.Sprintf("Capacity: %s", work_cache.Cap())
	fmt.Fprintf(w, response)

}

func AddWithTTL(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		fmt.Fprintf(w, "Ключ для добавление элемента в кеш не передан")
	}
	value := r.URL.Query().Get("value")
	if value == "" {
		fmt.Fprintf(w, "Значение для добавление элемента в кеш не передано")
	}

	duration := r.URL.Query().Get("duration")
	if duration == "" {
		fmt.Fprintf(w, "Время жизни для элемента кеша не задано")
	}
	i, err := strconv.Atoi(duration)
	if err != nil {
		// ... handle error
		panic(err)
	}

	var d time.Duration = time.Duration(i) * time.Second

	work_cache.AddWithTTL(key, value, d)
}
