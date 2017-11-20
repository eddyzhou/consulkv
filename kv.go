package consulkv

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
)

type ConfKV struct {
	kv          *api.KVPair
	valueMapper ValueMapper
}

type ValueMapper func(string) string

func (c *ConfKV) String() string {
	val := string(c.kv.Value)
	if c.valueMapper != nil {
		val = c.valueMapper(val)
	}
	return val
}

func (c *ConfKV) Validate(fn func(string) bool) bool {
	return fn(c.String())
}

func parseBool(str string) (value bool, err error) {
	switch str {
	case "1", "true", "TRUE", "True", "YES", "yes", "Yes", "y", "ON", "on", "On":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "n", "OFF", "off", "Off":
		return false, nil
	}
	return false, fmt.Errorf("parsing \"%s\": invalid syntax", str)
}

func (c *ConfKV) Bool() (bool, error) {
	return parseBool(c.String())
}

func (c *ConfKV) Float64() (float64, error) {
	return strconv.ParseFloat(c.String(), 64)
}

func (c *ConfKV) Int() (int, error) {
	return strconv.Atoi(c.String())
}

func (c *ConfKV) Int64() (int64, error) {
	return strconv.ParseInt(c.String(), 10, 64)
}

func (c *ConfKV) Uint() (uint, error) {
	u, e := strconv.ParseUint(c.String(), 10, 64)
	return uint(u), e
}

func (c *ConfKV) Uint64() (uint64, error) {
	return strconv.ParseUint(c.String(), 10, 64)
}

func (c *ConfKV) Duration() (time.Duration, error) {
	return time.ParseDuration(c.String())
}

func (c *ConfKV) Time(format string) (time.Time, error) {
	return time.Parse(format, c.String())
}

func (c *ConfKV) JsonObject() (map[string]interface{}, error) {
	var ret map[string]interface{}
	err := json.Unmarshal([]byte(c.String()), &ret)
	return ret, err
}

func (c *ConfKV) JsonArray() ([]interface{}, error) {
	var ret []interface{}
	err := json.Unmarshal([]byte(c.String()), &ret)
	return ret, err
}

func (c *ConfKV) MustString(defaultVal string) string {
	val := c.String()
	if len(val) == 0 {
		return defaultVal
	}
	return val
}

func (c *ConfKV) MustBool(defaultVal ...bool) bool {
	val, err := c.Bool()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

func (c *ConfKV) MustFloat64(defaultVal ...float64) float64 {
	val, err := c.Float64()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

func (c *ConfKV) MustInt(defaultVal ...int) int {
	val, err := c.Int()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

func (c *ConfKV) MustInt64(defaultVal ...int64) int64 {
	val, err := c.Int64()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

func (c *ConfKV) MustUint(defaultVal ...uint) uint {
	val, err := c.Uint()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

func (c *ConfKV) MustUint64(defaultVal ...uint64) uint64 {
	val, err := c.Uint64()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

func (c *ConfKV) MustDuration(defaultVal ...time.Duration) time.Duration {
	val, err := c.Duration()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

func (c *ConfKV) MustTime(format string, defaultVal ...time.Time) time.Time {
	val, err := c.Time(format)
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

func (c *ConfKV) MustJsonObject(defaultVal ...map[string]interface{}) map[string]interface{} {
	val, err := c.JsonObject()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}

func (c *ConfKV) MustJsonArray(defaultVal ...[]interface{}) []interface{} {
	val, err := c.JsonArray()
	if len(defaultVal) > 0 && err != nil {
		return defaultVal[0]
	}
	return val
}
