package service

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
)

const CREDIT_PARTY_MSISDN_100 string = "263775892100"
const CREDIT_PARTY_MSISDN_117 string = "263775892117"
const CREDIT_PARTY_MSISDN_104 string = "263775892104"
const CREDIT_PARTY_MSISDN_111 string = "263775892111"

func TestHappyFlow(t *testing.T) {
	getPayers()
	quotationResponse := createQuotation(getQuotationRequestTestData())
	transactionResponse := createTransaction(IntToString(quotationResponse.Id), getCreateTransactionRequestTestData(CREDIT_PARTY_MSISDN_100))
	attachmentResponse := attachDocumentToTransaction(IntToString(transactionResponse.Id))
	if attachmentResponse.TransactionId == transactionResponse.Id {
		fmt.Println("Happy flow attachment uploaded successfully")
		fmt.Println()
	}
	confirmTransaction(IntToString(transactionResponse.Id))
	time.Sleep(time.Second * 3)
	result := getTransaction(IntToString(transactionResponse.Id))
	if result.Status == STATUS_COMPLETED {
		fmt.Println("Happy flow result is ok")
	}
}

func TestPayerUnavailableFlow(t *testing.T) {
	quotationResponse := createQuotation(getQuotationRequestTestData())
	transactionResponse := createTransaction(IntToString(quotationResponse.Id), getCreateTransactionRequestTestData(CREDIT_PARTY_MSISDN_117))
	confirmTransaction(IntToString(transactionResponse.Id))
	time.Sleep(time.Second * 3)
	result := getTransaction(IntToString(transactionResponse.Id))
	if result.Status == STATUS_DECLINED_PAYER_CURRENTLY_UNAVAILABLE {
		fmt.Println("Payer Unavailable flow result is ok")
	}
}

func TestBarredBeneficiaryFlow(t *testing.T) {
	quotationResponse := createQuotation(getQuotationRequestTestData())
	transactionResponse := createTransaction(IntToString(quotationResponse.Id), getCreateTransactionRequestTestData(CREDIT_PARTY_MSISDN_104))
	confirmTransaction(IntToString(transactionResponse.Id))
	time.Sleep(time.Second * 3)
	result := getTransaction(IntToString(transactionResponse.Id))
	if result.Status == STATUS_DECLINED_BARRED_BENEFICIARY {
		fmt.Println("Barred Beneficiary flow result is ok")
	}
}

func TestLimitationOnTransactionFlow(t *testing.T) {
	quotationResponse := createQuotation(getQuotationRequestTestData())
	transactionResponse := createTransaction(IntToString(quotationResponse.Id), getCreateTransactionRequestTestData(CREDIT_PARTY_MSISDN_111))
	confirmTransaction(IntToString(transactionResponse.Id))
	time.Sleep(time.Second * 3)
	result := getTransaction(IntToString(transactionResponse.Id))
	if result.Status == STATUS_LIMITATIONS_ON_TRANSACTION_VALUE {
		fmt.Println("Limitation On Transaction flow result is ok")
	}
}

func TestAll(t *testing.T) {
	log.Printf(">>>>>>>>>>start testing Success flow")
	TestHappyFlow(t)
	log.Printf("\nstart testing Payer Unavailable flow")
	TestPayerUnavailableFlow(t)
	log.Printf("\nstart testing Barred Beneficiary flow")
	TestBarredBeneficiaryFlow(t)
	log.Printf("\nstart testing payer Limitation On Transaction flow")
	TestLimitationOnTransactionFlow(t)
}

func getQuotationRequestTestData() []byte {
	source := Source{10, "SGD", "SGP"}
	destination := Destination{1, "PHP"}
	quotation := QuotationRequest{GetExternalId(), "83", "SOURCE_AMOUNT", "C2C", source, destination}
	quotationJson, err := json.Marshal(quotation)
	if err != nil {
		fmt.Println("GetQuotationRequestBody err:", err)
	}
	return quotationJson
}

func getCreateTransactionRequestTestData(creditPartyMsisdn string) []byte {
	creditPartyIdentifier := CreditPartyIdentifier{creditPartyMsisdn, "0123456789", "ABCDEFGH"}
	sender := Sender{"Doe", "John", "SGP", "1970-01-01",
		"SGP", "MALE", "42 Rue des fleurs", "75000", "Paris",
		"FRA", "33712345678", "327113606@qq.com", "SOCIAL_SECURITY",
		"502-42-0158", "2016-01-01", "Residential Advisor"}
	beneficiary := Beneficiary{"Doe", "John", "SGP", "1970-01-01",
		"SGP", "MALE", "42 Rue des fleurs", "75000", "Paris",
		"FRA", "33712345678", "327113606@qq.com", "SOCIAL_SECURITY",
		"FRA", "502-42-0158", "Residential Advisor"}
	transactionRequest := TransactionRequest{creditPartyIdentifier, sender, beneficiary, GetExternalId()}
	transactionRequestJson, err := json.Marshal(transactionRequest)
	if err != nil {
		fmt.Println("GetCreateTransactionRequestBody err:", err)
	}
	return transactionRequestJson
}

func attachDocumentToTransaction(transactionId string) AttachmentResponse {
	return AttachmentToTransaction(transactionId, "invoice", "Thunes.pdf",
		"C:\\Users\\Administrator\\Desktop\\Thunes Demo\\Thunes.pdf")
}
