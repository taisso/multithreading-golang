package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var client = http.Client{
	Timeout: time.Second,
}

func main() {
	ch := make(chan string)

	go RequestBrasilAPI("01153000", ch)
	go RequestViaCep("01153000", ch)

	select {
	case msg := <-ch:
		Print(msg)
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}

}

func Print(value string) {
	fmt.Println("O Endereço é:")
	fmt.Println(value)
}

func Request(url string) (map[string]string, error) {
	req, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	var address map[string]string
	if err = json.NewDecoder(req.Body).Decode(&address); err != nil {
		return nil, err
	}

	return address, nil
}

func RequestBrasilAPI(cep string, ch chan<- string) {
	startTime := time.Now()

	address, err := Request("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		panic(err)
	}

	sub := time.Since(startTime)
	ch <- fmt.Sprintf("Tempo de resposta: %v\nCEP: %s\nRua: %s\nCidade: %s\nEstado: %s\nAPI: BrasilAPI",
		sub,
		address["cep"],
		address["street"],
		address["city"],
		address["state"],
	)
}

func RequestViaCep(cep string, ch chan<- string) {
	startTime := time.Now()

	url := fmt.Sprintf("http://viacep.com.br/ws/%v/json/", cep)
	address, err := Request(url)
	if err != nil {
		panic(err)
	}

	sub := time.Since(startTime)
	ch <- fmt.Sprintf("Tempo de resposta: %v\nCEP: %s\nRua: %s\nCidade: %s\nEstado: %s\nAPI: ViaCep",
		sub,
		address["cep"],
		address["logradouro"],
		address["localidade"],
		address["uf"],
	)
}
