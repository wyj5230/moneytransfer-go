package service

import (
	"encoding/json"
	"fmt"
	"time"
)

const CREDIT_PARTY_MSISDN_100 string = "263775892100"
const CREDIT_PARTY_MSISDN_117 string = "263775892117"
const CREDIT_PARTY_MSISDN_104 string = "263775892104"
const CREDIT_PARTY_MSISDN_111 string = "263775892111"

const STATUS_COMPLETED string = "70000"
const STATUS_DECLINED_PAYER_CURRENTLY_UNAVAILABLE string = "90400"
const STATUS_DECLINED_BARRED_BENEFICIARY string = "90201"
const STATUS_LIMITATIONS_ON_TRANSACTION_VALUE string = "90305"

func getPayers() []Payer {
	var payersResponse = HttpGet("https://api-mt.pre.thunes.com/v2/money-transfer/payers",
		GetApiKey(), GetApiSecret())
	var payers []Payer
	err := json.Unmarshal(payersResponse, &payers)
	if err != nil {
		fmt.Println("GetTransaction err:", err)
	}
	for i := 0; i < len(payers); i++ {
		fmt.Println("payer： " + payers[i].Name + ", currency:" + payers[i].Currency + ", id:" + FloatToString(payers[i].Id))
	}
	fmt.Println()
	return payers
}

func getTransaction(transactionId string) TransactionResponse {
	var transactionResponseByte = HttpGet("https://api-mt.pre.thunes.com/v2/money-transfer/transactions/"+transactionId,
		GetApiKey(), GetApiSecret())
	var transactionResponse TransactionResponse
	err := json.Unmarshal(transactionResponseByte, &transactionResponse)
	if err != nil {
		fmt.Println("GetTransaction err:", err)
	}
	fmt.Printf(">>>>>>>>>>Transaction %d retrieved, current status: %s \n\n",
		transactionResponse.Id, transactionResponse.StatusMessage)
	return transactionResponse
}

func GetPayerLimitation(payerId string) C2c {
	var payerResponseByte = HttpGet("https://api-mt.pre.thunes.com/v2/money-transfer/payers/"+payerId,
		GetApiKey(), GetApiSecret())
	var payer Payer
	err := json.Unmarshal(payerResponseByte, &payer)
	if err != nil {
		fmt.Println("getPayerLimitation err:", err)
	}
	return payer.TransactionType.C2c
}

func createQuotation() QuotationResponse {
	quotationRequest := GetQuotationRequestBody()
	var quotationResponseByte = HttpPost("https://api-mt.pre.thunes.com/v2/money-transfer/quotations",
		quotationRequest, GetApiKey(), GetApiSecret())
	var quotationResponse QuotationResponse
	err := json.Unmarshal(quotationResponseByte, &quotationResponse)
	if err != nil {
		fmt.Println("CreateQuotation err:", err)
	}
	fmt.Printf(">>>>>>>>>>Quotation created: %d \n", quotationResponse.Id)
	fmt.Printf("Please confirm you are paying: %g %s through payer: %s \n",
		quotationResponse.SentAmount.Amount, quotationResponse.SentAmount.Currency, quotationResponse.Payer.Name)
	fmt.Printf("The beneficiary will received: %g %s and the fee will be: %g %s \n\n",
		quotationResponse.Destination.Amount, quotationResponse.Destination.Currency,
		quotationResponse.Fee.Amount, quotationResponse.Fee.Currency)
	return quotationResponse
}

func createTransaction(quotationId string, msisdn string) TransactionResponse {
	transactionRequest := GetCreateTransactionRequestBody(msisdn)
	var transactionResponseByte = HttpPost("https://api-mt.pre.thunes.com/v2/money-transfer/quotations/"+
		quotationId+"/transactions",
		transactionRequest, GetApiKey(), GetApiSecret())
	var transactionResponse TransactionResponse
	err := json.Unmarshal(transactionResponseByte, &transactionResponse)
	if err != nil {
		fmt.Println("CreateTransaction err:", err)
	}
	fmt.Printf(">>>>>>>>>>Transaction %d created, current status %s \n",
		transactionResponse.Id, transactionResponse.StatusMessage)
	SendTransactionStatusUpdateEmail(transactionResponse)
	return transactionResponse
}

