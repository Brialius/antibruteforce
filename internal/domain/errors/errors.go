package errors

type AntiBruteForceError string

func (a AntiBruteForceError) Error() string {
	return string(a)
}

var (
	ErrNotFound = AntiBruteForceError("record not found")
)
