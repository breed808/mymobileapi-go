package mymobileapi

// The SMS message that should be sent. Must contain at least one message and at most 100 messages.
type Message struct {
	// The textual content of the SMS message
	Content string
	// The MSISDN (mobile number) to whom the messages is intended
	Destination string
	// A user defined ID which can be used to correlate messages sent to the API and receipts or replies
	// received for that message. Maximum length 100 characters
	CustomerId string
	// Variables related to the generation of a mobile document associated with a specific SMS message, and document template.
	Document struct {
		// The API template name related to the document template.
		Template string
		// The version of the document template to use. If not provided then the configured default (active) version is used.
		Version int
		// The authentication password for the document.
		// If not provided (Empty or Null) then the document will not require authentication.
		Password string
		// The values of the variables defined in the document template.
		Variables interface{}
	}
}

// BulkMessageRequest is a request to send a batch of SMS messages
type BulkMessageRequest struct {
	// A number of optional settings for a send
	SendOptions struct {
		// The sender ID (source address) to use when sending the messages specified. Maximum length is 11 characters.
		// Please note that a best effort is made to use the value specified but availability is ultimately
		// determined by the mobile network operator
		SenderId string
		// The type of duplication check to use (Defaults to 'None')
		DuplicateCheck string
		// The date and time in UTC when the messages should be sent.
		// Maximum date and time allowed is 3 months from the current date and time.
		// Value should be formatted according to ISO 8601 UTC format (yyyy-mm-ddThh:mm[:ss][Z])
		StartDeliveryUtc string
		// This value should be set in UTC if delivery should be staggered over a period of time.
		// Leave null if all messages should be sent as soon as possible.
		// Value, if set, must be greater than StartDeliveryUtc
		// and should be formatted according to ISO 8601 UTC format (yyyy-mm-ddThh:mm[:ss][Z])
		EndDeliveryUtc string
		// The name of the reply rule set that should be used for the send. Replies that are received as a result of
		// the send are then subject to rules configured in the reply rule set governing auto forwards and auto
		// responses
		ReplyRuleSetName string
		// A user defined value for the name of the campaign associated with this send. Only used by the system for
		// reporting purposes. Maximum value is 100 characters
		CampaignName string
		// A user defined value for the name of the cost center associated with this send. Only used by the system for
		// reporting purposes. Maximum value is 100 characters
		CostCentre string
		// True if SMS messages should not be sent to mobile subscribers that have opted-out as a result of a previous
		// send; otherwise false
		CheckOptOuts bool
		// True if this send should try shorten and track urls otherwise false.
		ShortenUrls bool
		// The amount of time in hours an SMS should remain valid.
		// The network will continue to try to deliver the SMS over the validity period.
		ValidityPeriod int
		// Should the messages be sent or is this a test of the API?
		// If TestMode = true then it will return the result but won’t send the data, and won’t appear in any sent report.
		TestMode bool
		// The api name of the reply rule to use for this send
		RuleName string
		// The specific reply rule version to use for this send (if null, the active version will be used).
		ReplyRuleVersion int
		// Any Extra reply forward emails for this send.
		ExtraForwardEmails string
	}
	Messages []Message
}

