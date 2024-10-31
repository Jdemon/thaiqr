package thaiqr_test

import (
	"github.com/Jdemon/thaiqr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateBillPaymentPayloadMustValid(t *testing.T) {
	expectedPayload := "00020101021230570016A00000067701011201153110400394751010206REF0010304REF253037645406555.555802TH62100706SCB001630437C6" // SCB
	qr := thaiqr.NewPromptPayQR()
	cmd := thaiqr.PromptPayBillPaymentQRCmd{
		BillerID:   "311040039475101",
		Ref1:       "REF001",
		Ref2:       "REF2",
		TerminalID: "SCB001",
		Amount:     "555.55",
	}
	actualPayload, err := qr.GenerateBillPaymentPayload(cmd)
	assert.Nil(t, err)
	assert.Equal(t, expectedPayload, actualPayload)

	result, err := qr.Reader(actualPayload)
	assert.Nil(t, err)
	assert.Equal(t, cmd.BillerID, result.BillPayment.BillerID)
	assert.Equal(t, cmd.Ref1, result.BillPayment.Reference1)
	assert.Equal(t, cmd.Ref2, result.BillPayment.Reference2)
	assert.Equal(t, cmd.TerminalID, result.AdditionalFields.TerminalID)
	assert.Equal(t, cmd.Amount, result.TransactionAmount)
	assert.Equal(t, thaiqr.CountryCodeTH, result.CountryCode)
	assert.Equal(t, "THB", result.TransactionCurrencyCode)
	assert.Equal(t, thaiqr.TransactionCurrencyTHB, result.TransactionCurrency)
	assert.Equal(t, actualPayload[len(actualPayload)-4:], result.CRC)
}

func TestGenerateBillPaymentPayloadInvalidAmount(t *testing.T) {
	qr := thaiqr.NewPromptPayQR()
	_, err := qr.GenerateBillPaymentPayload(thaiqr.PromptPayBillPaymentQRCmd{
		BillerID:   "311040039475101",
		Ref1:       "REF001",
		Ref2:       "REF2",
		TerminalID: "SCB001",
		Amount:     "1xxx.55",
	})
	assert.Error(t, err, "invalid amount")
}

func TestGeneratePromptPayMobileNoWithoutAmount(t *testing.T) {
	expectedPayload := "00020101021129370016A0000006770101110113006690976485653037645802TH63044D1B"
	qr := thaiqr.NewPromptPayQR()
	cmd := thaiqr.PromptPayQRCmd{
		ProxyID:   "0909764856",
		ProxyType: thaiqr.ProxyTypeMsisdn,
	}
	actualPayload, err := qr.GeneratePayload(cmd)
	assert.Nil(t, err)
	assert.Equal(t, expectedPayload, actualPayload)

	result, err := qr.Reader(actualPayload)
	assert.Nil(t, err)
	assert.Equal(t, "0066909764856", result.CreditTransfer.MSISDN)
	assert.Equal(t, cmd.Amount, result.TransactionAmount)
	assert.Equal(t, cmd.OTA, result.CreditTransfer.OTA)
	assert.Equal(t, thaiqr.CountryCodeTH, result.CountryCode)
	assert.Equal(t, "THB", result.TransactionCurrencyCode)
	assert.Equal(t, thaiqr.TransactionCurrencyTHB, result.TransactionCurrency)
	assert.Equal(t, actualPayload[len(actualPayload)-4:], result.CRC)
}

func TestGeneratePromptPayMobileNo(t *testing.T) {
	expectedPayload := "00020101021229370016A0000006770101110113006690976485653037645802TH540510.006304CF65"
	qr := thaiqr.NewPromptPayQR()
	cmd := thaiqr.PromptPayQRCmd{
		ProxyID:   "0909764856",
		ProxyType: thaiqr.ProxyTypeMsisdn,
		Amount:    "10.00",
	}
	actualPayload, err := qr.GeneratePayload(cmd)
	assert.Nil(t, err)
	assert.Equal(t, expectedPayload, actualPayload)

	result, err := qr.Reader(actualPayload)
	assert.Nil(t, err)
	assert.Equal(t, "0066909764856", result.CreditTransfer.MSISDN)
	assert.Equal(t, cmd.Amount, result.TransactionAmount)
	assert.Equal(t, cmd.OTA, result.CreditTransfer.OTA)
	assert.Equal(t, thaiqr.CountryCodeTH, result.CountryCode)
	assert.Equal(t, "THB", result.TransactionCurrencyCode)
	assert.Equal(t, thaiqr.TransactionCurrencyTHB, result.TransactionCurrency)
	assert.Equal(t, actualPayload[len(actualPayload)-4:], result.CRC)
}

func TestGeneratePromptPayMobileNoInvalidAmount(t *testing.T) {
	qr := thaiqr.NewPromptPayQR()
	_, err := qr.GeneratePayload(thaiqr.PromptPayQRCmd{
		ProxyID:   "0909764856",
		ProxyType: thaiqr.ProxyTypeMsisdn,
		Amount:    "10xx.00",
	})
	assert.Error(t, err, "invalid amount")
}

