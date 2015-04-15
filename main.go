package main

import (
	"encoding/json"
	"fmt"
	"github.com/s-kostyaev/webtop-protocol"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

var sharedEncoder *json.Encoder

func main() {
	setupLogger()
	reconnect()
	http.Handle("/", http.FileServer(http.Dir(staticDir)))
	http.HandleFunc("/command/", handleCommand)
	http.ListenAndServe(":3000", nil)
}

func reconnect() {
	connection, err := net.Dial("unix", socketPath)
	if err != nil {
		logger.Error(err.Error())
	}
	sharedEncoder = json.NewEncoder(connection)
	go answerReader(connection)
}

var answers = make(map[int]protocol.Answer)

func handleCommand(w http.ResponseWriter, r *http.Request) {
	newRequest := protocol.Request{}
	newRequest.Id = newId()
	containerIPs, err := net.LookupIP(r.Host)
	if err != nil {
		logger.Error(err.Error())
	}
	newRequest.Host = containerIPs[0].String()
	newRequest.Cmd = strings.TrimPrefix(r.URL.RequestURI(), "/command/")
	sharedEncoder.Encode(newRequest)
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
		logger.Error(answer.Status, answer.Error)
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprint(w, answer.Data)
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
	buf := make([]byte, 2048)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			logger.Error(err.Error())
			reconnect()
			return
		}
		currentAnswer := protocol.Answer{}
		err = json.Unmarshal(buf[0:n], &currentAnswer)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		answers[currentAnswer.Id] = currentAnswer
	}
}
