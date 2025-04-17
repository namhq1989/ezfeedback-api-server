package database

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type ArrayString []string

func (s *ArrayString) Scan(value interface{}) error {
	input := strings.Trim(value.(string), "{}")
	*s = strings.Split(input, ",")

	return nil
}

func (s *ArrayString) Value() (driver.Value, error) {
	return "{" + strings.Join(*s, ",") + "}", nil
}

type CountResult struct {
	Total int64 `json:"total"`
}

func SQLArrayStrings(values []string) string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = fmt.Sprintf("%s", value)
	}
	return "{" + strings.Join(result, ", ") + "}"
}

type ArrayInt []int

func (s *ArrayInt) Scan(value interface{}) error {
	input := strings.Trim(value.(string), "{}")
	arr := strings.Split(input, ",")

	intArray := make([]int, len(arr))
	for i, v := range arr {
		num, _ := strconv.Atoi(v)
		intArray[i] = num
	}

	*s = intArray
	return nil
}

func (s *ArrayInt) Value() (driver.Value, error) {
	strArray := make([]string, len(*s))

	for i, v := range *s {
		strArray[i] = strconv.Itoa(v)
	}
	return "{" + strings.Join(strArray, ",") + "}", nil
}
