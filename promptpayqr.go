package thaiqr

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	PayloadFormatEMVQRCPSMerchantPresentedMode = "01"
	POIMethodStatic                            = "11"
	POIMethodDynamic                           = "12"

	// BOTIDCreditTransferAID Credit Transfer Tag 29
	BOTIDCreditTransferAID   = "00"
	BOTIDMerchantMSISDN      = "01"
	BOTIDMerchantNationalID  = "02"
	BOTIDMerchantEWalletID   = "03"
	BOTIDMerchantBankAccount = "04"
	BOTIDMerchantOTA         = "05"

	// BOTIDBillPaymentAID Bill Payment Tag 30
	BOTIDBillPaymentAID      = "00"
	BOTIDBillPaymentBillerID = "01"
	BOTIDBillPaymentRef1     = "02"
	BOTIDBillPaymentRef2     = "03"

	BOTIDTag62TerminalID     = "07"
	GUIDPromptPay            = "A000000677010111"
	GUIDPromptPayBillPayment = "A000000677010112"
	TransactionCurrencyTHB   = "764"
	CountryCodeTH            = "TH"

	ProxyTypeEWalletID   = "EWALLETID"
	ProxyTypeNatID       = "NATID"
	ProxyTypeBankAccount = "BANKACCOUNT"
	ProxyTypeMsisdn      = "MSISDN"
)

const (
	IDPayloadFormat                     = "00"
	IDPOIMethod                         = "01"
	IDMerchantInformationBOT            = "29"
	IDMerchantInformationBOTBillPayment = "30"
	IDMerchantCategoryCode              = "52"
	IDTransactionCurrency               = "53"
	IDTransactionAmount                 = "54"
	IDCountryCode                       = "58"
	IDMerchantName                      = "59"
	IDMerchantCity                      = "60"
	IDPostalCode                        = "61"
	IDAdditionalFields                  = "62"
	IDCRC                               = "63"
)

type PromptPayQRCmd struct {
	ProxyID      string `json:"proxyId"`
	ProxyType    string `json:"proxyType"`
	Amount       string `json:"Amount"`
	OTA          string `json:"ota"`
	CountryCode  string `json:"countryCode"`
	CurrencyCode string `json:"currencyCode"`
}

type PromptPayBillPaymentQRCmd struct {
	BillerID     string `json:"billerId"`
	Ref1         string `json:"ref1"`
	Ref2         string `json:"ref2"`
	TerminalID   string `json:"terminalId"`
	Amount       string `json:"Amount"`
	CountryCode  string `json:"countryCode"`
	CurrencyCode string `json:"currencyCode"`
}

type PromptPayQRResults struct {
	PayloadFormatIndicator  string            `json:"payloadFormatIndicator"`
	PointOfInitiationMethod string            `json:"pointOfInitiationMethod"`
	CreditTransfer          *CreditTransfer   `json:"creditTransfer,omitempty"`
	BillPayment             *BillPayment      `json:"billPayment,omitempty"`
	MerchantCategoryCode    string            `json:"merchantCategoryCode,omitempty"`
	TransactionCurrency     string            `json:"transactionCurrency"`
	TransactionCurrencyCode string            `json:"transactionCurrencyCode"`
	TransactionAmount       string            `json:"transactionAmount"`
	CountryCode             string            `json:"countryCode"`
	MerchantName            string            `json:"merchantName,omitempty"`
	MerchantCity            string            `json:"merchantCity,omitempty"`
	PostalCode              string            `json:"postalCode,omitempty"`
	AdditionalFields        *AdditionalFields `json:"additionalFields,omitempty"`
	CRC                     string            `json:"crc"`
	Segments                *[]Segment        `json:"segments,omitempty"`
}

type CreditTransfer struct {
	AID         string     `json:"aid,omitempty"`
	MSISDN      string     `json:"msisdn,omitempty"`
	NationalID  string     `json:"nationalId,omitempty"`
	EWalletID   string     `json:"eWalletID,omitempty"`
	BankAccount string     `json:"bankAccount,omitempty"`
	OTA         string     `json:"ota,omitempty"`
	Segments    *[]Segment `json:"segments,omitempty"`
}

