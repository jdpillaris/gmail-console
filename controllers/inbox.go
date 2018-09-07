package controllers

import (
	"app/utils"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"

	"github.com/skratchdot/open-golang/open"
)

// Inbox controller
type Inbox struct {
	gmailService *gmail.Service
	user 		 string
}

func (c *Inbox) Create() {
	c.user = "me"

	credAbsPath, _ := filepath.Abs("assets/credentials.json")
	b, err := ioutil.ReadFile(credAbsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	c.gmailService = srv
}

func (c *Inbox) StoreMessages(startDate, endDate time.Time) {
	startTime := startDate.Format("2006/01/02")
	endTime := endDate.Format("2006/01/02")

	mailIds := c.getMailIds(startTime, endTime)

	c.downloadItemsInParallel(mailIds)
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokenFile := "token.json"
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	// fmt.Printf("Go to the following link in your browser then type the "+
	// 		"authorization code: \n%v\n", authURL)
	fmt.Printf("Type the authorization code from your browser: ")
	err := open.Run(authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

// func (c *Inbox) getMailIds() []string {
func (c *Inbox) getMailIds(startTimeStr, endTimeStr string) []string {
	mailIds := make([]string, 0)

	queryStr := "after:" + startTimeStr + " before:" + endTimeStr
	r, err := c.gmailService.Users.Messages.List(c.user).Q(queryStr).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}
	if len(r.Messages) == 0 {
		fmt.Println("No messages found.")
		return mailIds
	}
	
	fmt.Println("Storing message IDs")
	for _, m := range r.Messages {
		mailIds = append(mailIds, m.Id) 
	}

	return mailIds
}

// downloadItems downloads email items using email IDs
func (c *Inbox) downloadItems(mailIds []string) {
	start := time.Now()

	if len(mailIds) <= 0 {
		fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
		return
	}

	passPhrase := "password"
	myHash := utils.CreateHash(passPhrase)

	for _, mailId := range mailIds {
		fmt.Printf("Downloading mail ID - %s\n", mailId)
		p, _ := c.gmailService.Users.Messages.Get(c.user, mailId).Do()
		utils.EncryptFile(p, myHash)

		// For testing purposes: retrieve an encrypted file
		// utils.DecryptFile(p.Id, myHash)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

// downloadItemsInParallel downloads email items using email IDs "in parallel"
func (c *Inbox) downloadItemsInParallel(mailIds []string) {
	start := time.Now()

	if len(mailIds) <= 0 {
		fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
		return
	}

	passPhrase := "password"
	myHash := utils.CreateHash(passPhrase)

	ch := make(chan *gmail.Message)
	for _, mailId := range mailIds {
		go c.getMessagesInParallel(mailId, myHash, ch)
	}

	for range mailIds {
		utils.EncryptFile(<-ch, myHash)
		// For testing purposes: retrieve an encrypted file
		// utils.DecryptFile(p.Id, myHash)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func (c *Inbox) getMessagesInParallel(mailID, hashKey string, ch chan<-*gmail.Message) {
	fmt.Printf("Downloading mail ID - %s\n", mailID)
	p, _ := c.gmailService.Users.Messages.Get(c.user, mailID).Do()

	ch <- p
}
