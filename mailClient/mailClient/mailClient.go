package mailClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

type Ticket struct {
	Tags    []string `json:"tags"`
	Problem string   `json:"problem"`
}

func SendEmail(ticket *Ticket) {
	to := "botbugland@gmail.com"
	subject := "New Ticket"
	body := "Tags: " + strings.Join(ticket.Tags, ", ") + "\nProblem: " + ticket.Problem

	msg := "To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	// Set up authentication information.
	smtpHost := os.Getenv("SMTPHOST")
	smtpPort := os.Getenv("SMTPPORT")
	smtpUser := "botbugland@gmail.com"
	smtpPass := "DeineMudder123"
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	auth = nil

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("Email sent")
	fmt.Println("")
}

// func SendEmail(ticket *Ticket) {
// 	to := "buglandbot@gmail.com"
// 	subject := "New Ticket"
// 	body := "Tags: " + strings.Join(ticket.Tags, ", ") + "\nProblem: " + ticket.Problem

// 	msg := "To: " + to + "\n" +
// 		"Subject: " + subject + "\n\n" +
// 		body

// 	fmt.Println("Email Content: \n" + msg) // Print the email content to the console
// }

/*func FetchAndEmailTicket() {
	url := "http://" + os.Getenv("UNSOLVEDHOST") + ":" + os.Getenv("UNSOLVEDPORT") + "/data"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tickets []Ticket
	err = json.Unmarshal(body, &tickets)
	if err != nil {
		log.Fatal(err)
	}

	for _, ticket := range tickets {
		SendEmail(&ticket)

	}
}*/

func FetchAndPrintTickets() {
	url := "http://" + os.Getenv("UNSOLVEDHOST") + ":" + os.Getenv("UNSOLVEDPORT") + "/data"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tickets []Ticket
	err = json.Unmarshal(body, &tickets)
	if err != nil {
		log.Fatal(err)
	}

	for _, ticket := range tickets {
		SendEmail(&ticket)
	}
}
