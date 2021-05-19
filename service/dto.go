package service

type QuotationRequest struct {
	ExternalId      string      `json:"external_id"`
	PayerId         string      `json:"payer_id"`
	Mode            string      `json:"mode"`
	TransactionType string      `json:"transaction_type"`
	Source          Source      `json:"source"`
	Destination     Destination `json:"destination"`
}

type QuotationResponse struct {
	CreationDate    string      `json:"creation_date"`
	Destination     Destination `json:"destination"`
	ExpirationDate  string      `json:"expiration_date"`
	Fee             Fee         `json:"fee"`
	Id              int         `json:"id"`
	Mode            string      `json:"mode"`
	Payer           Payer       `json:"payer"`
	SentAmount      SentAmount  `json:"sent_amount"`
	Source          Source      `json:"source"`
	TransactionType string      `json:"transaction_type"`
	WholesaleFxRate float32     `json:"wholesale_fx_rate"`
}

type TransactionRequest struct {
	CreditPartyIdentifier CreditPartyIdentifier `json:"credit_party_identifier"`
	Sender                Sender                `json:"sender"`
	Beneficiary           Beneficiary           `json:"beneficiary"`
	ExternalId            string                `json:"external_id"`
}

type TransactionResponse struct {
	Id                        int                   `json:"id"`
	Status                    string                `json:"status"`
	StatusMessage             string                `json:"status_message"`
	StatusClass               string                `json:"status_class"`
	StatusClassMessage        string                `json:"status_class_message"`
	ExternalId                string                `json:"external_id"`
	ExternalCode              string                `json:"external_code"`
	TransactionType           string                `json:"transaction_type"`
	PayerTransactionReference string                `json:"payer_transaction_reference"`
	PayerTransactionCode      string                `json:"payer_transaction_code"`
	CreationDate              string                `json:"creation_date"`
	ExpirationDate            string                `json:"expiration_date"`
	CreditPartyIdentifier     CreditPartyIdentifier `json:"credit_party_identifier"`
	Source                    Source                `json:"source"`
	Destination               Destination           `json:"destination"`
	Payer                     Payer                 `json:"payer"`
	Sender                    Sender                `json:"sender"`
	SentAmount                SentAmount            `json:"sent_amount"`
	WholesaleFxRate           float32               `json:"wholesale_fx_rate"`
	RetailRate                string                `json:"retail_rate"`
	RetailFee                 string                `json:"retail_fee"`
	RetailFeeCurrency         string                `json:"retail_fee_currency"`
	Fee                       Fee                   `json:"fee"`
	PurposeOfRemittance       string                `json:"purpose_of_remittance"`
	DocumentReferenceNumber   string                `json:"document_reference_number"`
}

type AttachmentResponse struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	ContentType   string `json:"content_type"`
	TransactionId int    `json:"transaction_id"`
	Type          string `json:"type"`
}

type Source struct {
	Amount         float32 `json:"amount"`
	Currency       string  `json:"currency"`
	CountryIsoCode string  `json:"country_iso_code"`
}

type Destination struct {
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
}

type Fee struct {
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
}

type Payer struct {
	CountryIsoCode  string          `json:"country_iso_code"`
	Currency        string          `json:"currency"`
	Id              float32         `json:"id"`
	Name            string          `json:"name"`
	TransactionType TransactionType `json:"transaction_types"`
}

type SentAmount struct {
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
}

type CreditPartyIdentifier struct {
	Msisdn            string `json:"msisdn"`
	BankAccountNumber string `json:"bank_account_number"`
	SwiftBicCode      string `json:"swift_bic_code"`
}

type Sender struct {
	Lastname                  string `json:"lastname"`
	Firstname                 string `json:"firstname"`
	NationalityCountryIsoCode string `json:"nationality_country_iso_code"`
	DateOfBirth               string `json:"date_of_birth"`
	CountryOfBirthIsoCode     string `json:"country_of_birth_iso_code"`
	Gender                    string `json:"gender"`
	Address                   string `json:"address"`
	PostalCode                string `json:"postal_code"`
	City                      string `json:"city"`
	CountryIsoCode            string `json:"country_iso_code"`
	Msisdn                    string `json:"msisdn"`
	Email                     string `json:"email"`
	IdType                    string `json:"id_type"`
	IdNumber                  string `json:"id_number"`
	IdDeliveryDate            string `json:"id_delivery_date"`
	Occupation                string `json:"occupation"`
}

type Beneficiary struct {
	Lastname                  string `json:"lastname"`
	Firstname                 string `json:"firstname"`
	NationalityCountryIsoCode string `json:"nationality_country_iso_code"`
	DateOfBirth               string `json:"date_of_birth"`
	CountryOfBirthIsoCode     string `json:"country_of_birth_iso_code"`
	Gender                    string `json:"gender"`
	Address                   string `json:"address"`
	PostalCode                string `json:"postal_code"`
	City                      string `json:"city"`
	CountryIsoCode            string `json:"country_iso_code"`
	Msisdn                    string `json:"msisdn"`
	Email                     string `json:"email"`
	IdType                    string `json:"id_type"`
	IdCountryIsoCode          string `json:"id_country_iso_code"`
	IdNumber                  string `json:"id_number"`
	Occupation                string `json:"occupation"`
}

type TransactionType struct {
	C2c C2c `json:"C2C"`
}

type C2c struct {
	MaximumTransactionAmount int `json:"maximum_transaction_amount"`
	MinimumTransactionAmount int `json:"minimum_transaction_amount"`
}
