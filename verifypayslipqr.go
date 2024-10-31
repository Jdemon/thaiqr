package thaiqr

import (
	"errors"
	"fmt"
)

const (
	IDPayloadAPIID          = "00"
	IDPayloadSendingBankID  = "01"
	IDPayloadTransactionRef = "02"
	VerifyPaySlipAPIID      = "000001"
	CountryCodeLA           = "LA"
)

const (
	IDQrVerifyPayload     = "00"
	IDQrVerifyCountryCode = "51"
	IDQrVerifyCRC         = "91"
)

type VerifyPaySlipQRCmd struct {
	TransactionRef string
	SendingBankID  string
	CountryCode    string
}

type VerifyPaySlipQRResult struct {
	Payload     Payload    `json:"payload"`
	CountryCode string     `json:"countryCode"`
	CRC         string     `json:"crc"`
	Segments    *[]Segment `json:"segments,omitempty"`
}

type Payload struct {
	APIID          string     `json:"apiId"`
	TransactionRef string     `json:"transactionRef"`
	SendingBankID  string     `json:"sendingBankId"`
	Segments       *[]Segment `json:"segments,omitempty"`
}

// VerifyPaySlipQR represents a VerifyPaySlip QR code generator.
type VerifyPaySlipQR struct {
}

// NewVerifyPaySlipQR returns a new VerifyPaySlipQR instance.
func NewVerifyPaySlipQR() *VerifyPaySlipQR {
	return &VerifyPaySlipQR{}
}

// GeneratePayload generates a VerifyPaySlip QR code payload.
func (qr *VerifyPaySlipQR) GeneratePayload(cmd VerifyPaySlipQRCmd) (string, error) {
	sendingBankID := sanitizeTarget(cmd.SendingBankID)
	payload := []string{
		formatField(IDPayloadAPIID, VerifyPaySlipAPIID),
		formatField(IDPayloadSendingBankID, sendingBankID),
		formatField(IDPayloadTransactionRef, cmd.TransactionRef),
	}

	data := []string{
		formatField(IDQrVerifyPayload, serialize(payload)),
		formatField(IDQrVerifyCountryCode, cmd.CountryCode),
	}

	dataToCrc := fmt.Sprintf("%s%s%s", serialize(data), IDQrVerifyCRC, "04")
	data = append(data, formatField(IDQrVerifyCRC, checksum([]byte(dataToCrc))))

	return serialize(data), nil
}

func (qr *VerifyPaySlipQR) Reader(data string) (*VerifyPaySlipQRResult, error) {
	if data == "" || len(data) < 40 {
		return nil, invalidFormat()
	}

	if !VerifyPayloadChecksum(data) {
		return nil, errors.New("invalid checksum")
	}

	qrFields, qrSegments, err := deserialize(data)
	if err != nil {
		return nil, err
	}

	payloadFields, payloadSegments, err := deserialize(qrFields[IDQrVerifyPayload])
	if err != nil {
		return nil, err
	}

	countryCode := qrFields[IDQrVerifyCountryCode]

	apiID := payloadFields[IDPayloadAPIID]
	if apiID != VerifyPaySlipAPIID {
		return nil, invalidFormat()
	}

	return &VerifyPaySlipQRResult{
		Payload: Payload{
			APIID:          apiID,
			SendingBankID:  payloadFields[IDPayloadSendingBankID],
			TransactionRef: payloadFields[IDPayloadTransactionRef],
			Segments:       &payloadSegments,
		},
		CountryCode: countryCode,
		CRC:         qrFields[IDQrVerifyCRC],
		Segments:    &qrSegments,
	}, nil
}
