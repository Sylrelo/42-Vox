package main

import (
	"math"
)

type Vector3 struct {
	x float64
	y float64
	z float64
}

type Vector2 struct {
	x float64
	y float64
}

func (vec2 Vector2) dot(other Vector2) float64 {
	return vec2.x*other.x + vec2.y*other.y
}

func (vec3 Vector3) dot2(other Vector2) float64 {
	return vec3.x*other.x + vec3.y*other.y
}

func (vec3 Vector3) dot3(other Vector3) float64 {
	return vec3.x*other.x + vec3.y*other.y + vec3.z*other.z
}

func GetConstantVector(v int) Vector2 {
	//v is the value from the permutation table
	var h int = v & 3
	if h == 0 {
		return Vector2{1.0, 1.0}
	} else if h == 1 {
		return Vector2{-1.0, 1.0}
	} else if h == 2 {
		return Vector2{-1.0, -1.0}
	} else {
		return Vector2{1.0, -1.0}
	}
}

func Schneider(p Vector2) float64 {
	var k1 Vector2 = Vector2{23.14069263277926, 2.665144142690225}

	var res float64 = math.Cos(dot(p, k1)) * 12345.6789
	return res - math.Floor(res)
}

func Fade(t float64) float64 {
	return ((6*t-15)*t + 10) * t * t * t
}

func Lerp(t, a1, a2 float64) float64 {
	return a1 + t*(a2-a1)
}

func perlin2d(x, y float64, permutTable int) float64 {
	var X int = int(math.Floor(x)) % 256
	var Y int = int(math.Floor(y)) % 256

	var xf float64 = x - math.Floor(x)
	var yf float64 = y - math.Floor(y)

	var topRight Vector2 = Vector2{xf - 1.0, yf - 1.0}
	var topLeft Vector2 = Vector2{xf, yf - 1.0}
	var bottomRight Vector2 = Vector2{xf - 1.0, yf}
	var bottomLeft Vector2 = Vector2{xf, yf}

	//Select a value in the array for each of the 4 corners
	// Here is the lisible version, modified it to not allocated memory to specific variable which makes it slower //
	/*
		var valueTopRight int = perlinPermutations[permutTable][perlinPermutations[permutTable][X+1]+Y+1]
		var valueTopLeft int = perlinPermutations[permutTable][perlinPermutations[permutTable][X]+Y+1]
		var valueBottomRight int = perlinPermutations[permutTable][perlinPermutations[permutTable][X+1]+Y]
		var valueBottomLeft int = perlinPermutations[permutTable][perlinPermutations[permutTable][X]+Y]

		var dotTopRight float64 = topRight.dot(GetConstantVector(valueTopRight))
		var dotTopLeft float64 = topLeft.dot(GetConstantVector(valueTopLeft))
		var dotBottomRight float64 = bottomRight.dot(GetConstantVector(valueBottomRight))
		var dotBottomLeft float64 = bottomLeft.dot(GetConstantVector(valueBottomLeft))

		var u = Fade(xf)
	*/
	var v = Fade(yf)
	/*
		return Lerp(u, Lerp(v, dotBottomLeft, dotTopLeft), Lerp(v, dotBottomRight, dotTopRight))
	*/
	return Lerp(
		Fade(xf),
		Lerp(v,
			bottomLeft.dot(GetConstantVector(perlinPermutations[permutTable][perlinPermutations[permutTable][X]+Y])),
			topLeft.dot(GetConstantVector(perlinPermutations[permutTable][perlinPermutations[permutTable][X]+Y+1]))),
		Lerp(v,
			bottomRight.dot(GetConstantVector(perlinPermutations[permutTable][perlinPermutations[permutTable][X+1]+Y])),
			topRight.dot(GetConstantVector(perlinPermutations[permutTable][perlinPermutations[permutTable][X+1]+Y+1]))),
	)

}

