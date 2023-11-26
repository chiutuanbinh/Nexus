package common

// This will work with numberator in little endian order
func Modulo(numerator []byte, denominator int) int {
	//a % b = ((a // 256) % b) * (256 % b) + (a % 256) % b
	result := 0
	for i := len(numerator) - 1; i >= 0; i-- {
		result *= (256 % denominator)
		result %= denominator
		result += int((numerator[i] % byte(denominator)))
		result %= denominator
	}
	return result
}