type BillPayment struct {
	AID        string     `json:"aid,omitempty"`
	BillerID   string     `json:"billerId,omitempty"`
	Reference1 string     `json:"reference1,omitempty"`
	Reference2 string     `json:"reference2,omitempty"`
	Segments   *[]Segment `json:"segments,omitempty"`
}

type AdditionalFields struct {
	TerminalID string     `json:"terminalId,omitempty"`
	Segments   *[]Segment `json:"segments,omitempty"`
}

// PromptPayQR represents a PromptPay QR code generator.
type PromptPayQR struct{}

// NewPromptPayQR returns a new PromptPayQR instance.
func NewPromptPayQR() *PromptPayQR {
	return &PromptPayQR{}
}

// GeneratePayload generates a PromptPay QR code payload.
func (qr *PromptPayQR) GeneratePayload(cmd PromptPayQRCmd) (string, error) {
	proxyID := sanitizeTarget(cmd.ProxyID)
	tagProxyType := determineTargetType(cmd.ProxyType)

	merchantInfoData := []string{
		formatField(BOTIDCreditTransferAID, GUIDPromptPay),
		formatField(tagProxyType, formatTarget(proxyID)),
	}

	if strings.TrimSpace(cmd.OTA) != "" {
		merchantInfoData = append(merchantInfoData, formatField(BOTIDMerchantOTA, cmd.OTA))
	}

	amount := strings.TrimSpace(cmd.Amount)

	currencyNo := TransactionCurrencyTHB
	if cmd.CurrencyCode != "" && currencyCode[cmd.CurrencyCode] != "" {
		currencyNo = currencyCode[cmd.CurrencyCode]
	}

	data := []string{
		formatField(IDPayloadFormat, PayloadFormatEMVQRCPSMerchantPresentedMode),
		formatField(IDPOIMethod, ifThenElse(amount != "", POIMethodDynamic, POIMethodStatic).(string)),
		formatField(IDMerchantInformationBOT, serialize(merchantInfoData)),
		formatField(IDTransactionCurrency, currencyNo),
		formatField(IDCountryCode, ifThenElse(cmd.CountryCode != "", cmd.CountryCode, CountryCodeTH).(string)),
	}
	if amount != "" {
		amountFormat, err := formatAmount(amount)
		if err != nil {
			return "", err
		}
		data = append(data, formatField(IDTransactionAmount, amountFormat))
	}

	dataToCrc := fmt.Sprintf("%s%s%s", serialize(data), IDCRC, "04")
	data = append(data, formatField(IDCRC, checksum([]byte(dataToCrc))))

	return serialize(data), nil
}

// GenerateBillPaymentPayload generates a PromptPay bill payment QR code payload.
func (qr *PromptPayQR) GenerateBillPaymentPayload(cmd PromptPayBillPaymentQRCmd) (string, error) {
	billerID := sanitizeTarget(cmd.BillerID)
	amount := strings.TrimSpace(cmd.Amount)

	data := []string{
		formatField(IDPayloadFormat, PayloadFormatEMVQRCPSMerchantPresentedMode),
		formatField(IDPOIMethod, ifThenElse(amount != "", POIMethodDynamic, POIMethodStatic).(string)),
		formatField(IDMerchantInformationBOTBillPayment, serialize([]string{
			formatField(BOTIDBillPaymentAID, GUIDPromptPayBillPayment),
			formatField(BOTIDBillPaymentBillerID, billerID),
			formatField(BOTIDBillPaymentRef1, cmd.Ref1),
			formatField(BOTIDBillPaymentRef2, cmd.Ref2),
		})),
	}

	currencyNo := TransactionCurrencyTHB
	if cmd.CurrencyCode != "" && currencyCode[cmd.CurrencyCode] != "" {
		currencyNo = currencyCode[cmd.CurrencyCode]
	}
	data = append(data, formatField(IDTransactionCurrency, currencyNo))

	if amount != "" {
		amountFormat, err := formatAmount(amount)
		if err != nil {
			return "", err
		}
		data = append(data, formatField(IDTransactionAmount, amountFormat))
	}

	data = append(data, formatField(IDCountryCode, ifThenElse(cmd.CountryCode != "", cmd.CountryCode, CountryCodeTH).(string)))

	if strings.TrimSpace(cmd.TerminalID) != "" {
		data = append(data,
			formatField(IDAdditionalFields, serialize([]string{
				formatField(BOTIDTag62TerminalID, cmd.TerminalID),
			})),
		)
	}

	dataToCrc := fmt.Sprintf("%s%s%s", serialize(data), IDCRC, "04")
	data = append(data, formatField(IDCRC, checksum([]byte(dataToCrc))))

	return serialize(data), nil
}

