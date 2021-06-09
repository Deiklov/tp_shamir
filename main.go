//Процессу на stdin приходят строки, содержащие интересующие нас URL. Каждый такой URL нужно дернуть и посчитать кол-во вхождений строки "Go" в ответе. В конце работы приложение выводит на экран общее кол-во найденных строк "Go" во всех запросах, например:
//$ echo -e 'https://golang.org\nhttps://golang.org\nhttps://golang.org' | go counter.go
//Count for https://golang.org: 9
//Count for https://golang.org: 9
//Count for https://golang.org: 9
//Total: 27
//
//Введенный URL должен начать обрабатываться сразу после вычитывания и параллельно с вычитыванием следующего. URL должны обрабатываться параллельно, но не более k=5 одновременно. Обработчики url-ов не должны порождать лишних горутин, т.е. если k=1000 а обрабатываемых URL-ов нет, не должно создаваться 1000 горутин. Нужно обойтись без глобальных переменных и использовать только стандартные библиотеки.
//Поддерживает ключи: k - количество горутин обрабатывающих урлы (по дефолту 5) q - поисковое слово (по дефлоту go)
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func HandleCommonString(common string, tasks chan<- string) {
	scanner := bufio.NewScanner(strings.NewReader(common))
	for scanner.Scan() {
		//прочитали строку и скинули пулу
		//fmt.Println(scanner.Text())
		tasks <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("closing tasks queue")
	close(tasks)
}

func worker(group *sync.WaitGroup, tasks <-chan string, results chan<- int) {
	client := http.Client{
		Timeout: 15 * time.Second,
	}
	//все воркеры смотрят в один канал, забирают из него таски, пишут в общий канал results
	for v := range tasks {
		//делаем http request и подсчет числа слов
		resp, err := client.Get(v)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		html, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		//отправляем кол-во слов в стопку с результатами
		countWrd := strings.Count(string(html), "go")

		fmt.Printf("Count for %s : %d\n", v, countWrd)
		results <- countWrd
		//ждем, пока таски не закончатся
	}
	group.Done()
}

func main() {
	var wg sync.WaitGroup
	n := 5
	tasks := make(chan string, 7)
	results := make(chan int, 7)
	input := "https://golang.org/doc/\nhttps://golang.org\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/\nhttps://golang.org\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/\nhttps://golang.org\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/\nhttps://golang.org\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/\nhttps://golang.org\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/\nhttps://golang.org\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/\nhttps://golang.org\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/\nhttps://golang.org\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/effective_go\nhttps://golang.org/doc/effective_go\n"

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(&wg, tasks, results)
	}

	go HandleCommonString(input, tasks)

	//ждем пока все воркеры отработают
	var sum int

	go func() {
		wg.Wait()
		close(results)
	}()

	//читаем уже из закрытого канала
	for v := range results {
		//fmt.Println(v)
		sum += v
	}

	fmt.Printf("total %d\n", sum)
}
