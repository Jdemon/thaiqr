package thaiqr

var currencyCode = map[string]string{
	"IDR": "360", // Indonesia
	"MMK": "104", // Myanmar (Burma)
	"BND": "096", // Brunei
	"KHR": "116", // Cambodia
	"LAK": "418", // Laos
	"MYR": "458", // Malaysia
	"PHP": "608", // Philippines
	"SGD": "702", // Singapore
	"THB": "764", // Thailand
	"VND": "704", // Vietnam
}
var currencyNoMap map[string]string

func GetCurrencyCode(no string) string {
	if currencyNoMap == nil {
		currencyNoMap = make(map[string]string)
		for code, number := range currencyCode {
			currencyNoMap[number] = code
		}
	}
	return currencyNoMap[no]
}
