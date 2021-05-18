package service

import "log"

func TestAll() {
	log.Printf(">>>>>>>>>>start testing Success flow")
	TestSuccessFlow()
	log.Printf("\nstart testing Payer Unavailable flow")
	TestPayerUnavailableFlow()
	log.Printf("\nstart testing Barred Beneficiary flow")
	TestBarredBeneficiaryFlow()
	log.Printf("\nstart testing payer Limitation On Transaction flow")
	TestLimitationOnTransactionFlow()
}
