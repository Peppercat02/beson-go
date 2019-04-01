package helper

import (
    "fmt"
    "bytes"
    "encoding/hex"
    "log"
    "strconv"
    "regexp"
)

var BYTE_MAX byte = 255

func a () {
    fmt.Println()
}

const HEX_FORMAT_CHECKER string = "^0x[0-9a-fA-F]+$";
var HEX_MAP_I = map[rune]uint8 {
    '0':0, '1':1, '2':2, '3':3, '4':4, '5':5, '6':6, '7':7, '8':8, '9':9,
    'a':10, 'b':11, 'c':12, 'd':13, 'e':14, 'f':15,
    'A':10, 'B':11, 'C':12, 'D':13, 'E':14, 'F':15,
};

func HexStringToBytes(s string) []byte {
    isMatch, _ := regexp.MatchString(HEX_FORMAT_CHECKER, s)
    if isMatch {
        s = s[2:]
        decoded, err := hex.DecodeString(s)
        if err != nil {
            log.Fatal(err)
            return nil
        }
        return decoded
    }
    return nil
}

// TODO : BinaryStringToBytes

// TODO : DecimalStringToBytes

func Concat(segments ...[]byte) []byte {
    buffer := bytes.NewBuffer(make([]byte, 0))

    for _, seg := range segments {
        buffer.Write(seg)
    }

    newBytes := make([]byte, buffer.Len())
    buffer.Read(newBytes)

    return newBytes
}

func LeftShift(value []byte, bits uint, padding uint8)  {
    if bits == 0 {
        return
    }
    if padding == 0 {
        padding = 0x00
    } else {
        padding = 0xFF
    }

    valueLength := uint(len(value))
    
    if bits >= valueLength * 8 {
        var off uint
        for off = 0; off < valueLength; off++ {
            value[off] = byte(padding)
        }
        return
    }

    byteNum := bits >> 3
    bitNum := bits & 7
    
    // copy bits
    var i uint
    for i = valueLength - 1; int(i) >= int(byteNum); i-- {
        high := (value[i-byteNum] & (1 << (8 - bitNum) - 1)) << bitNum
        var low byte = 0
        if int(i-byteNum-1) >= 0 {
            low = (value[i-byteNum-1] & ^(1 << (8 - bitNum) - 1)) >> (8 - bitNum)
        }
        value[i] = high | low
    }
    
    // padding
    for i = 0; i < byteNum; i++ {
        value[i] = padding
    }
    if padding > 0 {
        value[byteNum] = value[byteNum] | (1 << bitNum - 1)
    } else {
        value[byteNum] = value[byteNum] & ^(1 << bitNum - 1)
    }
}

func RightShift(value []byte, bits uint, padding uint8)  {
    if bits == 0 {
        return
    }
    if padding == 0 {
        padding = 0x00
    } else {
        padding = 0xFF
    }

    valueLength := uint(len(value))
    
    if bits >= valueLength * 8 {
        var off uint
        for off = 0; off < valueLength; off++ {
            value[off] = byte(padding)
        }
        return
    }

    byteNum := bits >> 3
    bitNum := bits & 7

    // copy bits
    var i uint
    for i = 0; i < valueLength - byteNum; i++ {
        var high byte = 0
        if i+byteNum+1 < valueLength {
            high = (value[i+byteNum+1] & (1 << bitNum - 1)) << (8 - bitNum)
        }
        low := ((value[i+byteNum] & ^(1 << bitNum - 1))) >> bitNum
        value[i] = high | low
    }
    

    // padding
    for i = valueLength - 1; i >= valueLength - byteNum; i-- {
        value[i] = padding
    }
    if padding > 0 {
        value[valueLength - byteNum - 1] = value[valueLength - byteNum - 1] | ^(1 << (8 - bitNum) - 1)
    } else {
        value[valueLength - byteNum - 1] = value[valueLength - byteNum - 1] & (1 << (8 - bitNum) - 1)
    }
}

func Not(value []byte) {
    for i := 0; i < len(value); i++ {
        value[i] = ^value[i];
    }
}

func And(a []byte, b []byte) {
    for i := 0; i < len(a); i++ {
        a[i] = a[i] & b[i];
    }
}

func Or(a []byte, b []byte) {
    for i := 0; i < len(a); i++ {
        a[i] = a[i] | b[i];
    }
}

func Xor(a []byte, b []byte) {
    for i := 0; i < len(a); i++ {
        a[i] = a[i] ^ b[i];
    }
}

func Add(a []byte, b []byte) {
    var carry, nextCarry byte = 0, 0
    for i := 0; i < len(a); i++ {
        if i < len(b) && a[i] > BYTE_MAX - b[i] - carry {
            nextCarry = 1
        } else {
            nextCarry = 0
        }
        if i < len(b) {
            a[i] = a[i] + b[i] + carry
        } else {
            a[i] = a[i] + carry
        }
        carry = nextCarry
    }
}

