package math32

func Sin(x float32) float32 {
	// glibc implementations
	const pio4 = Pi / 4
	var n int
	p := sincosTable0
	switch {
	case absToP12s(x) < absToP12s(pio4):
		s := x * x
		if absToP12s(x) < absToP12s(0x1p-12) {
			// Force underflow for tiny y.
			// if absToP12s(x) < absToP12s(0x1p-126) {
			// math force eval // what does this do?
			// }
			return x
		}
		return sinPolyS(x, s, p, 0)
	case absToP12s(x) < absToP12s(120):
		x, n = reducePi(x, p)
		s := p.sign[n&3]

		if n&2 != 0 {
			p = sincosTable1
		}
		return sinPolyS(x*s, x*x, p, n)

	case absToP12s(x) < absToP12s(Inf(1)):
		xi := Float32bits(x)
		sign := xi >> 31

		x, n = reduceLargePi(xi)

		// Setup signs for sin and cos - include original sign.
		s := p.sign[(n+int(sign))&3]
		if (n+int(sign))&2 != 0 {
			p = sincosTable1
		}
		return sinPolyS(x*s, x*x, p, n)
	default:
		return NaN()
	}
}
