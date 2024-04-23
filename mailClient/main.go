package main

import ("mailClient/mailClient"
		"net/http"
)

func main() {
	//mailClient.FetchAndEmailTicket()
	mailClient.FetchAndPrintTicket()
	http.HandleFunc("/trigger", mailClient.ReceiveAndPrintTicket)
}