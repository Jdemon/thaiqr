package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Jdemon/thaiqr"
	"image"
	"image/png"
	"os"
)

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

	qrVerify, err := thaiqr.GenerateQR(payload)
	if err != nil {
		return
	}
	bytesToImage(*qrVerify, "VerifyQR.png")

	ppQR := thaiqr.NewPromptPayQR()
	if err != nil {
		return
	}
	ppPayload, err := ppQR.GeneratePayload(thaiqr.PromptPayQRCmd{
		ProxyID:   "0909764856",
		ProxyType: thaiqr.ProxyTypeMsisdn,
		Amount:    "1000000.00",
	})
	if err != nil {
		return
	}
	ppData, err := ppQR.Reader(ppPayload)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ppMarshal, err := json.Marshal(ppData)
	if err != nil {
		return
	}

	fmt.Printf("promptpay: %s\n", string(ppMarshal))

	qrPP, err := thaiqr.GenerateQRWithThaiQRLogo(ppPayload)
	if err != nil {
		return
	}
	bytesToImage(*qrPP, "PromptPayQR.png")
}

func bytesToImage(imgByte []byte, name string) {
	img, _, _ := image.Decode(bytes.NewReader(imgByte))

	//save the imgByte to file
	out, err := os.Create("assets/" + name)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = png.Encode(out, img)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
