# Bugbot
---
## Overview

This repository is included among the projects developed by the 2022 IT class students at BHH & ITECH Hamburg.

BugBot is a chatbot designed to facilitate customer communication with a chosen company regarding issues with its products.
---
## Project Description

The features of the chatbot have been designed such that all components are implemented as microservices, utilizing the following languages and technologies:

| Name     | Description                                             | Services                                       |
| -------- | ------------------------------------------------------- | ---------------------------------------------- |
| Golang   | native & imported modules used                          | Client, API, Tokenizer, Matcher, E-Mail client |
| HTML     | used for a webpage                                      | Web-Frontend                                   |
| MAUI.NET | used for a windows/android app                          | Application                                    |
| MongoDB  |                                                         | DB, Unsolved DB                                |
| Mail-Pit | for testing purposes of the mail toolchain              | E-Mail client                                  |
| Docker   | Containerization of all dedicated services              | Containerization of all dedicated services     |
| Ansible  | using IaC to deploy the whole project with one playbook |                                                |
|Javascript| creating a webapplication                               | Webinterface                                   |

---
### Services

The chatbot includes the following services:

    **Client**: This service handles user input and communicates with the Rest-API to forward the user’s message.

    **API**: This service receives data from the client, marshals it into a JSON format, and sends it to the Tokenizer.

    **Tokenizer**: This service tokenizes the marshaled data received from the Rest-API and forwards it to the Matcher.

    **Matcher**: This service identifies keywords in the tokenized data and queries a database to compare them with existing entries. The response from the database is then returned to the API.

    **Web Application** & **App**: This service displays the response from the API as a string to the user.

    **Unsolved Database**: If the Matcher is unable to find a match for the input, the input is stored in this database for future reference.

    **Response Database**: Knowledge Database in which the matcher identifies solutions

    **Email Client**: This service automatically sends an email to an employee as a ticket when the Matcher is unable to find a match for the input.

    **Mailpit**: For testing the mail client
   

In addition to the web application, there is also a standalone application available for Android and Windows platforms.
---
### Workflow

The workflow of the chatbot is as follows:

    1. The client receives the user’s input and sends it to the Rest-API.
    2. The Rest-API marshals the data into a JSON format and sends it to the Tokenizer.
    3. The Tokenizer tokenizes the data and sends it to the Matcher.
    4. The Matcher identifies keywords in the tokenized data and queries a database to compare them with existing entries.
    5. If a match is found, the Matcher sends the response back to the API, which then sends it to the client to be displayed in the web application.
    6. If a match is not found, the Matcher sends the input to the Unsolved Database and the Email Client, which sends an email to an employee as a ticket.

This modular design allows each component to be developed, tested, and deployed independently, improving the overall robustness and scalability of the chatbot.

---
## Deploy

Clone the repository and copy all .env_examples into your .env-Files and change them to your satisfaction.

To start the program enter '''ansible-playbook deploy.yaml'''

---
## Module

This project is related to the module "Progamminglanguages and Methods" of the BHH/ITECH.



