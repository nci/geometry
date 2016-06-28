package geometry

import (
	"gopkg.in/mgo.v2/bson"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"reflect"
	"regexp"
	"strings"
)

type LineString []Point


func (l LineString) Equals(ls LineString) bool {
	for i, point := range(l) {
		if !point.Equals(ls[i]) {
			return false
		}
	}
	return true
}

func (l LineString) WKB(end binary.ByteOrder) []byte {
	buf := new(bytes.Buffer)
	numPoints := uint32(len(l))
	binary.Write(buf, end, &numPoints)
	for _, point := range l {
		binary.Write(buf, end, point.WKB(end))
	}
	return buf.Bytes()
}

func (l LineString) WKT() string {
	out := "("

	for i, point := range l {
		if i == 0 {
			out += point.WKT()
		} else {
			out += fmt.Sprintf(",%s", point.WKT())
		}
	}
	out += ")"

	return out
}

func (l LineString) JSON() string {
	out := "["

	for i, point := range l {
		if i == 0 {
			out += point.JSON()
		} else {
			out += fmt.Sprintf(",%s", point.JSON())
		}
	}
	out += "]"

	return out
}

func (l LineString) MarshalWKB(mode uint8) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, endian[mode], &mode)

	lsId := uint32(2)
	binary.Write(buf, endian[mode], &lsId)

	enc := l.WKB(endian[mode])
	binary.Write(buf, endian[mode], &enc)

	return buf.Bytes()
}

func (l *LineString) UnmarshalWKB(in []byte) error {
	buf := bytes.NewBuffer(in)

	var end uint8
	err := binary.Read(buf, binary.BigEndian, &end)
	if err != nil {
		return fmt.Errorf("Error reading geometry: %s", err)
	}

	var wkbType uint32
	//err = binary.Read(buf, binary.BigEndian, &wkbType)
	err = binary.Read(buf, endian[end], &wkbType)
	if err != nil || wkbType != 2 {
		return fmt.Errorf("Not a LineString: %s", err)
	}

	*l, err = ExtractWKBLineString(buf, endian[end])

	return err
}

func (l LineString) MarshalWKT() string {
	return fmt.Sprintf("LINESTRING %s", l.WKT())
}

func (l *LineString) UnmarshalWKT(in string) error {
	//POLYGON ((4 9.5, 2 9.5, 4 5.5, 4 9.5, 4 9.5))
	regExp := `^LINESTRING\s+(?P<points>\(.*\))$`

	r := regexp.MustCompile(regExp)
	match := r.FindStringSubmatch(in)
	var err error
	*l, err = ExtractWKTLineString(match[1])

	return err
}