func Sub(a []byte, b []byte) {
    newB := make([]byte, len(b))
    copy(newB, b)
    TwosComplement(newB)
    Add(a, newB)
}

func Multiply(a []byte, b []byte) {
    ans := make([]byte, len(a) + len(b))
    bits := nbits(b)

    var i uint
    for i = bits - 1; int(i) >= 0; i-- {
        byteNum := i >> 3
        bitNum := i & 7

        LeftShift(ans, 1, 0)
        if (b[byteNum] & (1 << bitNum)) > 0 {
            Add(ans, a)
        }
    }
    copy(a, ans)
}

func Divide(a []byte, b []byte, signed bool) []byte {
    quotient := make([]byte, len(a))
    remainder := make([]byte, len(a))
    copy(remainder, a)
    divider := make([]byte, len(b))
    copy(divider, b)
    
    if IsZero(b) {
        log.Fatal("Divisor cannot be zero.")
    }
    if Compare(a, b) < 0 {
        for i := 0; i < len(a); i++ {
            a[i] = 0
        }
        return remainder
    }

    var negA, negB bool
    if signed {
        negA = IsNegative(a)
        if negA {
            TwosComplement(a)
        }
        negB = IsNegative(b)
        if negB {
            TwosComplement(b)
        }
    }

    var dPadding int = 0
    var rPadding int = 0
    var count int = len(remainder) * 8

    for count > 0 {
        count--
        if (remainder[len(remainder) - 1] & 0x80) != 0 {
            break
        }
        LeftShift(remainder, 1, 0)
        rPadding++
    }

    copy(remainder, a)
    count = len(divider) * 8

    for count > 0 {
        count--
        if (divider[len(divider) - 1] & 0x80) != 0 {
            break
        }
        LeftShift(divider, 1, 0)
        dPadding++
    }
    
    RightShift(divider, uint(rPadding), 0)
    count = dPadding - rPadding + 1

    for count > 0 {
        count--

        if Compare(remainder, divider) >= 0 {
            Sub(remainder, divider)
            quotient[0] = quotient[0] | 0x01
        }
        if count > 0 {
            LeftShift(quotient, 1, 0)
            RightShift(divider, 1, 0)
        }
    }

    if negA != negB {
       TwosComplement(quotient)
    }

    copy(a, quotient)
    return remainder
}

// TODO : Mod

func Compare(a []byte, b []byte) int {
    if len(a) == 0 && len(b) == 0 {
        return 0
    }

    var valA, valB byte
    for i := max(len(a), len(b)) - 1; i >= 0; i-- {
        if i < len(a) {
            valA = a[i]
        } else {
            valA = 0
        }
        if i < len(b) {
            valB = b[i]
        } else {
            valB = 0
        }

        if valA == valB {
            continue
        }
        
        if valA > valB {
            return 1
        } else {
            return -1
        }
    }
    return 0
}

func IsZero(value []byte) bool {
    isZero := true
    for _, v := range value {
        isZero = isZero && (v == 0)
    }
    return isZero
}

func IsNegative(value []byte) bool {
    return (value[len(value) - 1] & 0x80) > 0
}

func TwosComplement(value []byte) {
    var carry, nextCarry byte = 1, 0
    for i := 0; i < len(value); i++ {
        if ^value[i] > BYTE_MAX - carry {
            nextCarry = 1
        } else {
            nextCarry = 0
        }
        value[i] = ^value[i] + carry
        carry = nextCarry
    }
}

func ToBinaryString(value []byte) string {
    str := ""
    for i := len(value) - 1; i >= 0; i-- {
        subStr := strconv.FormatUint(uint64(value[i]), 2)
        str = str + paddingZero(subStr, 8)
    }
    return str
}

func ToHexString(value []byte) string {
    str := ""
    for i := len(value) - 1; i >= 0; i-- {
        subStr := strconv.FormatUint(uint64(value[i]), 16)
        str = str + paddingZero(subStr, 2)
    }
    return str
}

// TODO : ToDecimalString(value []byte, signed bool)

func paddingZero(data string, length int) string {
    zeros := length - len(data)
    padded := ""
    for zeros > 0 {
        padded = padded + "0";
        zeros--
    }

    return padded + data;
}

func genMask(bits uint) uint8 {
    if bits > 8 {
        return 0xFF
    }

    var val uint8 = 0
    for bits > 0 {
        val = ((val << 1) | 1) >> 0
        bits--
    }
    return val;
}

func nbits(value []byte) uint {
    var byteNum, bitNum int = 0, 0
    for i := len(value) - 1; i >= 0; i-- {
        if value[i] != 0 {
            byteNum = i
            break
        }
    }
    for i := 7; i >= 0; i-- {
        if value[byteNum] & (1 << byte(i)) > 0 {
            bitNum = i
            break
        } 
    }
    bits := byteNum * 8 + bitNum + 1
    return uint(bits)
}

func max(a int, b int) int {
    if a < b {
        return b
    }
    return a
}
