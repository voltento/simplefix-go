package fix

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

// Value is an interface implementing all basic methods required to process field values of FIX messages.
type Value interface {
	// ToBytes returns a byte representation of a field value.
	ToBytes() []byte
	WriteBytes(writer *bytes.Buffer) bool

	// FromBytes parses values stored in a byte array.
	FromBytes([]byte) error

	// Value returns a field value.
	Value() interface{}

	// String returns a string representation of a value.
	String() string

	// IsNull is used to check whether a field value is not filled
	IsNull() bool

	// IsNull is used to check whether a field value is empty.
	IsEmpty() bool

	// Set replaces a specified field value with a value of a corresponding type.
	Set(d interface{}) error
}

// Raw is a structure that is used to convert the message data to a byte array.
type Raw struct {
	value []byte
}

// NewRaw creates a new instance of a Raw object.
func NewRaw(v []byte) *Raw {
	return &Raw{
		value: v,
	}
}

func (v *Raw) ToBytes() []byte {
	return v.value
}

func (v *Raw) WriteBytes(writer *bytes.Buffer) bool {
	_, _ = writer.Write(v.value)
	return true
}

func (v *Raw) FromBytes(d []byte) (err error) {
	v.value = d
	return nil
}

func (v *Raw) IsNull() bool {
	return v.value == nil
}
func (v *Raw) IsEmpty() bool {
	return v.value == nil
}
func (v *Raw) Value() interface{} {
	return v.value
}

// Set parses and assigns the field value stored as a byte array.
func (v *Raw) Set(d interface{}) error {
	if res, ok := d.([]byte); ok {
		v.value = res
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "Byte Array")
}

func (v *Raw) String() string {
	return string(v.value)
}

// String is a structure used for converting string values.
type String struct {
	value string
	valid bool
}

func NewString(v string) *String {
	return &String{value: v, valid: true}
}

// Set parses and assigns the field value stored as a string.
func (v *String) Set(d interface{}) error {
	if d == nil {
		v.valid = false
		return nil
	}

	if res, ok := d.(string); ok {
		v.value = res
		v.valid = true
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "String")
}

func (v *String) ToBytes() []byte {
	if !v.valid || v.value == "" {
		return nil
	}
	return []byte(v.value)
}
func (v *String) WriteBytes(writer *bytes.Buffer) bool {
	if !v.valid || v.value == "" {
		return false
	}
	_, _ = writer.WriteString(v.value)
	return true
}

func (v *String) IsNull() bool {
	return !v.valid
}
func (v *String) IsEmpty() bool {
	return !v.valid || v.value == ""
}
func (v *String) Value() interface{} {
	return v.value
}

func (v *String) FromBytes(d []byte) (err error) {
	if d == nil {
		v.valid = false
		return nil
	}

	v.valid = true
	v.value = string(d)

	return nil
}

func (v *String) String() string {
	return v.value
}

// Int is a structure used for converting integer values.
type Int struct {
	value int
	valid bool
}

func NewInt(value int) *Int {
	return &Int{valid: true, value: value}
}

func (v *Int) IsNull() bool {
	return !v.valid
}
func (v *Int) IsEmpty() bool {
	return !v.valid
}

// Set parses and assigns the field value stored as an integer number.
func (v *Int) Set(d interface{}) error {
	if d == nil {
		v.valid = false
		return nil
	}

	if res, ok := d.(int); ok {
		v.value = res
		v.valid = true
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "Int")
}

func (v *Int) Value() interface{} {
	return v.value
}
func (v *Int) String() string {
	return strconv.Itoa(v.value)
}

func (v *Int) FromBytes(d []byte) (err error) {
	if d == nil {
		v.valid = false
		return nil
	}

	v.valid = true
	v.value, err = bytesToInt(d)

	return err
}

func (v *Int) ToBytes() []byte {
	if !v.valid {
		return nil
	}
	return intToBytes(v.value)
}
func (v *Int) WriteBytes(writer *bytes.Buffer) bool {
	if !v.valid {
		return false
	}
	_, _ = writer.Write(intToBytes(v.value))
	return true
}

