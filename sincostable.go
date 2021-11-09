package math32

type sincos struct {
	sign               [4]float32 // Sign of sine in quadrants 0..3.
	hpiInv             float32
	hpi                float32
	c0, c1, c2, c3, c4 float32 // Cosine polynomial
	s1, s2, s3         float32 // Sine polynomial
}

var sincosTable0 = &sincos{
	sign:   [4]float32{1, -1, -1, 1},
	hpiInv: 0x1.45F306DC9C883p+23, // TOINTINTRINSICS=0
	hpi:    0x1.921FB54442D18p0,
	c0:     0x1p0,
	c1:     -0x1.ffffffd0c621cp-2,
	c2:     0x1.55553e1068f19p-5,
	c3:     -0x1.6c087e89a359dp-10,
	c4:     0x1.99343027bf8c3p-16,
	s1:     -0x1.555545995a603p-3,
	s2:     0x1.1107605230bc4p-7,
	s3:     -0x1.994eb3774cf24p-13,
}

var sincosTable1 = &sincos{
	sign:   [4]float32{1, -1, -1, 1},
	hpiInv: 0x1.45F306DC9C883p+23, // TOINTINTRINSICS=0
	hpi:    0x1.921FB54442D18p0,
	c0:     -0x1p0,
	c1:     0x1.ffffffd0c621cp-2,
	c2:     -0x1.55553e1068f19p-5,
	c3:     0x1.6c087e89a359dp-10,
	c4:     -0x1.99343027bf8c3p-16,
	s1:     -0x1.555545995a603p-3,
	s2:     0x1.1107605230bc4p-7,
	s3:     -0x1.994eb3774cf24p-13,
}

// Return the sine of inputs X and X2 (X squared) using the polynomial P.
// N is the quadrant, and if odd the cosine polynomial is used.
func sinPolyS(x, x2 float32, p *sincos, n int) float32 {
	if n&1 == 0 {
		x3 := x * x2
		s1 := p.s2 + x2*p.s3

		x7 := x3 * x2
		s := x + x3*p.s1
		return s + x7*s1
	}

	x4 := x2 * x2
	c2 := p.c3 + x2*p.c4
	c1 := p.c0 + x2*p.c1

	x6 := x4 * x2
	c := c1 + x4*p.c2
	return c + x6*c2
}

func reducePi(x float32, p *sincos) (float32, int) {
	// 	Use scaled float to int conversion with explicit rounding.
	// hpi_inv is prescaled by 2^24 so the quadrant ends up in bits 24..31.
	// This avoids inaccuracies introduced by truncating negative values.
	r := x * p.hpiInv
	n := int((int32(r) + 0x800000) >> 24)
	return x - float32(n)*p.hpi, n
}

var invpio4 = [24]uint32{
	0xa2, 0xa2f9, 0xa2f983, 0xa2f9836e,
	0xf9836e4e, 0x836e4e44, 0x6e4e4415, 0x4e441529,
	0x441529fc, 0x1529fc27, 0x29fc2757, 0xfc2757d1,
	0x2757d1f5, 0x57d1f534, 0xd1f534dd, 0xf534ddc0,
	0x34ddc0db, 0xddc0db62, 0xc0db6295, 0xdb629599,
	0x6295993c, 0x95993c43, 0x993c4390, 0x3c439041,
}

func reduceLargePi(xi uint32) (float32, int) {
	const pi63 = 0x1.921FB54442D18p-62

	arr := invpio4[(xi>>26)&15:]
	shift := int((xi >> 23) & 7)
	xi = (xi & 0xffffff) | 0x800000
	xi <<= shift

	var res0, res1, res2, n uint64
	res0 = uint64(xi) * uint64(arr[0])
	res1 = uint64(xi) * uint64(arr[4])
	res2 = uint64(xi) * uint64(arr[8])
	res0 = (res2 >> 32) | (res0 << 32)
	res0 += res1

	n = (res0 + (1 << 61)) >> 62
	res0 -= n << 62
	x := int64(res0)
	return float32(x) * pi63, int(n)
}
