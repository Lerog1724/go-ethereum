package arith384

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular,
// there is no security guarantees such as constant time implementation
// or side-channel attack resistance
// /!\ WARNING /!\

import (
    "math/bits"
)

type Element [4]uint64

// Limbs number of 64 bits words needed to represent Element
const Limbs = 4

// Bytes number bytes needed to represent Element
const Bytes = Limbs * 8

// Mul z = x * y mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func MulMod(z, x, y, mod *Element, modinv uint64) {
    _mulGeneric(z, x, y, mod, modinv)
}

// Add z = x + y mod q
func AddMod(z, x, y, mod *Element) {
    _addGeneric(z, x, y, mod)
}

// Sub  z = x - y mod q
func SubMod(z, x, y, mod *Element) {
    _subGeneric(z, x, y, mod)
}

// Generic (no ADX instructions, no AMD64) versions of multiplication and squaring algorithms

func _mulGeneric(z, x, y, mod *Element, modinv uint64) {

    var t [4]uint64
    var c [3]uint64
    {
        // round 0
        v := x[0]
        c[1], c[0] = bits.Mul64(v, y[0])
        m := c[0] * modinv
        c[2] = madd0(m, mod[0], c[0])
        c[1], c[0] = madd1(v, y[1], c[1])
        c[2], t[0] = madd2(m, mod[1], c[2], c[0])
        c[1], c[0] = madd1(v, y[2], c[1])
        c[2], t[1] = madd2(m, mod[2], c[2], c[0])
        c[1], c[0] = madd1(v, y[3], c[1])
        t[3], t[2] = madd3(m, mod[3], c[0], c[2], c[1])
    }
    {
        // round 1
        v := x[1]
        c[1], c[0] = madd1(v, y[0], t[0])
        m := c[0] * modinv
        c[2] = madd0(m, mod[0], c[0])
        c[1], c[0] = madd2(v, y[1], c[1], t[1])
        c[2], t[0] = madd2(m, mod[1], c[2], c[0])
        c[1], c[0] = madd2(v, y[2], c[1], t[2])
        c[2], t[1] = madd2(m, mod[2], c[2], c[0])
        c[1], c[0] = madd2(v, y[3], c[1], t[3])
        t[3], t[2] = madd3(m, mod[3], c[0], c[2], c[1])
    }
    {
        // round 2
        v := x[2]
        c[1], c[0] = madd1(v, y[0], t[0])
        m := c[0] * modinv
        c[2] = madd0(m, mod[0], c[0])
        c[1], c[0] = madd2(v, y[1], c[1], t[1])
        c[2], t[0] = madd2(m, mod[1], c[2], c[0])
        c[1], c[0] = madd2(v, y[2], c[1], t[2])
        c[2], t[1] = madd2(m, mod[2], c[2], c[0])
        c[1], c[0] = madd2(v, y[3], c[1], t[3])
        t[3], t[2] = madd3(m, mod[3], c[0], c[2], c[1])
    }
    {
        // round 3
        v := x[3]
        c[1], c[0] = madd1(v, y[0], t[0])
        m := c[0] * modinv
        c[2] = madd0(m, mod[0], c[0])
        c[1], c[0] = madd2(v, y[1], c[1], t[1])
        c[2], z[0] = madd2(m, mod[1], c[2], c[0])
        c[1], c[0] = madd2(v, y[2], c[1], t[2])
        c[2], z[1] = madd2(m, mod[2], c[2], c[0])
        c[1], c[0] = madd2(v, y[3], c[1], t[3])
        z[3], z[2] = madd3(m, mod[3], c[0], c[2], c[1])
    }

    // TODO can make the following faster and constant time

    // if z > q --> z -= q
    // note: this is NOT constant time
    if !(z[3] < mod[3] || (z[3] == mod[3] && (z[2] < mod[2] || (z[2] == mod[2] && (z[1] < mod[1] || (z[1] == mod[1] && (z[0] < mod[0] || (z[0] == mod[0] && (z[0] < mod[0]))))))))) {
        var b uint64
        z[0], b = bits.Sub64(z[0], mod[0], 0)
        z[0], b = bits.Sub64(z[0], mod[0], b)
        z[1], b = bits.Sub64(z[1], mod[1], b)
        z[2], b = bits.Sub64(z[2], mod[2], b)
        z[3], _ = bits.Sub64(z[3], mod[3], b)
    }
}

func _addGeneric(z, x, y, mod *Element) {
    var carry uint64

    z[0], carry = bits.Add64(x[0], y[0], 0)
    z[1], carry = bits.Add64(x[1], y[1], carry)
    z[2], carry = bits.Add64(x[2], y[2], carry)
    z[3], _ = bits.Add64(x[3], y[3], carry)

    // if z > q --> z -= q
    // note: this is NOT constant time
    if !(z[3] < mod[3] || (z[3] == mod[3] && (z[2] < mod[2] || (z[2] == mod[2] && (z[1] < mod[1] || (z[1] == mod[1] && (z[0] < mod[0] || (z[0] == mod[0] && (z[0] < mod[0]))))))))) {
        var b uint64
        z[0], b = bits.Sub64(z[0], mod[0], 0)
        z[0], b = bits.Sub64(z[0], mod[0], b)
        z[1], b = bits.Sub64(z[1], mod[1], b)
        z[2], b = bits.Sub64(z[2], mod[2], b)
        z[3], _ = bits.Sub64(z[3], mod[3], b)
    }
}

func _subGeneric(z, x, y, mod *Element) {
    var b uint64
    z[0], b = bits.Sub64(x[0], y[0], 0)
    z[0], b = bits.Sub64(x[0], y[0], b)
    z[1], b = bits.Sub64(x[1], y[1], b)
    z[2], b = bits.Sub64(x[2], y[2], b)
    z[3], b = bits.Sub64(x[3], y[3], b)
    if b != 0 {
        var c uint64
        z[0], c = bits.Add64(z[0], mod[0], 0)
        z[0], c = bits.Add64(z[0], mod[0], c)
        z[1], c = bits.Add64(z[1], mod[1], c)
        z[2], c = bits.Add64(z[2], mod[2], c)
        z[3], _ = bits.Add64(z[3], mod[3], c)
    }
}
