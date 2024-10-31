## How to install go package

1. Run go-cli command

```shell
# latest
go get github.com/Jdemon/thaiqr.git

# specific version
go get github.com/Jdemon/thaiqr.git@v1.0.1

## How to Generate QR Payload

### PromptPay QR Payload
``` go
func main() {
	qr := thaiqr.NewPromptPayQR()
	payload, err := qr.GeneratePayload(thaiqr.PromptPayQRCmd{
		ProxyID:   "004999014280076",
		ProxyType: thaiqr.ProxyTypeEWalletID,
		Amount:    "10.00",
	})
	
	// Payload: 00020101021229390016A000000677010111031500499901428007653037645802TH540510.0063046D71
	fmt.Println("Payload: " + payload)
}
```

### Verify Pay Slip QR Payload
``` go
func main() {
	cmd := thaiqr.VerifyPaySlipQRCmd{
		TransactionRef: "2023113077352422",
		SendingBankID:  "006",
		CountryCode:    thaiqr.CountryCodeTH,
	}
	qr := thaiqr.NewVerifyPaySlipQR()
	payload, err := qr.GeneratePayload(cmd)
	
	// Payload: 003700060000010103006021620231130773524225102TH9104EC49
	fmt.Println("Payload: " + payload)
}
```

## How to read QR Payload

### PromptPay QR Payload Reader
``` go
func main() {
    pp_qr := thaiqr.NewPromptPayQR()
	pp_payload, err := pp_qr.GeneratePayload(thaiqr.PromptPayQRCmd{
		ProxyID:   "0909764856",
		ProxyType: thaiqr.ProxyTypeMsisdn,
		Amount:    "1000000.00",
	})
	pp_data, err := pp_qr.Reader(pp_payload)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	pp_marshal, err := json.Marshal(pp_data)
	if err != nil {
		return
	}

	fmt.Printf("promptpay: %s\n", string(pp_marshal))
}
```

### Verify Pay Slip QR Payload Reader
``` go
func main() {
    payload := "003700060000010103006021620231130773524225102TH9104EC49"
    qr := thaiqr.NewVerifyPaySlipQR()
    data, err := qr.Reader(payload)
    if err != nil {
    fmt.Println(err.Error())
    return
    }
    marshal, err := json.Marshal(data)
    if err != nil {
    return
    }
    
    fmt.Printf("qrverify: %s\n", string(marshal))
}
```

## How to Generate QR Image

``` go
qrBtyes, err := thaiqr.GenerateQRWithThaiQRLogo(payload)
```
![PromptpayQR.png](assets%2FPromptpayQR.png)

``` go
qrBtyes, err := thaiqr.GenerateQR(payload)
```
![VerifyQR.png](assets%2FVerifyQR.png)


## Donate

![PromptpayQR.png](assets%2FPromptpayQR.png)
