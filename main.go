package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func main() {
	req1 := make(chan string)
	req2 := make(chan string)

	go fetchAPI("https://viacep.com.br/ws/01153000/json/", req2, "viacep")
	go fetchAPI("https://brasilapi.com.br/api/cep/v1/01153000", req1, "brasilapi")

	select {
	case res := <-req1:
		fmt.Println("Resposta recebida da BrasilAPI: ", res)
	case res := <-req2:
		fmt.Println("Resposta recebida da ViaCEP: ", res)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: não tivemos resposta em 1 segundo.")
	}
}

func fetchAPI(url string, ch chan<- string, source string) {
	defer close(ch)
	resp, err := http.Get(url)
	if err != nil {
		ch <- "Erro ao fazer requisição para " + source + ": " + err.Error()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- "Erro de status code na " + source + ": " + strconv.Itoa(resp.StatusCode) + " " + resp.Status
		return
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- "Erro ao ler resposta da " + source + ": " + err.Error()
		return
	}

	ch <- string(res)
}
