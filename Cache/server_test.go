package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestAddItemToCache(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/add?key=1&value=1")
	if err != nil {
		log.Fatalf("Произошла ошибка при отправке GET запроса на http://localhost:8080/add. Не передан ключ или значения для добавления элемента в кег: %s\n", err)
	} else {
		fmt.Println("Данные успешно добавлены в lru кеш")
	}
	fmt.Println(resp)
}

func TestGetItemToCache(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/get?key=1")
	if err != nil {
		log.Fatalf("Произошла ошибка при отправке GET запроса на http://localhost:8080/get. Не передан ключ: %s\n", err)
	}
	fmt.Println(resp)
}

func TestRemoveItemToCache(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/remove?key=1")
	if err != nil {
		log.Fatalf("Произошла ошибка при отправке GET запроса на http://localhost:8080/get. Не передан ключ: %s\n", err)
	}
	fmt.Println("Элемент удален по ключу из lru кеша", resp)
}

func TestClearAllItemToCache(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/remove?key=1")
	if err != nil {
		log.Fatalf("Произошла ошибка при отправке GET запроса на http://localhost:8080/get. Не передан ключ: %s\n", err)
	}
	fmt.Println("lru кеш очищен", resp)
}

func TestGetCapacityCache(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/cap")
	if err != nil {
		log.Fatalf("Произошла ошибка при отправке GET запроса на http://localhost:8080/get: %s\n", err)
	}
	fmt.Println("Capacity кеша: ", resp)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanRunes)
	var buf bytes.Buffer
	for scanner.Scan() {
		buf.WriteString(scanner.Text())
	}
	fmt.Println(buf.String())

}

func TestAddItemWithTTLToCache(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/add_with_ttl?key=1&value=1&duration=30")
	if err != nil {
		log.Fatalf("Произошла ошибка при отправке GET запроса на http://localhost:8080/add. Не передан ключ, значение, или время жизни для элемента кеша: %s\n", err)
	} else {
		fmt.Println("Данные успешно добавлены в lru кеш с учетом времени жихни")
	}
	fmt.Println(resp)
}
