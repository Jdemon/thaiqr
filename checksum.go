package thaiqr

import (
	"fmt"
)

// checksum calculates the checksum according to ISO/IEC 13239 with the specified polynomial and initial value.
func checksum(data []byte) string {
	// Polynomial '1021' in hexadecimal
	polynomial := uint16(0x1021)
	// Initial value 'FFFF' in hexadecimal
	_checksum := uint16(0xFFFF)

	// Iterate over each byte in the data
	for _, b := range data {
		_checksum ^= uint16(b) << 8

		for i := 0; i < 8; i++ {
			if (_checksum & 0x8000) != 0 {
				_checksum = (_checksum << 1) ^ polynomial
			} else {
				_checksum <<= 1
			}
		}
	}

	// Format the checksum as a hexadecimal string with leading zeros
	return fmt.Sprintf("%04X", _checksum)
}

func VerifyPayloadChecksum(data string) bool {
	if data == "" {
		return false
	}
	payload, crc := splitData(data)
	return crc == checksum([]byte(payload))
}
