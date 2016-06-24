package geometry

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type MultiPolygon []Polygon

func (m MultiPolygon) WKB(end binary.ByteOrder) []byte {
	buf := new(bytes.Buffer)
	numPolys := uint32(len(m))
	binary.Write(buf, end, &numPolys)
	for _, p := range m {
		binary.Write(buf, end, p.WKB(end))
	}
	return buf.Bytes()
}

func (m MultiPolygon) WKT() string {
	out := "("

	for i, poly := range m {
		if i > 0 {
			out += ","
		}
		out += poly.WKT()
	}
	out += ")"

	return out
}

func (m MultiPolygon) JSON() string {
	out := "["

	for i, poly := range m {
		if i > 0 {
			out += ","
		}
		out += poly.JSON()
	}
	out += "]"

	return out
}

func (m MultiPolygon) MarshalWKB(mode uint8) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, endian[mode], &mode)

	mId := uint32(6)
	binary.Write(buf, endian[mode], &mId)

	enc := m.WKB(endian[mode])
	binary.Write(buf, endian[mode], &enc)

	return buf.Bytes()
}

func (m *MultiPolygon) UnmarshalWKB(in []byte) error {
	buf := bytes.NewBuffer(in)

	var end uint8
	err := binary.Read(buf, binary.BigEndian, &end)
	if err != nil {
		return fmt.Errorf("Problem reading geometry: %s", err)
	}

	var wkbType uint32
	//err = binary.Read(buf, binary.BigEndian, &wkbType)
	err = binary.Read(buf, binary.LittleEndian, &wkbType)
	if err != nil || wkbType != 6 {
		return fmt.Errorf("Not a MultiPolygon: %s", err)
	}

	*m, err = ExtractWKBMultiPolygon(buf, endian[end])

	return err
}

func (p MultiPolygon) MarshalWKT() string {
	return fmt.Sprintf("MULTIPOLYGON (%s)", p.WKT())
}

func (m *MultiPolygon) UnmarshalWKT(in string) error {
	//MULTIPOLYGON (((4 9.5, 2 9.5, 4 5.5, 4 9.5)), ((8 9.5, 6 9.5, 8 5.5, 8 9.5)))
	regExp := `^MULTIPOLYGON\s+\((?P<multipolygon>\(\(.*\)\))\)$`

	r := regexp.MustCompile(regExp)
	match := r.FindStringSubmatch(in)
	var err error
	*m, err = ExtractWKTMultiPolygon(match[1])

	return err
}

func (m MultiPolygon) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"type": "MultiPolygon", "coordinates":%s}`, m.JSON())), nil
}

func (m *MultiPolygon) UnmarshalJSON(in []byte) error {
	dict := make(map[string]interface{})
	err := json.Unmarshal(in, &dict)
	if err != nil {
		return err
	}
	*m, err = Interface2MultiPolygon(dict["coordinates"])

	return err
}

func Interface2MultiPolygon(a interface{}) (MultiPolygon, error) {
	if reflect.TypeOf(a).Kind() != reflect.Slice {
		return nil, errors.New("Wrong type for coordinates for Multipolygon.")
	}

	s := reflect.ValueOf(a)
	m := MultiPolygon{}
	for i := 0; i < s.Len(); i++ {
		p, err := Interface2Polygon(s.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		m = append(m, p)
	}

	return m, nil
}

func ExtractWKTMultiPolygon(in string) (MultiPolygon, error) {
	//((4 9.5, 2 9.5, 4 5.5, 4 9.5)), ((8 9.5, 6 9.5, 8 5.5, 8 9.5))
	polygons := strings.SplitAfter(in, ")),")
	m := MultiPolygon{}
	for _, polyStr := range polygons {
		p, _ := ExtractWKTPolygon(strings.Trim(polyStr, ", "))
		m = append(m, p)
	}

	return m, nil
}

func ExtractWKBMultiPolygon(buf *bytes.Buffer, end binary.ByteOrder) (MultiPolygon, error) {
	var numPolys uint32
	err := binary.Read(buf, end, &numPolys)
	if err != nil {
		return nil, err
	}

	ps := make([]Polygon, int(numPolys))

	for i := 0; i < int(numPolys); i++ {
		poly, err := ExtractWKBPolygon(buf)
		if err != nil {
			return nil, err
		}
		ps[i] = poly
	}

	return MultiPolygon(ps), nil
}
