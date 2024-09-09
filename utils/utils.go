package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	LocServer *time.Location
)

func init() {
	LocServer = MustLoadLocation(os.Getenv("HOST_LOCATION"))
}

// PUBLIC METHOD
func MustLoadLocation(name string) *time.Location {
	l, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("time util: could not load location %s: %s", name, err.Error()))
	}
	return l
}

func TimeNow() time.Time {
	return time.Now().In(LocServer)
}

func ConvertStrToInt(number string, defaultResult int) int {
	result, err := strconv.Atoi(number)
	if err != nil {
		return defaultResult
	}
	return result
}

func UnescapeString(text string) string {
	res, _ := url.QueryUnescape(strings.TrimSpace(strings.ToLower(text)))
	return res
}

func DecodeHttpResponse(data io.ReadCloser, dest interface{}) {
	json.NewDecoder(data).Decode(dest)
}

func ObjToByte(data interface{}) []byte {
	v, _ := json.Marshal(data)
	return v
}

func ConvertStructToMap(v any) map[string]interface{} {
	var result map[string]interface{}
	vByte, _ := json.Marshal(v)
	json.Unmarshal(vByte, &result)

	return result
}

func ConvertMapToStruct(data map[string]interface{}, target interface{}) error {
	byteData, _ := json.Marshal(data)
	if err := json.Unmarshal(byteData, target); err != nil {
		return err
	}
	return nil
}

func PtrToValue[T any](valuePtr *T) T {
	var v T
	if valuePtr != nil {
		v = *valuePtr
	}
	return v
}

func CustomError(message string) error {
	return errors.New(message)
}

func ConvertErrorToMap(key string, err error) map[string]string {
	return map[string]string{key: err.Error()}
}

func GenerateSlug(text string) string {
	result := strings.ToLower(strings.TrimSpace(text))
	result = strings.Replace(result, " ", "-", -1)
	return result
}