func confirmTransaction(transactionId string) TransactionResponse {
	var transactionResponseByte = HttpPost("https://api-mt.pre.thunes.com/v2/money-transfer/transactions/"+
		transactionId+"/confirm",
		nil, GetApiKey(), GetApiSecret())
	var transactionResponse TransactionResponse
	err := json.Unmarshal(transactionResponseByte, &transactionResponse)
	if err != nil {
		fmt.Println("ConfirmTransaction err:", err)
	}
	fmt.Printf(">>>>>>>>>>Transaction %d confirmed, current status: %s \n",
		transactionResponse.Id, transactionResponse.StatusMessage)
	SendTransactionStatusUpdateEmail(transactionResponse)
	return transactionResponse
}

func HandlerCallback(transactionResponseByte []byte) {
	var transactionResponse TransactionResponse
	err := json.Unmarshal(transactionResponseByte, &transactionResponse)
	if err != nil {
		fmt.Println("HandlerCallback err:", err)
	}
	SendTransactionStatusUpdateEmail(transactionResponse)
}

func SendTransactionStatusUpdateEmail(transactionResponse TransactionResponse) {
	subject := "Transaction update :" + IntToString(transactionResponse.Id)
	content := "Status for transaction id:" + IntToString(transactionResponse.Id) +
		" is updated to " + transactionResponse.StatusMessage + "."
	switch transactionResponse.Status {
	case STATUS_COMPLETED:
		content += "\ntransaction is successful."
	case STATUS_DECLINED_PAYER_CURRENTLY_UNAVAILABLE:
		content += "\npayer is currently unavailable, please try again later."
	case STATUS_DECLINED_BARRED_BENEFICIARY:
		content += "\nbeneficiary is barred, we are sorry for the inconvenience caused."
	case STATUS_LIMITATIONS_ON_TRANSACTION_VALUE:
		c2c := GetPayerLimitation("37")
		content += "\ntransaction amount exceeds the limitation, payer's transaction maximum is: " +
			IntToString(c2c.MaximumTransactionAmount) + " and minimum is: " + IntToString(c2c.MinimumTransactionAmount)
	}
	recipient := transactionResponse.Sender.Email
	sendEmail(subject, content, recipient)
}

func TestSuccessFlow() {
	getPayers()
	quotationResponse := createQuotation()
	transactionResponse := createTransaction(IntToString(quotationResponse.Id), CREDIT_PARTY_MSISDN_100)
	confirmTransaction(IntToString(transactionResponse.Id))
	time.Sleep(time.Second * 3)
	getTransaction(IntToString(transactionResponse.Id))
}

func TestPayerUnavailableFlow() {
	quotationResponse := createQuotation()
	transactionResponse := createTransaction(IntToString(quotationResponse.Id), CREDIT_PARTY_MSISDN_117)
	confirmTransaction(IntToString(transactionResponse.Id))
	time.Sleep(time.Second * 3)
	getTransaction(IntToString(transactionResponse.Id))
}

func TestBarredBeneficiaryFlow() {
	quotationResponse := createQuotation()
	transactionResponse := createTransaction(IntToString(quotationResponse.Id), CREDIT_PARTY_MSISDN_104)
	confirmTransaction(IntToString(transactionResponse.Id))
	time.Sleep(time.Second * 3)
	getTransaction(IntToString(transactionResponse.Id))
}

func TestLimitationOnTransactionFlow() {
	quotationResponse := createQuotation()
	transactionResponse := createTransaction(IntToString(quotationResponse.Id), CREDIT_PARTY_MSISDN_111)
	confirmTransaction(IntToString(transactionResponse.Id))
	time.Sleep(time.Second * 3)
	getTransaction(IntToString(transactionResponse.Id))
}
