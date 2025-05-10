package mt19937

const (
	NN              = 312
	MM              = 156
	MATRIX_A uint64 = 0xB5026F5AA96619E9
	UM       uint64 = 0xFFFFFFFF80000000
	LM       uint64 = 0x7FFFFFFF
)

var (
	mt    [NN]uint64
	mti   int = NN + 1 // mti==NN+1 means not initialized
	mag01     = [2]uint64{0, MATRIX_A}
)

func Init_genrand64(seed uint64) {
	mt[0] = seed
	for mti = 1; mti < NN; mti++ {
		mt[mti] = (6364136223846793005*(mt[mti-1]^(mt[mti-1]>>62)) + uint64(mti))
	}
}

/* initialize by an array with array-length */
/* init_key is the array for initializing keys */
/* key_length is its length */
func Init_by_array64(init_key []uint64, key_length uint64) {
	var i, j, k uint64
	Init_genrand64(19650218)
	i = 1
	j = 0
	if NN > key_length {
		k = NN
	} else {
		k = key_length
	}
	for ; k != 0; k-- {
		mt[i] = (mt[i] ^ ((mt[i-1] ^ (mt[i-1] >> 62)) * 3935559000370003845)) + init_key[j] + j /* non linear */
		i++
		j++
		if i >= NN {
			mt[0] = mt[NN-1]
			i = 1
		}
		if j >= key_length {
			j = 0
		}
	}
	for k = NN - 1; k != 0; k-- {
		mt[i] = (mt[i] ^ ((mt[i-1] ^ (mt[i-1] >> 62)) * 2862933555777941757)) - i /* non linear */
		i++
		if i >= NN {
			mt[0] = mt[NN-1]
			i = 1
		}
	}

	mt[0] = (1 << 63) /* MSB is 1; assuring non-zero initial array */
}

func Genrand64_int64() uint64 {
	var i int = 0
	var x uint64 = 0

	if mti >= NN { /* generate NN words at one time */

		/* if init_genrand64() has not been called, */
		/* a default initial seed is used     */
		if mti == NN+1 {
			Init_genrand64(5489)
		}

		for i = 0; i < NN-MM; i++ {
			x = (mt[i] & UM) | (mt[i+1] & LM)
			mt[i] = mt[i+MM] ^ (x >> 1) ^ mag01[(int)(x&1)]
		}
		for ; i < NN-1; i++ {
			x = (mt[i] & UM) | (mt[i+1] & LM)
			mt[i] = mt[i+(MM-NN)] ^ (x >> 1) ^ mag01[(int)(x&1)]
		}
		x = (mt[NN-1] & UM) | (mt[0] & LM)
		mt[NN-1] = mt[MM-1] ^ (x >> 1) ^ mag01[(int)(x&1)]
		mti = 0
	}
	x = mt[mti]
	mti++
	x ^= (x >> 29) & 0x5555555555555555
	x ^= (x << 17) & 0x71D67FFFEDA60000
	x ^= (x << 37) & 0xFFF7EEE000000000
	x ^= (x >> 43)

	return x
}

/* generates a random number on [0, 2^63-1]-interval */
func Genrand64_int63() int64 {
	return int64((Genrand64_int64() >> 1))
}

/* generates a random number on [0,1]-real-interval */
func Genrand64_real1() float64 {
	return float64(Genrand64_int64()>>11) * (1.0 / 9007199254740991.0)
}

/* generates a random number on [0,1)-real-interval */
func Genrand64_real2() float64 {
	return float64(Genrand64_int64()>>11) * (1.0 / 9007199254740992.0)
}

/* generates a random number on (0,1)-real-interval */
func Genrand64_real3() float64 {
	return float64(Genrand64_int64()>>12) + float64(0.5)*(1.0/4503599627370496.0)
}

func InitEX() {
	var arr = []uint64{0x12345, 0x23456, 0x34567, 0x45678}
	var length uint64 = 4
	Init_by_array64(arr, length)
}
