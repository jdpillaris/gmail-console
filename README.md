# gmail-console
Console application to retrieve items in Gmail mailbox for specific time frame

To test run the application, run the following command:

```
go run main.go
```

To create the final binary for deployment:

```
go build
```

# Application Decisions
1. Parameters for this console app are: mailbox name, mailbox password, start of retrieval period, end of retrieval period. 
Since the mailbox name/password can change, they can be provided in the default HTML browser and not from the command 
line. Hence the console prompts the user for 2 parameters - start & end of retrieval period.
2. This default browser only for the  launched the first time to authenticate with Google's OAuth2 authorization server. 
So, a token file is generated oinly for the first time and overwritten for every new account login .
3. Each email item is stored into a file which takes the email ID as its name.
4. The content of each email item is taken directly from the *gmail.Message struct defined in Gmail API library
5. For encryption, a fixed passphrase "password" is used in downloadItems() function inside controllers/inbox.go. 
Using a simple MD5 hash, this passphrase produces 32 byte hashes.