func (l LineString) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"type": "LineString", "coordinates":%s}`, l.JSON())), nil
}

func (l *LineString) UnmarshalJSON(in []byte) error {
	dict := make(map[string]interface{})
	err := json.Unmarshal(in, &dict)

	if err != nil {
		return err
	}
	*l, err = Interface2LineString(dict["coordinates"])

	return err
}

type LinearRing []Point

func (r LinearRing) Equals(lr LinearRing) bool {
	for i, point := range(r) {
		if !point.Equals(lr[i]) {
			return false
		}
	}
	return true
}

// GetBSON implements bson.Getter.
func (r LinearRing) GetBSON() (interface{}, error) {
	out := make([]Point, len(r)+1)
	for i, point := range(r) {
		out[i] = point
	}
	out[len(r)] = r[0]
	return out, nil
}

// SetBSON implements bson.Setter.
func (r *LinearRing) SetBSON(raw bson.Raw) error {

	out := make(map[string]interface{})
	bsonErr := raw.Unmarshal(&out)

	if bsonErr == nil {
		aux := make([]Point, len(out)-1)
		for i:=0; i<len(aux); i++ {
			key := strconv.Itoa(i)
			aux[i] = Point{out[key].(map[string]interface{})["x"].(float64), 
				       out[key].(map[string]interface{})["y"].(float64)} 	               
		}

		*r = aux
		return nil
	} else {
		return bsonErr
	}
}

func (r LinearRing) WKB(end binary.ByteOrder) []byte {
	buf := new(bytes.Buffer)
	numPoints := uint32(len(r) + 1)
	binary.Write(buf, end, &numPoints)
	for _, point := range r {
		binary.Write(buf, end, point.WKB(end))
	}
	binary.Write(buf, end, r[0].WKB(end))
	return buf.Bytes()
}

func (r LinearRing) WKT() string {
	out := "("

	for i, point := range r {
		if i == 0 {
			out += point.WKT()
		} else {
			out += fmt.Sprintf(",%s", point.WKT())
		}
	}
	out += fmt.Sprintf(",%s)", r[0].WKT())

	return out
}

func (r LinearRing) JSON() string {
	out := "["

	for i, point := range r {
		if i == 0 {
			out += point.JSON()
		} else {
			out += fmt.Sprintf(",%s", point.JSON())
		}
	}
	out += fmt.Sprintf(",%s]", r[0].JSON())

	return out
}

func Interface2LineString(a interface{}) (LineString, error) {

	if reflect.TypeOf(a).Kind() != reflect.Slice {
		return nil, errors.New("Wrong type for coordinates.")
	}

	s := reflect.ValueOf(a)

	if s.Len() < 3 {
		return nil, errors.New("LineString of wrong dimension. Should have at least 3 Points")
	}

	l := LineString{}
	for i := 0; i < s.Len(); i++ {
		point, err := Interface2Point(s.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		l = append(l, *point)
	}

	return l, nil
}

func Interface2LinearRing(a interface{}) (LinearRing, error) {
	if reflect.TypeOf(a).Kind() != reflect.Slice {
		return nil, errors.New("Wrong type for coordinates.")
	}

	s := reflect.ValueOf(a)

	if s.Len() < 3 {
		return nil, errors.New("LinearRing of wrong dimension. Should have at least 3 Points")
	}

	r := LinearRing{}
	for i := 0; i < s.Len()-1; i++ {
		point, err := Interface2Point(s.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		r = append(r, *point)
	}

	return r, nil
}

func ExtractWKTLineString(in string) (LineString, error) {
	//LINESTRING (4 9.5, 2 9.5, 4 5.5, 4 9.5, 4 9.5)

	points := strings.Split(strings.Trim(in, "()"), ",")
	line := LineString{}
	for _, pointStr := range points[:len(points)] {
		point, _ := ExtractWKTPoint(strings.Trim(pointStr, " "))
		line = append(line, *point)
	}

	return line, nil
}

func ExtractWKTLinearRing(in string) (LinearRing, error) {
	//LINESTRING (4 9.5, 2 9.5, 4 5.5, 4 9.5, 4 9.5)

	points := strings.Split(strings.Trim(in, "()"), ",")
	ring := LinearRing{}
	for _, pointStr := range points[:len(points)-1] {
		point, _ := ExtractWKTPoint(strings.Trim(pointStr, " "))
		ring = append(ring, *point)
	}

	return ring, nil
}

func ExtractWKBLineString(buf *bytes.Buffer, end binary.ByteOrder) (LineString, error) {
	var numPoints uint32
	err := binary.Read(buf, end, &numPoints)
	if err != nil {
		return nil, err
	}

	ls := make([]Point, int(numPoints))

	for i := 0; i < int(numPoints); i++ {
		point, err := ExtractWKBPoint(buf, end)
		if err != nil {
			return nil, err
		}
		ls[i] = *point
	}

	return LineString(ls), nil
}

func ExtractWKBLinearRing(buf *bytes.Buffer, end binary.ByteOrder) (LinearRing, error) {
	var numPoints uint32
	err := binary.Read(buf, end, &numPoints)
	if err != nil {
		return nil, err
	}

	lr := make([]Point, int(numPoints)-1)

	for i := 0; i < int(numPoints); i++ {
		point, err := ExtractWKBPoint(buf, end)
		if err != nil {
			return nil, err
		}
		if i < int(numPoints)-1 {
			lr[i] = *point
		}
	}

	return LinearRing(lr), nil
}
