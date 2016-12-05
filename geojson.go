package geometry

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type Geometry interface {
	WKT() string
	WKB(binary.ByteOrder) []byte
	MarshalWKB(uint8) []byte
	UnmarshalWKB([]byte) error
	MarshalWKT() string
	UnmarshalWKT(string) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type TypeExtractor struct {
	Type     string           `json:"type"`
	Geometry *json.RawMessage `json:"geometry"`
}

type Feature struct {
	Type     string   `json:"type"`
	Geometry Geometry `json:"geometry"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

func (f *Feature) UnmarshalJSON(in []byte) error {
	featType := TypeExtractor{}
	err := json.Unmarshal(in, &featType)
	if err != nil {
		return err
	}

	geomType := TypeExtractor{}
	err = json.Unmarshal(*featType.Geometry, &geomType)
	if err != nil {
		return err
	}

	switch geomType.Type {
	case "Point":
		var point Point
		err = json.Unmarshal(*featType.Geometry, &point)
		if err != nil {
			return err
		}
		*f = Feature{Type: "Feature", Geometry: &point}

	case "LineString":
		var ls LineString
		err = json.Unmarshal(*featType.Geometry, &ls)
		if err != nil {
			return err
		}
		*f = Feature{Type: "Feature", Geometry: &ls}

	case "Polygon":
		var poly Polygon
		err = json.Unmarshal(*featType.Geometry, &poly)
		if err != nil {
			return err
		}
		*f = Feature{Type: "Feature", Geometry: &poly}
	default:
		return fmt.Errorf("json Unmarshal Feature: Geometry %s not recognised", string(*featType.Geometry))
	}

	return nil
}
