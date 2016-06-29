package geometry

import (
	"gopkg.in/mgo.v2/bson"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Point struct {
	X float64
	Y float64
}

func (p Point) Equals(q Point) bool {
	if p == q {
		return true
	}
	return false
}

var endian map[uint8]binary.ByteOrder = map[uint8]binary.ByteOrder{0: binary.BigEndian, 1: binary.LittleEndian}

// GetBSON implements bson.Getter.
func (p Point) GetBSON() (interface{}, error) {
	return []float64{p.X, p.Y}, nil
}

// SetBSON implements bson.Setter.
func (p *Point) SetBSON(raw bson.Raw) error {
	out := make(map[string]interface{})
	bsonErr := raw.Unmarshal(&out)
	if bsonErr == nil {
		*p = Point{X: out["0"].(float64), Y: out["1"].(float64)}
		return nil
	} else {
		return bsonErr
	}
}

func (p Point) WKB(end binary.ByteOrder) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, end, &p)
	return buf.Bytes()
}

func (p Point) WKT() string {
	return fmt.Sprintf("%g%s%g", p.X, " ", p.Y)
}

func (p Point) JSON() string {
	return fmt.Sprintf("%s%g%s%g%s", "[", p.X, ",", p.Y, "]")
}

func (p Point) MarshalWKB(mode uint8) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, endian[mode], &mode)

	pointId := uint32(1)
	binary.Write(buf, endian[mode], &pointId)

	enc := p.WKB(endian[mode])
	binary.Write(buf, endian[mode], &enc)

	return buf.Bytes()
}

func (p *Point) UnmarshalWKB(in []byte) error {
	buf := bytes.NewBuffer(in)

	var end uint8
	err := binary.Read(buf, binary.BigEndian, &end)
	if err != nil {
		return fmt.Errorf("Error reading geometry: %s", err)
	}

	var wkbType uint32
	err = binary.Read(buf, endian[end], &wkbType)
	if err != nil || wkbType != 1 {
		return fmt.Errorf("Not a Point: %s", err)
	}

	point, err := ExtractWKBPoint(buf, endian[end])
	*p = *point

	return err
}

func (p Point) MarshalWKT() string {
	return fmt.Sprintf("POINT (%s)", p.WKT())
}

func (p *Point) UnmarshalWKT(in string) error {
	regExp := `^POINT\s+\((?P<point>.*)\)$`

	r := regexp.MustCompile(regExp)
	match := r.FindStringSubmatch(in)
	point, err := ExtractWKTPoint(match[1])
	*p = *point

	return err
}

func (p Point) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"type": "Point", "coordinates": %s}`, p.JSON())), nil
}

func (p *Point) UnmarshalJSON(in []byte) error {
	dict := make(map[string]interface{})
	err := json.Unmarshal(in, &dict)

	if err != nil {
		return err
	}

	pt, err := Interface2Point(dict["coordinates"])
	*p = *pt

	return err
}

func Interface2Point(a interface{}) (*Point, error) {
	if reflect.TypeOf(a).Kind() != reflect.Slice {
		return nil, errors.New("Wrong type for coordinates.")
	}

	s := reflect.ValueOf(a)

	if s.Len() != 2 {
		return nil, errors.New("Point of wrong dimension")
	}

	return &Point{s.Index(0).Interface().(float64), s.Index(1).Interface().(float64)}, nil

}

func ExtractWKBPoint(buf *bytes.Buffer, end binary.ByteOrder) (*Point, error) {
	var X, Y float64
	//err := binary.Read(buf, binary.BigEndian, &X)
	err := binary.Read(buf, end, &X)
	if err != nil {
		return nil, fmt.Errorf("Error reading: %s", err)
	}
	//err = binary.Read(buf, binary.BigEndian, &Y)
	err = binary.Read(buf, end, &Y)
	if err != nil {
		return nil, fmt.Errorf("Error reading: %s", err)
	}
	return &Point{X, Y}, nil
}

func ExtractWKTPoint(in string) (*Point, error) {
	points := strings.Split(in, " ")
	if len(points) != 2 {
		return nil, errors.New("input not recognised as WKT Point")
	}

	X, err := strconv.ParseFloat(points[0], 64)
	if err != nil {
		return nil, err
	}

	Y, err := strconv.ParseFloat(points[1], 64)
	if err != nil {
		return nil, err
	}
	return &Point{X, Y}, nil
}
