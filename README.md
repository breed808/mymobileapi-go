# mymobileapi-go

Go API Client to interact with the mymobileapi [REST API](https://mymobileapi.readme.io/reference).

The current implementation requires Go 1.14 and uses modules.

## Install 

`import "github.com/breed808/mymobileapi-go"`

## Usage

```go
client := mymobileapi.Client()
client, err := gandiapi.NewClient("Client ID", "supersecret", false)
if err != nil {
	log.Fatal("error")
}

// Query account balance
_, err := client.GetBalance()
if err != nil {
	log.Fatal(err)
}

// Send SMS message to recipient(s)
req := BulkMessageRequest{
    Messages: []Message{
        Message{
            Content:     "Hello World!",
            Destination: "0412345678",
        },
    },
}

resp, err := client.SendBulkMessages(req)
if err != nil {
	log.Fatal(err)
}
```

Note that when the client is initialised, the library automatically authenticates to mymobileapi and generates an authentication token.
This token expires after 24 hours. Short-lived scripts/processes using this library need no be concerned of the expiry, but longer-running scripts/processes will need to re-authenticate every 24 hours.

```go
// Reauthenticate if authentication token has expired
if time.Now().After(client.AuthTokenExpiry) {
    client.Authenticate()
}

// Query account balance
_, err := client.GetBalance()
if err != nil {
	log.Fatal(err)
}
```
