package unique_alphanumeric_n_length_stateless

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const length = 7
const sampleSize = 100_000

func base62Encode(data []byte) string {
	const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bigInt := new(big.Int).SetBytes(data)
	var result strings.Builder
	for bigInt.Cmp(big.NewInt(0)) > 0 && result.Len() < length {
		mod := new(big.Int)
		bigInt.DivMod(bigInt, big.NewInt(62), mod)
		result.WriteByte(base62[mod.Int64()])
	}
	res := result.String()
	if len(res) < length {
		res += strings.Repeat("0", length-len(res))
	}
	return res[:length]
}

func generateMD5Timestamp() string {
	now := time.Now().UnixNano()
	randBytes := make([]byte, 4)
	rand.Read(randBytes) // Add some randomness
	data := append([]byte(strconv.FormatInt(now, 10)), randBytes...)
	hash := md5.Sum(data)
	return base62Encode(hash[:])[:length]
}

func generateRandom() string {
	b := make([]byte, length)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return string(b)
}

func generateUUIDv4() string {
	u := uuid.New()
	h := md5.Sum(u[:]) // Reduce entropy for base62 conversion
	return base62Encode(h[:])
}

func generateWithRandEntropy() string {
	randomPart, _ := rand.Int(rand.Reader, big.NewInt(1000))
	seed := fmt.Sprintf("%d-%d", time.Now().UnixNano(), randomPart.Int64())
	h := md5.Sum([]byte(seed))
	return base62Encode(h[:])
}
