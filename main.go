package main

import (
	"encoding/json"
	"fmt"
	"github.com/s-kostyaev/webtop-protocol"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	setupLogger()
	go reconnect()
	http.Handle("/", http.FileServer(http.Dir(staticDir)))
	http.HandleFunc("/command/", handleCommand)
	http.ListenAndServe(":80", nil)
}

func reconnect() {
	// remove socket file if exist
	os.Remove(dataSocketPath)
	listener, err := net.Listen("unix", dataSocketPath)
	if err != nil {
		logger.Error(err.Error())
		go reconnect()
		return
	}
	connection, err := listener.Accept()
	if err != nil {
		logger.Error("Access error: %s\n", err.Error())
		go reconnect()
		return
	}

	go answerReader(connection)
}

var answers = make(map[int]protocol.Answer)

func handleCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	newRequest := protocol.Request{}
	newRequest.Id = newId()
	containerIPs, err := net.LookupIP(r.Host)
	if err != nil {
		logger.Error(err.Error())
	}
	newRequest.Host = containerIPs[0].String()
	newRequest.Cmd = strings.TrimPrefix(r.URL.RequestURI(), "/command/")
	connection, err := net.Dial("unix", commandSocketPath)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, protocol.Answer{}.Data)
		delete(answers, newRequest.Id)
		return
	}
	defer connection.Close()
	encoder := json.NewEncoder(connection)

	encoder.Encode(newRequest)
	startTime := time.Now()
	// FIXME: get timeout from config
	timeout := (3 * time.Second)
	answer, ok := answers[newRequest.Id]
	for !ok {
		time.Sleep(30 * time.Millisecond)
		answer, ok = answers[newRequest.Id]
		// if no answer send error after timeout from config
		if time.Since(startTime) >= timeout {
			break
		}
	}
	if answer.Status != "ok" {
		logger.Error("%#v %#v", answer.Status, answer.Error)
		w.WriteHeader(http.StatusInternalServerError)
	}
	buf, _ := json.Marshal(answer.Data)
	w.Write(buf)
	delete(answers, newRequest.Id)
}

const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)
const minInt = -maxInt - 1

var currentId = minInt

func newId() int {
	if currentId == maxInt {
		currentId = minInt
	}
	currentId++
	return currentId
}

func answerReader(r io.Reader) {
	decoder := json.NewDecoder(r)
	for {
		currentAnswer := protocol.Answer{}
		err := decoder.Decode(&currentAnswer)
		if err != nil {
			if err.Error() != "EOF" {
				logger.Error(err.Error())
			}
			go reconnect()
			return
		}
		answers[currentAnswer.Id] = currentAnswer
	}
}
