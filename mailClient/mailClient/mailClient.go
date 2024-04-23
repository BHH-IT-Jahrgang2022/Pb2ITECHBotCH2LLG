package mailClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Ticket struct {
	Tags    []string `json:"tags"`
	Problem string   `json:"problem"`
}

func SendEmail(ticket *Ticket) {
	to := "buglandbot@pm.me"
	subject := "New Ticket"
	body := "Tags: " + strings.Join(ticket.Tags, ", ") + "\nProblem: " + ticket.Problem

	msg := "To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	cmd := exec.Command("/usr/sbin/sendmail", "-t")
	cmd.Stdin = strings.NewReader(msg)
	err := cmd.Run()

	if err != nil {
		log.Printf("sendmail error: %s", err)
		return
	}

	log.Print("Email sent")
	fmt.Println("")
}

func FetchAndEmailTicket() {
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