func TestGeneratePromptPayCIDWithoutAmount(t *testing.T) {
	expectedPayload := "00020101021129370016A0000006770101110213110060146718253037645802TH6304C7EE"
	qr := thaiqr.NewPromptPayQR()
	cmd := thaiqr.PromptPayQRCmd{
		ProxyID:   "1100601467182",
		ProxyType: thaiqr.ProxyTypeNatID,
	}
	actualPayload, err := qr.GeneratePayload(cmd)
	assert.Nil(t, err)
	assert.Equal(t, expectedPayload, actualPayload)

	result, err := qr.Reader(actualPayload)
	assert.Nil(t, err)
	assert.Equal(t, cmd.ProxyID, result.CreditTransfer.NationalID)
	assert.Equal(t, cmd.Amount, result.TransactionAmount)
	assert.Equal(t, cmd.OTA, result.CreditTransfer.OTA)
	assert.Equal(t, thaiqr.CountryCodeTH, result.CountryCode)
	assert.Equal(t, "THB", result.TransactionCurrencyCode)
	assert.Equal(t, thaiqr.TransactionCurrencyTHB, result.TransactionCurrency)
	assert.Equal(t, actualPayload[len(actualPayload)-4:], result.CRC)
}

func TestGeneratePromptPayCID(t *testing.T) {
	expectedPayload := "00020101021229370016A0000006770101110213110060146718253037645802TH540510.0063049A7C"
	qr := thaiqr.NewPromptPayQR()
	cmd := thaiqr.PromptPayQRCmd{
		ProxyID:   "1100601467182",
		ProxyType: thaiqr.ProxyTypeNatID,
		Amount:    "10.00",
	}
	actualPayload, err := qr.GeneratePayload(cmd)
	assert.Nil(t, err)
	assert.Equal(t, expectedPayload, actualPayload)

	result, err := qr.Reader(actualPayload)
	assert.Nil(t, err)
	assert.Equal(t, cmd.ProxyID, result.CreditTransfer.NationalID)
	assert.Equal(t, cmd.Amount, result.TransactionAmount)
	assert.Equal(t, cmd.OTA, result.CreditTransfer.OTA)
	assert.Equal(t, thaiqr.CountryCodeTH, result.CountryCode)
	assert.Equal(t, "THB", result.TransactionCurrencyCode)
	assert.Equal(t, thaiqr.TransactionCurrencyTHB, result.TransactionCurrency)
	assert.Equal(t, actualPayload[len(actualPayload)-4:], result.CRC)
}

func TestGeneratePromptPayCIDInvalidAmount(t *testing.T) {

	qr := thaiqr.NewPromptPayQR()
	_, err := qr.GeneratePayload(thaiqr.PromptPayQRCmd{
		ProxyID:   "1100601467182",
		ProxyType: thaiqr.ProxyTypeNatID,
		Amount:    "10xx.00",
	})
	assert.Error(t, err, "invalid amount")
}

func TestGeneratePromptPayEWalletWithoutAmount(t *testing.T) {
	expectedPayload := "00020101021129390016A000000677010111031500499901428007653037645802TH63044541"
	qr := thaiqr.NewPromptPayQR()
	cmd := thaiqr.PromptPayQRCmd{
		ProxyID:   "004999014280076",
		ProxyType: thaiqr.ProxyTypeEWalletID,
	}
	actualPayload, err := qr.GeneratePayload(cmd)
	assert.Nil(t, err)
	assert.Equal(t, expectedPayload, actualPayload)

	result, err := qr.Reader(actualPayload)
	assert.Nil(t, err)
	assert.Equal(t, cmd.ProxyID, result.CreditTransfer.EWalletID)
	assert.Equal(t, cmd.Amount, result.TransactionAmount)
	assert.Equal(t, cmd.OTA, result.CreditTransfer.OTA)
	assert.Equal(t, thaiqr.CountryCodeTH, result.CountryCode)
	assert.Equal(t, "THB", result.TransactionCurrencyCode)
	assert.Equal(t, thaiqr.TransactionCurrencyTHB, result.TransactionCurrency)
	assert.Equal(t, actualPayload[len(actualPayload)-4:], result.CRC)
}

func TestGeneratePromptPayEWallet(t *testing.T) {
	expectedPayload := "00020101021229390016A000000677010111031500499901428007653037645802TH540510.0063046D71"
	qr := thaiqr.NewPromptPayQR()
	cmd := thaiqr.PromptPayQRCmd{
		ProxyID:   "004999014280076",
		ProxyType: thaiqr.ProxyTypeEWalletID,
		Amount:    "10.00",
	}
	actualPayload, err := qr.GeneratePayload(cmd)
	assert.Nil(t, err)
	assert.Equal(t, expectedPayload, actualPayload)

	result, err := qr.Reader(actualPayload)
	assert.Nil(t, err)
	assert.Equal(t, cmd.ProxyID, result.CreditTransfer.EWalletID)
	assert.Equal(t, cmd.Amount, result.TransactionAmount)
	assert.Equal(t, cmd.OTA, result.CreditTransfer.OTA)
	assert.Equal(t, thaiqr.CountryCodeTH, result.CountryCode)
	assert.Equal(t, "THB", result.TransactionCurrencyCode)
	assert.Equal(t, thaiqr.TransactionCurrencyTHB, result.TransactionCurrency)
	assert.Equal(t, actualPayload[len(actualPayload)-4:], result.CRC)
}

func TestVerifyPayloadChecksum(t *testing.T) {
	isValid := thaiqr.VerifyPayloadChecksum("00020101021229390016A000000677010111031500499901428007653037645802TH540510.0063046D71")
	assert.True(t, isValid)
}
