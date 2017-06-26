package mango

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// ErrPreAuthorizationFailed custom error returned in case of failed payIn.
type ErrPreAuthorizationFailed struct {
	preauthorizationID string
	msg                string
}

func (e *ErrPreAuthorizationFailed) Error() string {
	return fmt.Sprintf("preAuthorization %s failed: %s ", e.preauthorizationID, e.msg)
}

// PreAuthorization ...
type PreAuthorization struct {
	ProcessReply
	AuthorID              string
	DebitedFunds          Money
	PaymentStatus         string //WAITING, CANCELED, EXPIRED, VALIDATED
	ExecutionType         string
	SecureMode            string
	CardID                string
	SecureModeNeeded      bool
	SecureModeRedirectURL string
	SecureModeReturnURL   string
	PayInID               string

	service *MangoPay
}

func (p *PreAuthorization) String() string {
	return struct2string(p)
}

// NewPreAuthorization ...
func (m *MangoPay) NewPreAuthorization(author Consumer, amount Money, secureMode, cardID, returnURL string) (*PreAuthorization, error) {
	msg := "new preauthorization: "
	if author == nil {
		return nil, errors.New(msg + "nil author")
	}
	id := consumerId(author)
	if id == "" {
		return nil, errors.New(msg + "author has empty Id")
	}

	u, err := url.Parse(returnURL)
	if err != nil {
		return nil, errors.New(msg + err.Error())
	}
	p := &PreAuthorization{
		AuthorID:            id,
		DebitedFunds:        amount,
		SecureMode:          secureMode,
		CardID:              cardID,
		SecureModeReturnURL: u.String(),
		service:             m,
	}

	return p, nil
}

// Save ...
func (p *PreAuthorization) Save() error {
	data := JsonObject{}
	j, err := json.Marshal(p)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(j, &data); err != nil {
		return err
	}

	// Force float64 to int conversion after unmarshalling.
	for _, field := range []string{"CreationDate", "ExecutionDate"} {
		data[field] = int(data[field].(float64))
	}

	// Fields not allowed when creating a tranfer.
	for _, field := range []string{"Id", "Status", "PaymentStatus", "ResultCode", "ResultMessage", "ExecutionType", "SecureModeNeeded", "SecureModeRedirectUrl", "ExpirationDate", "PayInId"} {
		delete(data, field)
	}

	tr, err := p.service.anyRequest(new(PreAuthorization), actionCreatePreAuthorization, data)
	if err != nil {
		return err
	}
	serv := p.service
	*p = *(tr.(*PreAuthorization))
	p.service = serv

	if p.Status == "FAILED" {
		return &ErrPreAuthorizationFailed{p.Id, p.ResultMessage}
	}
	return nil
}
