package error

type Error interface {
	Row() int
	Col() int
	Code() int
	Message() string
	String() string
}
