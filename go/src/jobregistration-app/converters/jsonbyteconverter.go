package converters

import "encoding/json"

//JSONByteConverter simply wraps the standard json marshalling code into an interface, so
//that other formats might be substituted later.
type JSONByteConverter struct{}

//ContentType returns the default JSON accept type.
func (converter *JSONByteConverter) ContentType() string {
	return "application/json"
}

//ToBytes takes a generic object and attempts to convert it to a JSON byte string.
func (converter *JSONByteConverter) ToBytes(object interface{}) ([]byte, error) {
	return json.Marshal(object)
}

//FromBytes takes a default object, and then tries to serialize the
//it from the JSON byte string.
func (converter *JSONByteConverter) FromBytes(bytes []byte, object interface{}) error {
	return json.Unmarshal(bytes, object)
}
