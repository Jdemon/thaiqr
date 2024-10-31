package thaiqr

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Segment struct {
	RawValue string `json:"rawValue"`
	ID       string `json:"id"`
	Length   int    `json:"length"`
	Value    string `json:"value"`
}

// sanitizeTarget removes non-numeric characters from the target value.
func sanitizeTarget(value string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(value, "")
}

// serialize concatenates a slice of strings into a single string.
func serialize(values []string) string {
	return strings.Join(values, "")
}

// f formats an ID and a value into a string.
func formatField(id, value string) string {
	ext := "00" + strconv.Itoa(len(value))
	return id + ext[len(ext)-2:] + value
}

// formatTarget sanitizes and formats the target value.
func formatTarget(value string) string {
	value = sanitizeTarget(value)
	if len(value) >= 13 {
		return value
	}

	re := regexp.MustCompile(`^0`)
	value = re.ReplaceAllString(value, "66")
	value = "0000000000000" + value

	return value[len(value)-13:]
}

// formatAmount converts the amount to a formatted string.
func formatAmount(amount string) (string, error) {
	if f, err := strconv.ParseFloat(amount, 32); err == nil {
		return fmt.Sprintf("%.2f", f), nil
	}
	return "", errors.New("invalid amount")
}

// ifThenElse returns 'a' if the condition is true, otherwise 'b'.
func ifThenElse(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func deserialize(data string) (map[string]string, []Segment, error) {
	if data == "" {
		return nil, nil, invalidFormat()
	}

	segments := make([]Segment, 0)

	remain := data
	results := make(map[string]string)
	for remain != "" {
		var value string
		keyAndLen, payload, err := getAndRemove(remain, 4)
		if err != nil {
			return nil, nil, err
		}
		key, lengthStr, err := getAndRemove(keyAndLen, 2)
		if err != nil {
			return nil, nil, err
		}
		length, err := parseInt(lengthStr)
		if err != nil {
			return nil, nil, err
		}
		value, remain, err = getAndRemove(payload, length)
		if err != nil {
			return nil, nil, err
		}
		segments = append(segments, Segment{
			RawValue: fmt.Sprintf("%s%s%s", key, lengthStr, value),
			ID:       key,
			Length:   length,
			Value:    value,
		})
		results[key] = value
	}

	return results, segments, nil
}

func getAndRemove(data string, length int) (string, string, error) {
	if len(data) < length {
		return "", "", invalidFormat()
	}
	key := data[:length]
	data = data[length:]
	return key, data, nil
}

func parseInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func splitData(data string) (string, string) {
	splitIndex := len(data) - 4
	return data[:splitIndex], data[splitIndex:]

}

func invalidFormat() error {
	return errors.New("invalid format")
}
