package validation

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringNotEmpty(t *testing.T) {
	assert.NoError(t, StringNotEmpty("test"))
	assert.Error(t, StringNotEmpty(""))
	assert.True(t, errors.Is(StringNotEmpty(""), ErrEmptyString))
}

func TestStringLength(t *testing.T) {
	assert.NoError(t, StringLength("test", 3, 5))
	assert.Error(t, StringLength("test", 5, 10))
	assert.Error(t, StringLength("test", 1, 3))
	assert.NoError(t, StringLength("test", 1, 0)) // No max limit
}

func TestStringMatches(t *testing.T) {
	regex := regexp.MustCompile(`^[a-z]+$`)
	assert.NoError(t, StringMatches("test", regex, "lowercase letters only"))
	assert.Error(t, StringMatches("Test", regex, "lowercase letters only"))
	assert.Error(t, StringMatches("test123", regex, "lowercase letters only"))
}

func TestName(t *testing.T) {
	assert.NoError(t, Name("valid-name"))
	assert.NoError(t, Name("valid-name-123"))
	assert.Error(t, Name("inv"))           // Too short
	assert.Error(t, Name("Invalid_Name"))  // Invalid format
	assert.Error(t, Name("invalid-name-")) // Ends with hyphen
}

func TestEmail(t *testing.T) {
	assert.NoError(t, Email("test@example.com"))
	assert.NoError(t, Email("test.user+tag@example.co.uk"))
	assert.Error(t, Email(""))
	assert.Error(t, Email("invalid-email"))
	assert.Error(t, Email("invalid@"))
	assert.Error(t, Email("@example.com"))
}

func TestDomain(t *testing.T) {
	assert.NoError(t, Domain("example.com"))
	assert.NoError(t, Domain("sub.example.co.uk"))
	assert.Error(t, Domain(""))
	assert.Error(t, Domain("invalid"))
	assert.Error(t, Domain(".com"))
	assert.Error(t, Domain("example."))
}

func TestIP(t *testing.T) {
	assert.NoError(t, IP("192.168.1.1"))
	assert.NoError(t, IP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")) // IPv6
	assert.Error(t, IP(""))
	assert.Error(t, IP("256.256.256.256")) // Invalid IPv4
	assert.Error(t, IP("not-an-ip"))
}

func TestNumber(t *testing.T) {
	assert.NoError(t, Number("123"))
	assert.NoError(t, Number("0"))
	assert.Error(t, Number(""))
	assert.Error(t, Number("abc"))
	assert.Error(t, Number("123.45")) // Not an integer
}

func TestNumberInRange(t *testing.T) {
	assert.NoError(t, NumberInRange("5", 1, 10))
	assert.NoError(t, NumberInRange("1", 1, 10))  // Min boundary
	assert.NoError(t, NumberInRange("10", 1, 10)) // Max boundary
	assert.Error(t, NumberInRange("0", 1, 10))    // Below min
	assert.Error(t, NumberInRange("11", 1, 10))   // Above max
	assert.Error(t, NumberInRange("abc", 1, 10))  // Not a number
	assert.NoError(t, NumberInRange("100", 1, 0)) // No max limit
}

func TestUnevenNumber(t *testing.T) {
	assert.NoError(t, UnevenNumber("1"))
	assert.NoError(t, UnevenNumber("3"))
	assert.NoError(t, UnevenNumber("101"))
	assert.Error(t, UnevenNumber("2"))
	assert.Error(t, UnevenNumber("100"))
	assert.Error(t, UnevenNumber("abc"))
}

func TestFilePath(t *testing.T) {
	assert.NoError(t, FilePath("/path/to/file"))
	assert.NoError(t, FilePath("./file"))
	assert.NoError(t, FilePath("file.txt"))
	assert.Error(t, FilePath(""))
	assert.Error(t, FilePath("../path/to/file")) // Contains ..
}

func TestURL(t *testing.T) {
	assert.NoError(t, URL("https://example.com"))
	assert.NoError(t, URL("http://example.com/path?query=value"))
	assert.Error(t, URL(""))
	assert.Error(t, URL("example.com"))       // Missing protocol
	assert.Error(t, URL("ftp://example.com")) // Invalid protocol
}
