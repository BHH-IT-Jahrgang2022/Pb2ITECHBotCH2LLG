package mailClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	//"net/smtp"
)

type Ticket struct {
	Tags    []string `json:"tags"`
	Problem string   `json:"problem"`
}

type InputData struct {
	Input string `json:"input"`
	Timestamp string `json:"timestamp"`
}

/*func SendEmail(ticket *Ticket) {
	to := "botbugland@gmail.com"
	subject := "New Ticket"
	body := "Tags: " + strings.Join(ticket.Tags, ", ") + "\nProblem: " + ticket.Problem

	msg := "To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

//buglandbot@gmail.com
//DeinerMudder123
// Lower security settings in the gmail app -> https://support.google.com/mail/thread/5621336/bad-credentials-using-gmail-smtp?hl=en

	// Set up authentication information.
	smtpHost := "smtp.gmail.com."
	smtpPort := "587"
	smtpUser := "botbugland@gmail.com"
	smtpPass := "DeineMudder123"
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("Email sent")
	fmt.Println("")
}*/

func PrintEmail(ticket *Ticket) {
	to := "buglandbot@gmail.com"
	subject := "New Ticket"
	body := "Tags: " + strings.Join(ticket.Tags, ", ") + "\nProblem: " + ticket.Problem

	msg := "To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	fmt.Println("Email Content: \n" + msg) // Print the email content to the console
}


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

func ReceiveAndPrintTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var inputData InputData
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticket := &Ticket{
		Tags:    []string{inputData.Input, inputData.Timestamp},
		Problem: "Received input from endpoint",
	}

	PrintEmail(ticket)
}

func FetchAndPrintTicket() {
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
		PrintEmail(&ticket)
	}
}