// Uint is a structure used for converting values to the uint64 type.
type Uint struct {
	value uint64
	valid bool
}

func NewUint(value uint64) *Uint {
	return &Uint{value: value, valid: true}
}

// Set parses and assigns the field value stored as a uint64 number.
func (v *Uint) Set(d interface{}) error {
	if d == nil {
		v.valid = false
		return nil
	}

	if res, ok := d.(uint64); ok {
		v.value = res
		v.valid = true
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "Uint")
}

func (v *Uint) IsNull() bool {
	return !v.valid
}
func (v *Uint) IsEmpty() bool {
	return !v.valid
}
func (v *Uint) FromBytes(d []byte) (err error) {
	if d == nil {
		v.valid = false
		return nil
	}

	v.valid = true
	v.value, err = bytesToUint(d)

	return err
}

func (v *Uint) Value() interface{} {
	return v.value
}

func (v *Uint) String() string {
	return fmt.Sprintf("%d", v.value)
}

func (v *Uint) ToBytes() []byte {
	if !v.valid {
		return nil
	}
	return uintToBytes(v.value)
}
func (v *Uint) WriteBytes(writer *bytes.Buffer) bool {
	if !v.valid {
		return false
	}
	_, _ = writer.Write(uintToBytes(v.value))
	return true
}

// Float is a structure used for converting values to the float64 type.
type Float struct {
	source []byte
	value  float64
	valid  bool
}

func NewFloat(value float64) *Float {
	return &Float{value: value, valid: true}
}

func (v *Float) IsNull() bool {
	return !v.valid
}
func (v *Float) IsEmpty() bool {
	return !v.valid
}
func (v *Float) Value() interface{} {
	return v.value
}

func (v *Float) FromBytes(d []byte) (err error) {
	if d == nil {
		v.valid = false
		return nil
	}

	v.valid = true
	v.source = d
	v.value, err = bytesToFloat(d)

	return err
}

func (v *Float) ToBytes() []byte {
	if !v.valid {
		return nil
	}
	if v.source != nil {
		return v.source
	}
	return floatToBytes(v.value)
}
func (v *Float) WriteBytes(writer *bytes.Buffer) bool {
	if !v.valid {
		return false
	}

	if v.source != nil {
		_, _ = writer.Write(v.source)
	} else {
		_, _ = writer.Write(strconv.AppendFloat(make([]byte, 0, 64), v.value, 'f', -1, 64))
	}
	return true
}
func (v *Float) String() string {
	return fmt.Sprintf("%f", v.value)
}

// Set parses and assigns the field value stored as a float64 number.
func (v *Float) Set(d interface{}) error {
	if d == nil {
		v.valid = false
		return nil
	}

	if res, ok := d.(float64); ok {
		v.value = res
		v.valid = true
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "Float")
}

// Time is a structure used for converting date-time values.
type Time struct {
	value time.Time
	valid bool
}

func NewTime(value time.Time) *Time {
	return &Time{value: value, valid: true}
}

// Set parses and assigns the field value stored in the date-time format.
func (v *Time) Set(d interface{}) error {
	if d == nil {
		v.valid = false
		return nil
	}

	if res, ok := d.(time.Time); ok {
		v.value = res
		v.valid = true
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "Date-Time")
}

func (v *Time) IsNull() bool {
	return !v.valid
}
func (v *Time) IsEmpty() bool {
	return !v.valid
}

func (v *Time) Value() interface{} {
	return v.value
}

func (v *Time) ToBytes() []byte {
	if !v.valid {
		return nil
	}
	return timeToBytes(v.value)
}
func (v *Time) WriteBytes(writer *bytes.Buffer) bool {
	if !v.valid {
		return false
	}
	_, _ = writer.Write(timeToBytes(v.value))
	return true
}
func (v *Time) FromBytes(d []byte) (err error) {
	if d == nil {
		v.valid = false
		return nil
	}

	v.valid = true
	v.value, err = time.Parse(TimeLayout, string(d))

	return err
}