// determineTargetType determines the type of the target based on its length.
func determineTargetType(proxyType string) string {
	switch targetType := strings.ToUpper(proxyType); {
	case targetType == ProxyTypeEWalletID:
		return BOTIDMerchantEWalletID
	case targetType == ProxyTypeNatID:
		return BOTIDMerchantNationalID
	case targetType == ProxyTypeBankAccount:
		return BOTIDMerchantBankAccount
	default:
		return BOTIDMerchantMSISDN
	}
}

func (qr *PromptPayQR) Reader(data string) (*PromptPayQRResults, error) {
	if data == "" {
		return nil, invalidFormat()
	}

	if !VerifyPayloadChecksum(data) {
		return nil, errors.New("invalid checksum")
	}

	qrFields, qrSegments, err := deserialize(data)
	if err != nil {
		return nil, err
	}

	payloadFormatIndicator := qrFields[IDPayloadFormat]
	if payloadFormatIndicator != "01" {
		return nil, invalidFormat()
	}
	poiMethod := qrFields[IDPOIMethod]
	if !slices.Contains([]string{POIMethodStatic, POIMethodDynamic}, poiMethod) {
		return nil, invalidFormat()
	}
	creditTransferData := qrFields[IDMerchantInformationBOT]
	billPaymentData := qrFields[IDMerchantInformationBOTBillPayment]

	creditTransferFields, creditTransferSegments, _ := deserialize(creditTransferData)
	billPaymentFields, billPaymentSegments, _ := deserialize(billPaymentData)

	merchantCategoryCode := qrFields[IDMerchantCategoryCode]
	transactionCurrency := qrFields[IDTransactionCurrency]
	transactionAmount := qrFields[IDTransactionAmount]
	merchantName := qrFields[IDMerchantName]
	merchantCity := qrFields[IDMerchantCity]
	postalCode := qrFields[IDPostalCode]
	countryCode := qrFields[IDCountryCode]
	if len(countryCode) != 2 {
		return nil, invalidFormat()
	}
	additionalData := qrFields[IDAdditionalFields]
	additionalFields, additionalSegments, _ := deserialize(additionalData)

	return &PromptPayQRResults{
		PayloadFormatIndicator:  payloadFormatIndicator,
		PointOfInitiationMethod: poiMethod,
		CreditTransfer: &CreditTransfer{
			AID:         creditTransferFields[BOTIDCreditTransferAID],
			MSISDN:      creditTransferFields[BOTIDMerchantMSISDN],
			NationalID:  creditTransferFields[BOTIDMerchantNationalID],
			EWalletID:   creditTransferFields[BOTIDMerchantEWalletID],
			BankAccount: creditTransferFields[BOTIDMerchantBankAccount],
			OTA:         creditTransferFields[BOTIDMerchantOTA],
			Segments:    &creditTransferSegments,
		},
		BillPayment: &BillPayment{
			AID:        billPaymentFields[BOTIDBillPaymentAID],
			BillerID:   billPaymentFields[BOTIDBillPaymentBillerID],
			Reference1: billPaymentFields[BOTIDBillPaymentRef1],
			Reference2: billPaymentFields[BOTIDBillPaymentRef2],
			Segments:   &billPaymentSegments,
		},
		MerchantCategoryCode:    merchantCategoryCode,
		TransactionCurrency:     transactionCurrency,
		TransactionCurrencyCode: GetCurrencyCode(transactionCurrency),
		TransactionAmount:       transactionAmount,
		CountryCode:             countryCode,
		MerchantName:            merchantName,
		MerchantCity:            merchantCity,
		PostalCode:              postalCode,
		AdditionalFields: &AdditionalFields{
			TerminalID: additionalFields[BOTIDTag62TerminalID],
			Segments:   &additionalSegments,
		},
		CRC:      qrFields[IDCRC],
		Segments: &qrSegments,
	}, nil
}
