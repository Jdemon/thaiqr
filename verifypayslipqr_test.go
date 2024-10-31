package thaiqr_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateVerifySlipUserStory(t *testing.T) {
	cmd := thaiqr.VerifyPaySlipQRCmd{
		TransactionRef: "1234567890123456789012345",
		SendingBankID:  "001",
		CountryCode:    thaiqr.CountryCodeTH,
	}
	expectedPayload := "004600060000010103001022512345678901234567890123455102TH910408DC"
	qr := thaiqr.NewVerifyPaySlipQR()
	actualPayload, err := qr.GeneratePayload(cmd)
	assert.Nil(t, err)
	assert.Equal(t, expectedPayload, actualPayload)

	result, err := qr.Reader(actualPayload)
	assert.Nil(t, err)
	assert.Equal(t, cmd.TransactionRef, result.Payload.TransactionRef)
	assert.Equal(t, cmd.SendingBankID, result.Payload.SendingBankID)
	assert.Equal(t, cmd.CountryCode, result.CountryCode)
	assert.Equal(t, expectedPayload[len(expectedPayload)-4:], result.CRC)
}

func TestGenerateVerifySlipLaoCountryCode(t *testing.T) {
	cmd := thaiqr.VerifyPaySlipQRCmd{
		TransactionRef: "1234567890123456789012345",
		SendingBankID:  "006",
		CountryCode:    thaiqr.CountryCodeLA,
	}
	expectedPayload := "004600060000010103006022512345678901234567890123455102LA91041E58"
	qr := thaiqr.NewVerifyPaySlipQR()
	actualPayload, err := qr.GeneratePayload(cmd)
	assert.Nil(t, err)
	assert.Equal(t, expectedPayload, actualPayload)

	result, err := qr.Reader(actualPayload)
	assert.Nil(t, err)
	assert.Equal(t, cmd.TransactionRef, result.Payload.TransactionRef)
	assert.Equal(t, cmd.SendingBankID, result.Payload.SendingBankID)
	assert.Equal(t, cmd.CountryCode, result.CountryCode)
	assert.Equal(t, expectedPayload[len(expectedPayload)-4:], result.CRC)
}

func TestGenerateVerifySlipKTB(t *testing.T) {
	cmd := thaiqr.VerifyPaySlipQRCmd{
		TransactionRef: "2023113077352422",
		SendingBankID:  "006",
		CountryCode:    thaiqr.CountryCodeTH,
	}
	expectedPayload := "003700060000010103006021620231130773524225102TH9104EC49"
	qr := thaiqr.NewVerifyPaySlipQR()
	actualPayload, err := qr.GeneratePayload(cmd)
	assert.Nil(t, err)
	assert.Equal(t, expectedPayload, actualPayload)

	result, err := qr.Reader(actualPayload)
	assert.Nil(t, err)
	assert.Equal(t, cmd.TransactionRef, result.Payload.TransactionRef)
	assert.Equal(t, cmd.SendingBankID, result.Payload.SendingBankID)
	assert.Equal(t, cmd.CountryCode, result.CountryCode)
	assert.Equal(t, expectedPayload[len(expectedPayload)-4:], result.CRC)
}