func (v *Time) String() string {
	return v.value.Format(TimeLayout)
}

const (
	True  = 'Y'
	False = 'N'
)

var trueByte = []byte{True}
var falseByte = []byte{False}

// Bool is a structure used for converting Boolean values.
type Bool struct {
	value bool
	valid bool
}

func (v *Bool) ToBytes() []byte {
	if !v.valid {
		return nil
	}

	if v.value {
		return trueByte
	}
	return falseByte
}
func (v *Bool) WriteBytes(writer *bytes.Buffer) bool {
	if !v.valid {
		return false
	}

	if v.value {
		_, _ = writer.Write(trueByte)
	} else {
		_, _ = writer.Write(falseByte)
	}
	return true
}

func (v *Bool) FromBytes(d []byte) error {
	if d == nil {
		v.valid = false
		return nil
	}

	v.valid = true
	v.value = len(d) > 0 && d[0] == True

	return nil
}

func (v *Bool) Value() interface{} {
	return v.value
}

func (v *Bool) String() string {
	if !v.valid {
		return ""
	}

	if v.value {
		return string(True)
	}
	return string(False)
}

func (v *Bool) IsNull() bool {
	return !v.valid
}
func (v *Bool) IsEmpty() bool {
	return !v.valid
}

// Set parses and assigns the field value stored in the Boolean format.
func (v *Bool) Set(d interface{}) error {
	if d == nil {
		v.valid = false
		return nil
	}

	if res, ok := d.(bool); ok {
		v.value = res
		v.valid = true
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "Boolean")
}

func bytesToUint(d []byte) (uint64, error) {
	if len(d) == 0 {
		return 0, errors.New("invalid input: empty byte slice")
	}
	var result uint64
	for i := 0; i < len(d); i++ {
		c := d[i]
		if c < '0' || c > '9' {
			return 0, errors.New("invalid input: non-numeric character")
		}
		digit := uint64(c - '0')
		if result > math.MaxUint64/10 || (result == math.MaxUint64/10 && digit > math.MaxUint64%10) {
			return 0, errors.New("overflow: value too large for uint64")
		}
		result = result*10 + digit
	}
	return result, nil
}
func uintToBytes(value uint64) []byte {
	return strconv.AppendUint(make([]byte, 0, 20), value, 10)
}

// not working with negative values
func timeToBytes(t time.Time) []byte {
	year, month, day := t.Date()
	if year < 0 {
		year = 0
	}
	hour, minute, second := t.Clock()
	milli := t.Nanosecond() / 1e6
	return []byte{byte('0' + year/1000),
		byte('0' + (year/100)%10),
		byte('0' + (year/10)%10),
		byte('0' + year%10),

		byte('0' + (month/10)%10),
		byte('0' + month%10),

		byte('0' + (day/10)%10),
		byte('0' + day%10),
		'-',
		byte('0' + (hour/10)%10),
		byte('0' + hour%10),
		':',
		byte('0' + (minute/10)%10),
		byte('0' + minute%10),
		':',
		byte('0' + (second/10)%10),
		byte('0' + second%10),
		'.',
		byte('0' + milli/100),
		byte('0' + (milli/10)%10),
		byte('0' + milli%10)}
}

func floatToBytes(f float64) []byte {
	if f == 0 {
		return []byte{'0'}
	}
	return strconv.AppendFloat(make([]byte, 0, 64), f, 'f', -1, 64)
}

var float64pow10 = [...]float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16,
}

var allowedFloat [256]bool

func init() {
	for _, ch := range "0123456789.eE+-" {
		allowedFloat[ch] = true
	}
}