// BulkMessageRequest is a request to send a batch of SMS messages
type GroupMessageRequest struct {
	// A number of optional settings for a send
	SendOptions struct {
		// The sender ID (source address) to use when sending the messages specified. Maximum length is 11 characters.
		// Please note that a best effort is made to use the value specified but availability is ultimately
		// determined by the mobile network operator
		SenderId string
		// The type of duplication check to use (Defaults to 'None')
		DuplicateCheck string
		// The date and time in UTC when the messages should be sent.
		// Maximum date and time allowed is 3 months from the current date and time.
		// Value should be formatted according to ISO 8601 UTC format (yyyy-mm-ddThh:mm[:ss][Z])
		StartDeliveryUtc string
		// This value should be set in UTC if delivery should be staggered over a period of time.
		// Leave null if all messages should be sent as soon as possible.
		// Value, if set, must be greater than StartDeliveryUtc
		// and should be formatted according to ISO 8601 UTC format (yyyy-mm-ddThh:mm[:ss][Z])
		EndDeliveryUtc string
		// The name of the reply rule set that should be used for the send. Replies that are received as a result of
		// the send are then subject to rules configured in the reply rule set governing auto forwards and auto
		// responses
		ReplyRuleSetName string
		// A user defined value for the name of the campaign associated with this send. Only used by the system for
		// reporting purposes. Maximum value is 100 characters
		CampaignName string
		// A user defined value for the name of the cost center associated with this send. Only used by the system for
		// reporting purposes. Maximum value is 100 characters
		CostCentre string
		// True if SMS messages should not be sent to mobile subscribers that have opted-out as a result of a previous
		// send; otherwise false
		CheckOptOuts bool
		// True if this send should try shorten and track urls otherwise false.
		ShortenUrls bool
		// The amount of time in hours an SMS should remain valid.
		// The network will continue to try to deliver the SMS over the validity period.
		ValidityPeriod int
		// Should the messages be sent or is this a test of the API?
		// If TestMode = true then it will return the result but won’t send the data, and won’t appear in any sent report.
		TestMode bool
		// The api name of the reply rule to use for this send
		RuleName string
		// The specific reply rule version to use for this send (if null, the active version will be used).
		ReplyRuleVersion int
		// Any Extra reply forward emails for this send.
		ExtraForwardEmails string
	}
	Message Message
	// The names of the groups that contacts should belong to. These contacts will be the recipients of
	// the SMS message. Should contain at least one group name.
	Groups []string
}

// BulkMessageResponse is returned by the API in response to a BulkMessageRequest
type BulkMessageResponse struct {
	// The total cost of the send excluding any VAT, GST or tax.
	Cost int `json:"cost"`
	// The remaining balance after the Cost has been deducted for the send.
	RemainingBalance int `json:"remainingBalance"`
	// The system generated ID for the send or batch of messages
	EventID int `json:"eventId"`
	// A sample message that was generated for the send. Purely informational.
	Sample string `json:"sample"`
	// The break down of the costs for the send grouped by mobile network and feature.
	CostBreakdown struct {
		Quantity int    `json:"quantity"`
		Cost     int    `json:"cost"`
		Network  string `json:"network"`
	} `json:"costBreakdown"`
	// The total number of messages that were successfully enqueued for delivery. Long messages may require
	// multiple SMS messages or parts to send a message and this is reflected in the Parts value.
	Messages int `json:"messages"`
	// The total number of SMS messages that were successfully enqueued for delivery.
	Parts       int `json:"parts"`
	ErrorReport struct {
		NoNetwork  int `json:"noNetwork"`
		Duplicates int `json:"duplicates"`
		OptedOuts  int `json:"optedOuts"`
		Faults     struct {
			RawDestination      string `json:"rawDestination"`
			ScrubbedDestination string `json:"scrubbedDestination"`
			CustomerID          string `json:"customerId"`
			ErrorMessage        string `json:"errorMessage"`
			Status              string `json:"status"`
		} `json:"faults"`
	} `json:"errorReport"`
}

// GetBalance returns the current balance of the account.
// Post Paid customers will always receive a balance of 1000000,
// as no credit deduction occurs on this account type.
func (c *Client) GetBalance() (int, error) {
	resp := struct {
		Balance int `json:"balance"`
	}{}

	_, err := c.get("Balance", nil, &resp)
	if err != nil {
		return 0, err
	}

	return resp.Balance, nil
}

// SendBulkMessages sends SMS message(s) to a multiple recipients.
func (c *Client) SendBulkMessages(data BulkMessageRequest) (resp BulkMessageResponse, err error) {
	_, err = c.post("BulkMessages", data, &resp)
	return
}

// SendGroupMessages send SMS message(s) to one or more specified groups.
func (c *Client) SendGroupMessages(data GroupMessageRequest) (resp BulkMessageResponse, err error) {
	_, err = c.post("GroupMessages", data, &resp)
	return
}
