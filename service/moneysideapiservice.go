package service

import (
	"encoding/json"
	"fmt"
)

const StatusCompleted string = "70000"
const StatusDeclinedPayerCurrentlyUnavailable string = "90400"
const StatusDeclinedBarredBeneficiary string = "90201"
const StatusLimitationsOnTransactionValue string = "90305"

func getPayers() []Payer {
	var payersResponse = HttpGet("https://api-mt.pre.thunes.com/v2/money-transfer/payers",
		GetApiKey(), GetApiSecret())
	var payers []Payer
	err := json.Unmarshal(payersResponse, &payers)
	if err != nil {
		//fmt.Println("GetPayers err:", err)
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

func getPayerLimitation(payerId string) C2c {
	var payerResponseByte = HttpGet("https://api-mt.pre.thunes.com/v2/money-transfer/payers/"+payerId,
		GetApiKey(), GetApiSecret())
	var payer Payer
	err := json.Unmarshal(payerResponseByte, &payer)
	if err != nil {
		fmt.Println("getPayerLimitation err:", err)
	}
	return payer.TransactionType.C2c
}

func CreateQuotation(quotationRequest QuotationRequest) QuotationResponse {
	//
	quotationRequest.ExternalId = GetExternalId()
	quotationRequestByte, quotationRequestErr := json.Marshal(quotationRequest)
	if quotationRequestErr != nil {
		fmt.Println("quotationRequest body err:", quotationRequestErr)
	}
	var quotationResponseByte = HttpPost("https://api-mt.pre.thunes.com/v2/money-transfer/quotations",
		quotationRequestByte, GetApiKey(), GetApiSecret())
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

func createTransaction(quotationId string, transactionRequest TransactionRequest) TransactionResponse {
	transactionRequest.ExternalId = GetExternalId()
	transactionRequestByte, transactionRequestErr := json.Marshal(transactionRequest)
	if transactionRequestErr != nil {
		fmt.Println("quotationRequest body err:", transactionRequestErr)
	}
	var transactionResponseByte = HttpPost("https://api-mt.pre.thunes.com/v2/money-transfer/quotations/"+
		quotationId+"/transactions",
		transactionRequestByte, GetApiKey(), GetApiSecret())
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

func AttachmentToTransaction(transactionId string, fileType string, fileName string, filePath string) AttachmentResponse {
	var attachmentResponseByte = attachmentUploadRequest("https://api-mt.pre.thunes.com/v2/money-transfer/transactions/"+
		transactionId+"/attachments", fileType, fileName, filePath, GetApiKey(), GetApiSecret())
	var attachmentResponse AttachmentResponse
	err := json.Unmarshal(attachmentResponseByte, &attachmentResponse)
	if err != nil {
		fmt.Println("ConfirmTransaction err:", err)
	}
	fmt.Printf(">>>>>>>>>>Attachment %s is attached to transaction:  %s \n",
		fileName, transactionId)
	return attachmentResponse
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
		" is updated to " + transactionResponse.StatusMessage + ".\n transaction detail: " +
		FloatToString(transactionResponse.Destination.Amount) + transactionResponse.Destination.Currency +
		". Creation date: " + transactionResponse.CreationDate + ".\n"
	switch transactionResponse.Status {
	case StatusCompleted:
		content += "\ntransaction is successful."
	case StatusDeclinedPayerCurrentlyUnavailable:
		content += "\npayer is currently unavailable, please try again later."
	case StatusDeclinedBarredBeneficiary:
		content += "\nbeneficiary is barred, we are sorry for the inconvenience caused."
	case StatusLimitationsOnTransactionValue:
		c2c := getPayerLimitation("37")
		content += "\ntransaction amount exceeds the limitation, payer's transaction maximum is: " +
			IntToString(c2c.MaximumTransactionAmount) + " and minimum is: " + IntToString(c2c.MinimumTransactionAmount)
	}
	recipient := transactionResponse.Sender.Email
	sendEmail(subject, content, recipient)
}