// inspired by https://github.com/valyala/fastjson/blob/6dae91c8e11a7fa6a257a550b75cba53ab81693e/fastfloat/parse.go#L203
func bytesToFloat(s []byte) (float64, error) {
	l := uint(len(s))
	if l == 0 {
		return 0, fmt.Errorf("invalid syntax: empty string")
	}
	i := uint(0)
	minus := s[0] == '-'
	if minus {
		i++
		if i >= l {
			return 0, errors.New("invalid syntax: ony minus sign")
		}
	}

	if l > 18 {
		return strconv.ParseFloat(string(s), 64)
	}

	d := uint64(0)
	j := i
	for i < l {
		if !allowedFloat[s[i]] {
			return 0, errors.New("invalid syntax: invalid character")
		}
		if s[i] >= '0' && s[i] <= '9' {
			d = d*10 + uint64(s[i]-'0')
			i++
			if i > 18 {
				// The integer part may be out of range for uint64.
				// Fall back to slow parsing.
				f, err := strconv.ParseFloat(string(s), 64)
				if err != nil && !math.IsInf(f, 0) {
					return 0, err
				}
				return f, nil
			}
			continue
		}
		break
	}
	f := float64(d)
	if s[0] != '.' || l == 1 {
		if i <= j {
			return 0, errors.New("invalid syntax: unparsable tail left")
		}
		if i >= l {
			if minus {
				return -f, nil
			}
			return f, nil
		}
	}

	if s[i] == '.' {
		// Parse fractional part.
		i++
		if i >= l {
			if l == 2 && s[0] == '0' {
				return 0, nil
			}
			return 0, errors.New("cannot parse fractional part")
		}
		k := i
		for i < l {
			if s[i] >= '0' && s[i] <= '9' {
				d = d*10 + uint64(s[i]-'0')
				i++
				if i-j >= uint(len(float64pow10)) {
					// The mantissa is out of range. Fall back to standard parsing.
					f, err := strconv.ParseFloat(string(s), 64)
					if err != nil && !math.IsInf(f, 0) {
						return 0, errors.New("cannot parse mantissa")
					}
					return f, nil
				}
				continue
			}
			break
		}
		if i < k {
			return 0, errors.New("cannot find mantissa")
		}
		// Convert the entire mantissa to a float at once to avoid rounding errors.
		f = float64(d) / float64pow10[i-k]
		if i >= l {
			// Fast path - parsed fractional number.
			if minus {
				return -f, nil
			}
			return f, nil
		}
	}
	if s[i] == 'e' || s[i] == 'E' {
		// Parse exponent part.
		i++
		if i >= l {
			return 0, errors.New("cannot parse exponent")
		}
		expMinus := false
		if s[i] == '+' || s[i] == '-' {
			expMinus = s[i] == '-'
			i++
			if i >= l {
				return 0, errors.New("cannot parse exponent from string")
			}
		}
		exp := int16(0)
		j := i
		for i < l {
			if s[i] >= '0' && s[i] <= '9' {
				exp = exp*10 + int16(s[i]-'0')
				i++
				if exp > 32 {
					// The exponent may be too big for float64.
					// Fall back to standard parsing.
					f, err := strconv.ParseFloat(string(s), 64)
					if err != nil && !math.IsInf(f, 0) {
						return 0, errors.New("cannot parse exponent by strconv")
					}
					return f, nil
				}
				continue
			}
			break
		}
		if i <= j {
			return 0, errors.New("cannot parse exponent tail")
		}
		if expMinus {
			exp = -exp
		}
		f *= math.Pow10(int(exp))
		if i >= l {
			if minus {
				return -f, nil
			}
			return f, nil
		}
	}
	return 0, errors.New("invalid syntax")
}

func bytesToInt(d []byte) (int, error) {
	if len(d) == 0 {
		return 0, errors.New("invalid input: empty byte slice")
	}

	var result int
	var sign = 1
	start := 0

	if d[0] == '-' {
		sign = -1
		start = 1
	} else if d[0] == '+' {
		start = 1
	}

	for i := start; i < len(d); i++ {
		c := d[i]
		if c < '0' || c > '9' {
			return 0, errors.New("invalid input: non-numeric character")
		}
		result = result*10 + int(c-'0')
	}

	return result * sign, nil
}
func intToBytes(value int) []byte {
	return strconv.AppendInt(make([]byte, 0, 20), int64(value), 10)
}
