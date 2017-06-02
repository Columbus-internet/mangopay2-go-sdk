package mango

import (
	"encoding/json"
	"errors"
	"log"
)

// KYCDocumentType represents type of KYC document
type KYCDocumentType int

// KYCDocumentType ...
const (
	IdentityProof KYCDocumentType = iota
	RegistrationProof
	ArticlesOfAssociation
	ShareholderDeclaration
	AddressProof
)

var kycDocumentTypes = map[KYCDocumentType]string{
	IdentityProof:          "IDENTITY_PROOF",
	RegistrationProof:      "REGISTRATION_PROOF",
	ArticlesOfAssociation:  "ARTICLES_OF_ASSOCIATION",
	ShareholderDeclaration: "SHAREHOLDER_DECLARATION",
	AddressProof:           "ADDRESS_PROOF",
}

type KYCDocumentStatus int

const (
	CREATED KYCDocumentStatus = iota
	VALIDATION_ASKED
	VALIDATED
	REFUSED
)

var KYCDocumentStatuses = map[KYCDocumentStatus]string{
	CREATED:          "CREATED",
	VALIDATION_ASKED: "VALIDATION_ASKED",
	VALIDATED:        "VALIDATED",
	REFUSED:          "REFUSED",
}

// KYCDocument ...
type KYCDocument struct {
	Id                   string
	CreationDate         int64
	Tag                  string
	UserId               string
	Status               string
	RefusedReasonMessage string
	RefusedReasonType    string

	service      *MangoPay
	documentType KYCDocumentType
}

// KYCDocumentList ...
type KYCDocumentList []*KYCDocument

// NewKYCDocument creates KYC document object for user
func (m *MangoPay) NewKYCDocument(user Consumer, kycDocumentType KYCDocumentType) (*KYCDocument, error) {
	id := consumerId(user)
	if id == "" {
		return nil, errors.New("unable to create KYC document: empty user ID")
	}
	d := &KYCDocument{
		UserId:       id,
		documentType: kycDocumentType,
	}
	d.service = m

	data := JsonObject{}
	j, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(j, &data); err != nil {
		return nil, err
	}

	for _, field := range []string{"Id", "CreationDate", "Status", "RefusedReasonMessage", "RefusedReasonType"} {
		delete(data, field)
	}

	data["Type"] = kycDocumentTypes[d.documentType]

	doc, err := d.service.anyRequest(new(KYCDocument), actionCreateUserKYCDocument, data)
	if err != nil {
		log.Println(data)
		return nil, err
	}
	serv := d.service
	*d = *(doc.(*KYCDocument))
	d.service = serv
	return d, nil
}

// KYCDocuments ...
func (m *MangoPay) KYCDocuments(user Consumer) (KYCDocumentList, error) {
	id := consumerId(user)
	if id == "" {
		return nil, errors.New("unable to get KYC documents list: empty user ID")
	}

	docList, err := m.anyRequest(new(KYCDocumentList), actionListUserKYCDocuments, JsonObject{"UserId": id})
	if err != nil {
		return nil, err
	}
	return *(docList.(*KYCDocumentList)), nil
}

// KYCDocument ...
func (m *MangoPay) KYCDocument(id string) (*KYCDocument, error) {
	d, err := m.anyRequest(new(KYCDocument), actionViewKYCDocument, JsonObject{"Id": id})
	if err != nil {
		return nil, err
	}
	return d.(*KYCDocument), nil
}

// AddPage ...
func (d *KYCDocument) AddPage(pageData string) error {
	data := JsonObject{}
	j, err := json.Marshal(d)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(j, &data); err != nil {
		return err
	}

	for _, field := range []string{"CreationDate", "Status", "RefusedReasonMessage", "RefusedReasonType"} {
		delete(data, field)
	}

	data["File"] = pageData

	_, err = d.service.anyRequest(new(interface{}), actionCreateUserKYCDocumentPage, data)
	if err != nil {
		return err
	}
	return nil
}

// Submit ...
func (d *KYCDocument) Submit(status KYCDocumentStatus) error {

	d.Status = KYCDocumentStatuses[status]

	data := JsonObject{}
	j, err := json.Marshal(d)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(j, &data); err != nil {
		return err
	}

	for _, field := range []string{"CreationDate", "RefusedReasonMessage", "RefusedReasonType"} {
		delete(data, field)
	}

	_, err = d.service.anyRequest(new(interface{}), actionSubmitUserKYCDocument, data)
	if err != nil {
		return err
	}
	return nil
}
