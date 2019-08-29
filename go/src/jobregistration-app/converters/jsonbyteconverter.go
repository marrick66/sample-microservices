package converters

import "json"

//JsonByteConverter simply wraps the standard json marshalling code into an interface, so
//that other formats might be substituted later.
type JsonByteConverter struct { }

func (converter *JsonByteConverter) ToBytes(object interface{}) []byte, error {
	return Json.Marshal(object)
}

func (converter *JsonByteConverter) FromBytes(bytes []byte, object interface{}) error {
	return Json.Unmarshal(bytes, object)
}