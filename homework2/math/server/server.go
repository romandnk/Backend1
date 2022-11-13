package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type client chan<- string

var (
	entering   = make(chan client)
	leaving    = make(chan client)
	messages   = make(chan string)
	game       = true
	ans        int
	expression string
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Print("Server has started")
	<-done
	log.Print("Server has stopped")

	go broadcaster()
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Print(err)
				continue
			}
			go handleConn(conn)
		}
	}()
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	nickName := ""
	ch := make(chan string)
	go clientWriter(conn, ch)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		nickName = input.Text()
		break
	}

	ch <- "You are " + nickName
	messages <- nickName + " has arrived"
	entering <- ch

	for input.Scan() {
		if game {
			generateExpression()
			messages <- nickName + ": " + input.Text()
			playerAns, err := strconv.Atoi(input.Text())
			if err != nil {
				continue
			}
			if playerAns == ans {
				messages <- "The game is over. " + nickName + " won!"
				game = false
			}
		}
	}

	leaving <- ch
	messages <- nickName + " has left"
}

func broadcaster() {
	clients := make(map[client]bool)

	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func generateExpression() {
	game = true
	operators := []rune("+-*")
	firstOperator := operators[rand.Intn(len(operators))]
	secondOperator := operators[rand.Intn(len(operators))]

	firstNumber := rand.Intn(100)
	secondNumber := rand.Intn(100)
	thirdNumber := rand.Intn(100)

	switch string(firstOperator) {
	case "+":
		ansFirstTwoNumbers := firstNumber + secondNumber
		switch string(secondOperator) {
		case "+":
			expression = strconv.Itoa(firstNumber) + string(firstOperator) +
				strconv.Itoa(secondNumber) + string(secondOperator) + strconv.Itoa(thirdNumber)
			ans = ansFirstTwoNumbers + thirdNumber
			messages <- "Solver the expression: " + expression
		case "-":
			expression = strconv.Itoa(firstNumber) + string(firstOperator) +
				strconv.Itoa(secondNumber) + string(secondOperator) + strconv.Itoa(thirdNumber)
			ans = ansFirstTwoNumbers - thirdNumber
			messages <- "Solver the expression: " + expression
		case "*":
			expression = strconv.Itoa(firstNumber) + string(firstOperator) +
				strconv.Itoa(secondNumber) + string(secondOperator) + strconv.Itoa(thirdNumber)
			ans = ansFirstTwoNumbers * thirdNumber
			messages <- "Solver the expression: " + expression
		}
	case "-":
		ansFirstTwoNumbers := firstNumber - secondNumber
		switch string(secondOperator) {
		case "+":
			expression = strconv.Itoa(firstNumber) + string(firstOperator) +
				strconv.Itoa(secondNumber) + string(secondOperator) + strconv.Itoa(thirdNumber)
			ans = ansFirstTwoNumbers + thirdNumber
			messages <- "Solver the expression: " + expression
		case "-":
			expression = strconv.Itoa(firstNumber) + string(firstOperator) +
				strconv.Itoa(secondNumber) + string(secondOperator) + strconv.Itoa(thirdNumber)
			ans = ansFirstTwoNumbers - thirdNumber
			messages <- "Solver the expression: " + expression
		case "*":
			expression = strconv.Itoa(firstNumber) + string(firstOperator) +
				strconv.Itoa(secondNumber) + string(secondOperator) + strconv.Itoa(thirdNumber)
			ans = ansFirstTwoNumbers * thirdNumber
			messages <- "Solver the expression: " + expression
		}
	case "*":
		ansFirstTwoNumbers := firstNumber * secondNumber
		switch string(secondOperator) {
		case "+":
			expression = strconv.Itoa(firstNumber) + string(firstOperator) +
				strconv.Itoa(secondNumber) + string(secondOperator) + strconv.Itoa(thirdNumber)
			ans = ansFirstTwoNumbers + thirdNumber
			messages <- "Solver the expression: " + expression
		case "-":
			expression = strconv.Itoa(firstNumber) + string(firstOperator) +
				strconv.Itoa(secondNumber) + string(secondOperator) + strconv.Itoa(thirdNumber)
			ans = ansFirstTwoNumbers - thirdNumber
			messages <- "Solver the expression: " + expression
		case "*":
			expression = strconv.Itoa(firstNumber) + string(firstOperator) +
				strconv.Itoa(secondNumber) + string(secondOperator) + strconv.Itoa(thirdNumber)
			ans = ansFirstTwoNumbers * thirdNumber
			messages <- "Solver the expression: " + expression
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
