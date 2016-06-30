package geometry

import (
	"gopkg.in/mgo.v2/bson"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type LineString []Point

type LineStringView struct {
	Type  string  `json:"type" bson:"type"`
	Coords [][]float64 `json:"coordinates" bson:"coordinates"`
}

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

func (r LineString) AsArray() [][]float64 {
	out := [][]float64{}	
	
	for _, point := range r {
		out = append(out, point.AsArray())
	}

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
	lExp := LineStringView{"LineString", l.AsArray()}
	return json.Marshal(lExp)
}

func (l *LineString) UnmarshalJSON(in []byte) error {
	lView := LineStringView{}
	err := json.Unmarshal(in, &lView)

	if err != nil {
		return err
	}
	*l, err = Slice2LineString(lView.Coords)

	return err
}

type LinearRing []Point

type LinearRingView struct {
	Type  string  `json:"type" bson:"type"`
	Coords [][]float64 `json:"coordinates" bson:"coordinates"`
}

func (r LinearRing) Equals(lr LinearRing) bool {
	for i, point := range(r) {
		if !point.Equals(lr[i]) {
			return false
		}
	}
	return true
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

func (r LinearRing) AsArray() [][]float64 {

	out := [][]float64{}	
	for _, point := range r {
		out = append(out, point.AsArray())
	}
	out = append(out, r[0].AsArray())

	return out
}


func (r LinearRing) GetBSON() (interface{}, error) {
	return LinearRingView{"LinearRing", r.AsArray()}, nil
}

func (r *LinearRing) SetBSON(raw bson.Raw) error {
	rView := LinearRingView{}
	err := raw.Unmarshal(&rView)
	if err != nil {
		return err
	}

	rout, err := Slice2LinearRing(rView.Coords)
	*r = rout

	return err
}

func Slice2LineString(ffSlice [][]float64) (LineString, error) {
	if len(ffSlice) < 3 {
		return nil, errors.New("LineString of wrong dimension. Should have at least 3 Points")
	}

	ls := LineString{}
	for _, fSlice := range(ffSlice) {
		point, err := Slice2Point(fSlice)
		if err != nil {
			return nil, err
		}
		ls = append(ls, *point)
	}

	return ls, nil
}

func Slice2LinearRing(ffSlice [][]float64) (LinearRing, error) {
	if len(ffSlice) < 3 {
		return nil, errors.New("LineString of wrong dimension. Should have at least 3 Points")
	}

	r := LinearRing{}
	for _, fSlice := range(ffSlice) {
		point, err := Slice2Point(fSlice)
		if err != nil {
			return nil, err
		}
		r = append(r, *point)
	}

	return r[:len(r)-1], nil
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
