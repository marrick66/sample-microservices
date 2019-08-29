package converters

type ByteConverter interface {
	ToBytes(object interface{}) ([]byte, error)
	FromBytes(bytes []byte, result interface{}) error
	ContentType() string
}
