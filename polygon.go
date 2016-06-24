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

type Polygon []LinearRing

func (p Polygon) WKB(end binary.ByteOrder) []byte {
	buf := new(bytes.Buffer)
	var enc uint8
	if end == binary.BigEndian {
		enc = 0
	} else {
		enc = 1
	}

	binary.Write(buf, end, &enc)

	pId := uint32(3)
	binary.Write(buf, end, &pId)

	numRings := uint32(len(p))
	binary.Write(buf, end, &numRings)
	for _, lr := range p {
		binary.Write(buf, end, lr.WKB(end))
	}
	return buf.Bytes()
}

func (p Polygon) WKT() string {
	out := "("

	for i, ring := range p {
		if i == 0 {
			out += ring.WKT()
		} else {
			out += fmt.Sprintf(",%s", ring.WKT())
		}
	}
	out += ")"

	return out
}

func (p Polygon) JSON() string {
	out := "["

	for i, ring := range p {
		if i == 0 {
			out += ring.JSON()
		} else {
			out += fmt.Sprintf(",%s", ring.JSON())
		}
	}
	out += "]"

	return out
}

func (p Polygon) MarshalWKB(mode uint8) []byte {
	buf := new(bytes.Buffer)

	enc := p.WKB(endian[mode])
	binary.Write(buf, endian[mode], &enc)

	return buf.Bytes()
}

func (p *Polygon) UnmarshalWKB(in []byte) error {
	buf := bytes.NewBuffer(in)

	var err error
	*p, err = ExtractWKBPolygon(buf)

	return err
}

func (p Polygon) MarshalWKT() string {
	return fmt.Sprintf("POLYGON (%s)", p.WKT())
}

func (p *Polygon) UnmarshalWKT(in string) error {
	//POLYGON ((4 9.5, 2 9.5, 4 5.5, 4 9.5, 4 9.5))
	regExp := `^POLYGON\s+(?P<points>\(\(.*\)\))$`

	r := regexp.MustCompile(regExp)
	match := r.FindStringSubmatch(in)
	var err error
	*p, err = ExtractWKTPolygon(match[1])

	return err
}

func (p Polygon) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"type": "Polygon", "coordinates":%s}`, p.JSON())), nil
}

func (p *Polygon) UnmarshalJSON(in []byte) error {
	dict := make(map[string]interface{})
	err := json.Unmarshal(in, &dict)

	if err != nil {
		return err
	}
	*p, err = Interface2Polygon(dict["coordinates"])

	return err
}

func Interface2Polygon(a interface{}) (Polygon, error) {
	if reflect.TypeOf(a).Kind() != reflect.Slice {
		return nil, errors.New("Wrong type for coordinates for Polygon.")
	}

	s := reflect.ValueOf(a)

	p := Polygon{}
	for i := 0; i < s.Len(); i++ {
		ring, err := Interface2LinearRing(s.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		p = append(p, ring)
	}

	return p, nil
}

func ExtractWKTPolygon(in string) (Polygon, error) {
	//POLYGON ((4 9.5, 2 9.5, 4 5.5, 4 9.5, 4 9.5))

	rings := strings.Split(strings.Trim(in, "()"), ",")
	p := Polygon{}
	for _, pointStr := range rings[:len(rings)] {
		ring, _ := ExtractWKTLinearRing(strings.Trim(pointStr, " "))
		p = append(p, ring)
	}

	return p, nil
}

func ExtractWKBPolygon(buf *bytes.Buffer) (Polygon, error) {
	//var bigEndian uint8
	var end uint8
	err := binary.Read(buf, binary.BigEndian, &end)
	if err != nil {
		return nil, fmt.Errorf("Problem reading geometry: %s", err)
	}

	var wkbType uint32
	//err = binary.Read(buf, binary.BigEndian, &wkbType)
	err = binary.Read(buf, endian[end], &wkbType)
	if err != nil || wkbType != 3 {
		return nil, fmt.Errorf("Not a Polygon: %s", err)
	}

	var numRings uint32
	//err := binary.Read(buf, binary.BigEndian, &numRings)
	err = binary.Read(buf, endian[end], &numRings)
	if err != nil {
		return nil, err
	}
	rs := make([]LinearRing, int(numRings))

	for i := 0; i < int(numRings); i++ {
		ring, err := ExtractWKBLinearRing(buf, endian[end])
		if err != nil {
			return nil, err
		}
		rs[i] = ring
	}

	return Polygon(rs), nil
}