func Noise3dPerlin(x, y, z float64, permutTable int) float64 {
	var ab float64 = perlin2d(x, y, permutTable)
	var bc float64 = perlin2d(y, z, permutTable)
	var ac float64 = perlin2d(x, z, permutTable)

	var ba float64 = perlin2d(y, x, permutTable)
	var cb float64 = perlin2d(z, y, permutTable)
	var ca float64 = perlin2d(z, x, permutTable)

	var abc float64 = ab + bc + ac + ba + cb + ca
	return abc * 0.16666666667
}

func CapC(color float64) uint8 {
	return uint8(math.Min(math.Max(color, 0), 255))
}

func TestNoise3dMinerals(x, y, z float64, permutTable int) float64 {
	x = math.Abs(x)
	y = math.Abs(y)
	z = math.Abs(z)

	var n float64 = 0.20
	var a float64 = 0.1575
	var freq float64 = 0.1

	for octave := 0; octave < 3; octave++ {
		var v float64 = a * Noise3dPerlin(float64(x)*freq, float64(y)*freq, float64(z)*freq, permutTable)
		n += v

		a *= 0.5
		freq *= 2.0
	}

	n += 1.0
	n *= 0.5

	n = math.Pow(n, 1.5)

	return n
}

func TestNoise3dClouds(x, y, z float64, permutTable int) float64 {
	x = math.Abs(x)
	y = math.Abs(y)
	z = math.Abs(z)

	var n float64 = 0.25
	var a float64 = 0.45
	var freq float64 = 0.055

	for octave := 0; octave < 5; octave++ {
		var v float64 = a * Noise3dPerlin(float64(x)*freq, float64(y)*freq, float64(z)*freq, permutTable)
		n += v

		a *= 0.5
		freq *= 2.0
	}

	n += 1.0
	n *= 0.5

	n = math.Pow(n, 1.5)

	return n
}

func dot(v1, v2 Vector2) float64 {
	return v1.x*v2.x + v1.y*v2.y
}

func Noise2dPerlin(x, y, n, a, freq float64, octave, permutTable int) float64 {
	x = math.Abs(x)
	y = math.Abs(y)

	for octave := 0; octave < 8; octave++ {
		var v float64 = a * perlin2d(float64(x)*freq, float64(y)*freq, permutTable)
		n += v

		a *= 0.5
		freq *= 2.0
	}

	n = (n + 1.0) * 0.5

	if n < 0 {
		n = 0
	} else if n > 1 {
		n = 1
	}
	n = math.Pow(n, 1.5)

	return n
}

var grad3 = [12]Vector3{
	{1, 1, 0}, {-1, 1, 0}, {1, -1, 0}, {-1, -1, 0},
	{1, 0, 1}, {-1, 0, 1}, {1, 0, -1}, {-1, 0, -1},
	{0, 1, 1}, {0, -1, 1}, {0, 1, -1}, {0, -1, -1},
}

func simplex2d(xin, yin float64, seed int) float64 {
	// Skewing and unskewing factors for 2, 3, and 4 dimensions
	var F2 = 0.5 * (math.Sqrt(3) - 1)
	var G2 = (3 - math.Sqrt(3)) * 0.166666667

	var n0, n1, n2 float64 // Noise contributions from the three corners
	// Skew the input space to determine which simplex cell we are in
	var s float64 = (xin + yin) * F2 // Hairy factor for 2D
	var i int = int(xin + s)
	var j int = int(yin + s)
	var t float64 = float64(i+j) * G2
	var x0 float64 = xin - float64(i) + t // The x, y distances from the cell origin
	var y0 float64 = yin - float64(j) + t
	// Determine which simplex we are in
	var i1, j1 float64 // Offsets for second (middle) corner of simplex in (i, j) coords
	if x0 > y0 {       // lower triangle, XY order: (0, 0) -> (1, 0) -> (1, 1)
		i1 = 1
		j1 = 0
	} else { // upper triangle, YX order: (0, 0) -> (0, 1) -> (1, 1)
		i1 = 0
		j1 = 1
	}
	// A step of (1,0) in (i,j) means a step of (1-c,-c) in (x,y), and
	// a step of (0,1) in (i,j) means a step of (-c,1-c) in (x,y), where
	// c = (3-sqrt(3))/6
	var x1 = x0 - i1 + G2 // Offsets for middle corner in (x,y)
	var y1 = y0 - j1 + G2
	var x2 = x0 - 1 + 2*G2 // Offsets for last corner in (x,y)
	var y2 = y0 - 1 + 2*G2
	// Work out the hashed gradient indices of the three simplex corners
	i &= 255
	j &= 255

	// Calculate the contribution from the three corners
	var t0 = 0.5 - x0*x0 - y0*y0
	if t0 < 0 {
		n0 = 0
	} else {
		t0 *= t0
		n0 = t0 * t0 * Seeds[seed].gradP[i+Seeds[seed].perm[j]].dot2(Vector2{x0, y0}) // (x,y) of grad3 used for 2D gradient
	}
	var t1 = 0.5 - x1*x1 - y1*y1
	if t1 < 0 {
		n1 = 0
	} else {
		t1 *= t1
		n1 = t1 * t1 * Seeds[seed].gradP[i+int(i1)+Seeds[seed].perm[j+int(j1)]].dot2(Vector2{x1, y1})
	}
	var t2 = 0.5 - x2*x2 - y2*y2
	if t2 < 0 {
		n2 = 0
	} else {
		t2 *= t2
		n2 = t2 * t2 * Seeds[seed].gradP[i+1+Seeds[seed].perm[j+1]].dot2(Vector2{x2, y2})
	}
	// Add contributions from each corner to get the final noise value
	// The result is scaled to return values in the interval [-1, 1]
	return 70 * (n0 + n1 + n2)
}

