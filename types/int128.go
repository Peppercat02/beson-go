package types

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
)

type Int128 struct {
    high uint64
    low uint64
}

func NewInt128(s string, base int) RootType {
    fmt.Println("....")
    // Empty string bad.
    if len(s) == 0 {
        return nil
    }

    // Pick off leading sign.
    neg := false
    if s[0] == '+' {
        s = s[1:]
    } else if s[0] == '-' {
        neg = true
        s = s[1:]
    }

    // Convert unsigned.
    un := NewUInt128(s, base).(*UInt128)

    newValue := &Int128 {
        high: 0,
        low: 0,
    }

    if neg {
        un.twosComplement(un)

        high := un.high
        low := un.low
        newValue.high = high
        newValue.low = low
    } else {
        high := un.high
        low := un.low
        newValue.high = high
        newValue.low = low
    }

    return newValue
}

func ToInt128(value interface{}) RootType {
    var low uint64
    switch value.(type) {
    case int8:
        low = uint64(value.(int8))
    case int16:
        low = uint64(value.(int16))
    case int32:
        low = uint64(value.(int32))
    case int64:
        low = uint64(value.(int64))
    default:
        return nil
    }

    newValue := &Int128 {
        high: 0,
        low: low,
    }
    return newValue
}

func (value *Int128) Rshift(bits uint) *Int128 {
    newValue := &Int128 {
        high: value.high,
        low: value.low,
    }
    value.rightShiftSigned(newValue, bits)
    return newValue
}

func (value *Int128) Lshift(bits uint) *Int128 {
    newValue := &Int128 {
        high: value.high,
        low: value.low,
    }
    value.leftShift(newValue, bits)
    return newValue
}

func (value *Int128) Not() *Int128 {
    newValue := &Int128 {
        high: value.high,
        low: value.low,
    }
    value.not(newValue)
    return newValue
}

func (value *Int128) Or(val *Int128) *Int128 {
    newValue := &Int128 {
        high: value.high,
        low: value.low,
    }
    newValue.or(newValue, val)
    return newValue
}

func (value *Int128) And(val *Int128) *Int128 {
    newValue := &Int128 {
        high: value.high,
        low: value.low,
    }
    newValue.and(newValue, val)
    return newValue
}

func (value *Int128) Xor(val *Int128) *Int128 {
    newValue := &Int128 {
        high: value.high,
        low: value.low,
    }
    newValue.xor(newValue, val)
    return newValue
}

func (value *Int128) Abs() *Int128 {
    newValue := &Int128 {
        high: value.high,
        low: value.low,
    }

    if value.isNegative(newValue) {
        value.twosComplement(newValue)
    }

    return newValue
}

func (value *Int128) Add(val *Int128) *Int128 {
    newValue := &Int128 {
        high: value.high,
        low: value.low,
    }
    newValue.add(newValue, val)
    return newValue
}

func (value *Int128) Sub(val *Int128) *Int128 {
    newValue := &Int128 {
        high: value.high,
        low: value.low,
    }
    newValue.sub(newValue, val)
    return newValue
}

func (value *Int128) Multiply(val *Int128) *Int128 {
    newValue := &Int128 {
        high: value.high,
        low: value.low,
    }
    newValue.multiply(newValue, val)
    return newValue
}

func (value *Int128) Divide(val *Int128) *Int128 {
    a := &Int128 {
        high: value.high,
        low: value.low,
    }
    b := &Int128 {
        high: val.high,
        low: val.low,
    }

    negA := value.isNegative(value)
    if negA {
        value.twosComplement(a)
    }

    negB := val.isNegative(val)
    if negB {
        val.twosComplement(b)
    }

    a.divide(a, b)
    if negA != negB {
        a.twosComplement(a)
    }

    return a
}

func (value *Int128) Modulo(val *Int128) *Int128 {
    a := &Int128 {
        high: value.high,
        low: value.low,
    }
    b := &Int128 {
        high: val.high,
        low: val.low,
    }

    negA := value.isNegative(value)
    if negA {
        a.twosComplement(a)
    }

    negB := val.isNegative(val)
    if negB {
        b.twosComplement(b)
    }

    remainder := a.divide(a, b)
    if negA {
        remainder.twosComplement(remainder)
    }

    return remainder
}

func (value *Int128) Compare(val *Int128) int {
    return value.compare(value, val)
}

func (value *Int128) IsZero() bool {
    return value.isZero(value)
}

func (value *Int128) IsNegative() bool {
    return value.isNegative(value)
}

func (value *Int128) ToString(base int) (string, error) {
    switch base {
    case 2:
        return value.toBinaryString(value), nil
    case 10:
        return value.toDecimalStringSigned(value), nil
    case 16:
        return value.toHexString(value), nil
    default:
        return "", errors.New("Only accepts 2 and 16 representations.")
    }
    return "", nil
}

func (value *Int128) ToBytes() *bytes.Buffer {
    bytesBuffer := bytes.NewBuffer(make([]byte, 0))
    binary.Write(bytesBuffer, binary.LittleEndian, value)

    return bytesBuffer
}

func (value *Int128) IsSigned() bool {
    return true
}

func (value *Int128) SetValue(str string, base int) {
    newValue := NewInt128(str, base).(*Int128)
    value.high = newValue.high
    value.low = newValue.low
}

func (value *Int128) High() uint64 {
    return value.high
}

func (value *Int128) SetHigh(high uint64) {
    value.high = high
}

func (value *Int128) Low() uint64 {
    return value.low
}

func (value *Int128) SetLow(low uint64) {
    value.low = low
}

func (value *Int128) ZERO() *Int128 {
    newValue := &Int128 {
        high: 0,
        low: 0,
    }

    return newValue;
}

func (value *Int128) MAX() *Int128 {
    newValue := &Int128 {
        high: 0x7FFFFFFFFFFFFFFF,
        low: 0xFFFFFFFFFFFFFFFF,
    }

    return newValue;
}

func (value *Int128) MIN() *Int128 {
    newValue := &Int128 {
        high: 0x8000000000000000,
        low: 0x8000000000000000,
    }

    return newValue;
}
