package errors

// AntiBruteForceError Business errors type
type AntiBruteForceError string

func (a AntiBruteForceError) Error() string {
	return string(a)
}

var (
	// ErrNotFound - record not found
	ErrNotFound = AntiBruteForceError("record not found")
	// ErrInvalidIP - invalid IP address
	ErrInvalidIP = AntiBruteForceError("invalid IP address")
	// ErrInvalidCIDR - invalid CIDR address
	ErrInvalidCIDR = AntiBruteForceError("invalid CIDR")
	// ErrBucketNotFound - bucket not found
	ErrBucketNotFound = AntiBruteForceError("bucket not found")
)