func Noise2dSimplex(x, y, n, a, freq float64, octave, seed int) float64 {
	x = math.Abs(x)
	y = math.Abs(y)

	for octave := 0; octave < 8; octave++ {
		var v float64 = a * simplex2d(float64(x)*freq, float64(y)*freq, seed)
		n += v

		a *= 0.5
		freq *= 2.0
	}

	n = (n + 1.0) * 0.5

	if n < 0 {
		n = 0
	} else if n >= 1 {
		n = 1
	}

	return n
}

var gradient3 = [36]float64{
	1, 1, 0,
	-1, 1, 0,
	1, -1, 0,

	-1, -1, 0,
	1, 0, 1,
	-1, 0, 1,

	1, 0, -1,
	-1, 0, -1,
	0, 1, 1,

	0, -1, 1,
	0, 1, -1,
	0, -1, -1}

func simplex3d(xin, yin, zin float64, seed int) float64 {
	var F3 float64 = 1.0 * 0.333333333
	var G3 float64 = 1.0 * 0.166666667
	var n0, n1, n2, n3 float64 // Noise contributions from the four corners

	// Skew the input space to determine which simplex cell we're in
	var s float64 = (xin + yin + zin) * F3 // Hairy factor for 2D
	var i = int(xin + s)
	var j = int(yin + s)
	var k = int(zin + s)

	var t = float64(i+j+k) * G3
	var x0 = xin - float64(i) + t // The x,y distances from the cell origin, unskewed.
	var y0 = yin - float64(j) + t
	var z0 = zin - float64(k) + t

	// For the 3D case, the simplex shape is a slightly irregular tetrahedron.
	// Determine which simplex we are in.
	var i1, j1, k1 int // Offsets for second corner of simplex in (i,j,k) coords
	var i2, j2, k2 int // Offsets for third corner of simplex in (i,j,k) coords
	if x0 >= y0 {
		if y0 >= z0 {
			i1 = 1
			j1 = 0
			k1 = 0
			i2 = 1
			j2 = 1
			k2 = 0
		} else if x0 >= z0 {
			i1 = 1
			j1 = 0
			k1 = 0
			i2 = 1
			j2 = 0
			k2 = 1
		} else {
			i1 = 0
			j1 = 0
			k1 = 1
			i2 = 1
			j2 = 0
			k2 = 1
		}
	} else {
		if y0 < z0 {
			i1 = 0
			j1 = 0
			k1 = 1
			i2 = 0
			j2 = 1
			k2 = 1
		} else if x0 < z0 {
			i1 = 0
			j1 = 1
			k1 = 0
			i2 = 0
			j2 = 1
			k2 = 1
		} else {
			i1 = 0
			j1 = 1
			k1 = 0
			i2 = 1
			j2 = 1
			k2 = 0
		}
	}
	// A step of (1,0,0) in (i,j,k) means a step of (1-c,-c,-c) in (x,y,z),
	// a step of (0,1,0) in (i,j,k) means a step of (-c,1-c,-c) in (x,y,z), and
	// a step of (0,0,1) in (i,j,k) means a step of (-c,-c,1-c) in (x,y,z), where
	// c = 1/6.
	var x1 = x0 - float64(i1) + G3 // Offsets for second corner
	var y1 = y0 - float64(j1) + G3
	var z1 = z0 - float64(k1) + G3

	var x2 = x0 - float64(i2) + 2*G3 // Offsets for third corner
	var y2 = y0 - float64(j2) + 2*G3
	var z2 = z0 - float64(k2) + 2*G3

	var x3 = x0 - 1 + 3*G3 // Offsets for fourth corner
	var y3 = y0 - 1 + 3*G3
	var z3 = z0 - 1 + 3*G3

	// Work out the hashed gradient indices of the four simplex corners
	i &= 255
	j &= 255
	k &= 255

	// Calculate the contribution from the four corners
	var t0 = 0.6 - x0*x0 - y0*y0 - z0*z0
	if t0 < 0 {
		n0 = 0
	} else {
		t0 *= t0
		n0 = t0 * t0 * Seeds[seed].gradP[i+Seeds[seed].perm[j+Seeds[seed].perm[k]]].dot3(Vector3{x0, y0, z0}) // (x,y) of grad3 used for 2D gradient
	}
	var t1 = 0.6 - x1*x1 - y1*y1 - z1*z1
	if t1 < 0 {
		n1 = 0
	} else {
		t1 *= t1
		n1 = t1 * t1 * Seeds[seed].gradP[i+i1+Seeds[seed].perm[j+j1+Seeds[seed].perm[k+k1]]].dot3(Vector3{x1, y1, z1})
	}
	var t2 = 0.6 - x2*x2 - y2*y2 - z2*z2
	if t2 < 0 {
		n2 = 0
	} else {
		t2 *= t2
		n2 = t2 * t2 * Seeds[seed].gradP[i+i2+Seeds[seed].perm[j+j2+Seeds[seed].perm[k+k2]]].dot3(Vector3{x2, y2, z2})
	}
	var t3 = 0.6 - x3*x3 - y3*y3 - z3*z3
	if t3 < 0 {
		n3 = 0
	} else {
		t3 *= t3
		n3 = t3 * t3 * Seeds[seed].gradP[i+1+Seeds[seed].perm[j+1+Seeds[seed].perm[k+1]]].dot3(Vector3{x3, y3, z3})
	}
	// Add contributions from each corner to get the final noise value.
	// The result is scaled to return values in the interval [-1,1].
	return 32 * (n0 + n1 + n2 + n3)
}

func Noise3dSimplexCavern(x, y, z float64, seed int) float64 {
	x = math.Abs(x)
	y = math.Abs(y)
	z = math.Abs(z)

	var n float64 = 0.25
	var a float64 = 0.45
	var freq float64 = 0.040

	for octave := 0; octave < 2; octave++ {
		var v float64 = a * simplex3d(float64(x)*freq, float64(y)*freq, float64(z)*freq, seed)
		n += v

		a *= 0.5
		freq *= 2.0
	}

	n = (n + 1.0) * 0.5

	n = math.Pow(n, 1.65)
	return n
}

func Noise3dMinerals(x, y, z float64, seed int) float64 {
	x = math.Abs(x)
	y = math.Abs(y)
	z = math.Abs(z)

	var n float64 = 0.20
	var a float64 = 0.1575
	var freq float64 = 0.075

	for octave := 0; octave < 3; octave++ {
		var v float64 = a * simplex3d(float64(x)*freq, float64(y)*freq, float64(z)*freq, seed)
		n += v

		a *= 0.5
		freq *= 2.0
	}

	n = (n + 1.0) * 0.5

	n = math.Pow(n, 1.65)

	return n
}
