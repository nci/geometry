package geometry

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestPointJSON(t *testing.T) {
	p := &Point{X: 4.0, Y: 9.5}

	out, err := json.Marshal(p)
	if err != nil {
		t.Errorf("JSON Point Test failed, error in JSON serialisation: %s", err)
	}

	var pout Point
	err = json.Unmarshal(out, &pout)

	if !pout.Equals(*p) {
		t.Errorf("JSON Point Test failed, expected: %+v, got: %+v", p, pout)
	}
}

func TestPointBSON(t *testing.T) {
	p := Point{X: 4.0, Y: 9.5}

	out, err := bson.Marshal(&p)
	if err != nil {
		t.Errorf("JSON Point Test failed, error in BSON serialisation: %s", err)
	}

	var pout Point
	err = bson.Unmarshal(out, &pout)

	if !pout.Equals(p) {
		t.Errorf("JSON Point Test failed, expected: %+v, got: %+v", p, pout)
	}
}

func TestPointWKT(t *testing.T) {
	p := &Point{X: 4.0, Y: 9.5}

	wktPoint := p.MarshalWKT()

	var pout Point
	err := pout.UnmarshalWKT(wktPoint)
	if err != nil {
		t.Errorf("WKT Point Test failed, error in WKT deserialisation: %s", err)
	}
	if !pout.Equals(*p) {
		t.Errorf("WKT Point Test failed, expected: %+v, got: %+v", p, pout)
	}
}

func TestPointWKB(t *testing.T) {
	p := &Point{X: 4.0, Y: 9.5}

	wkbPoint := p.MarshalWKB(1)

	var pout Point
	err := pout.UnmarshalWKB(wkbPoint)
	if err != nil {
		t.Errorf("WKB Point Test failed, error in WKB deserialisation: %s", err)
	}

	if !pout.Equals(*p) {
		t.Errorf("WKB Point Test failed, expected: %+v, got: %+v", *p, pout)
	}
}

func TestLinearRingBSON(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	ls := LinearRing{p1, p2, p3}

	out, err := bson.Marshal(ls)
	if err != nil {
		t.Errorf("JSON LineString Test failed, error in JSON serialisation: %s", err)
	}
	var lsout LinearRing
	err = bson.Unmarshal(out, &lsout)
	if !lsout.Equals(ls) {
		t.Errorf("BSON LineString Test failed, expected: %+v, got: %+v", ls, lsout)
	}
}

func TestLineStringJSON(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	ls := LineString{p1, p2, p3}

	out, err := json.Marshal(ls)
	if err != nil {
		t.Errorf("JSON LineString Test failed, error in JSON serialisation: %s", err)
	}
	var lsout LineString
	err = json.Unmarshal(out, &lsout)

	if !lsout.Equals(ls) {
		t.Errorf("JSON LineString Test failed, expected: %+v, got: %+v", ls, lsout)
	}
}

func TestLineStringWKT(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	ls := LineString{p1, p2, p3}

	wktLineString := ls.MarshalWKT()

	var lsout LineString
	err := lsout.UnmarshalWKT(wktLineString)
	if err != nil {
		t.Errorf("WKT LineString Test failed, error in WKT deserialisation: %s", err)
	}

	if !lsout.Equals(ls) {
		t.Errorf("JSON LineString Test failed, expected: %+v, got: %+v", ls, lsout)
	}
}

func TestLineStringWKB(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	ls := LineString{p1, p2, p3}

	wkbLineString := ls.MarshalWKB(1)

	var lsout LineString
	err := lsout.UnmarshalWKB(wkbLineString)
	if err != nil {
		t.Errorf("WKB LineString Test failed, error in WKT deserialisation: %s", err)
	}

	if !lsout.Equals(ls) {
		t.Errorf("JSON LineString Test failed, expected: %+v, got: %+v", ls, lsout)
	}
}

func TestPolygonJSON(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	p := Polygon{LinearRing{p1, p2, p3}}

	out, err := json.Marshal(p)
	if err != nil {
		t.Errorf("JSON Polygon Test failed, error in JSON serialisation: %s", err)
	}
	var pout Polygon
	err = json.Unmarshal(out, &pout)

	if !pout.Equals(p) {
		t.Errorf("JSON Polygon Test failed, expected: %+v, got: %+v", p, pout)
	}
}

func TestPolygonWKT(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	p := Polygon{LinearRing{p1, p2, p3}}

	wktPolygon := p.MarshalWKT()

	var pout Polygon
	err := pout.UnmarshalWKT("POLYGON ((124.825741335906 -27.4770818851779,124.76778384003 -28.3670436543958,125.79920096678 -28.4138865451055,125.848920710608 -27.5235483019135,124.825741335906 -27.4770818851779))")
	err = pout.UnmarshalWKT(wktPolygon)
	if err != nil {
		t.Errorf("WKT LineString Test failed, error in WKT deserialisation: %s", err)
	}

	if !pout.Equals(p) {
		t.Errorf("WKT Polygon Test failed, expected: %+v, got: %+v", p, pout)
	}
}

func TestPolygonWKB(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	p := Polygon{LinearRing{p1, p2, p3}}

	wkbPolygon := p.MarshalWKB(1)

	var pout Polygon
	err := pout.UnmarshalWKB(wkbPolygon)
	if err != nil {
		t.Errorf("WKB Polygon Test failed, error in WKB deserialisation: %s", err)
	}

	if !pout.Equals(p) {
		t.Errorf("WKT Polygon Test failed, expected: %+v, got: %+v", p, pout)
	}
}

func TestMultiPolygonJSON(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	p4 := Point{X: 8.0, Y: 9.5}
	p5 := Point{X: 6.0, Y: 9.5}
	p6 := Point{X: 8.0, Y: 5.5}
	m := MultiPolygon{Polygon{LinearRing{p1, p2, p3}}, Polygon{LinearRing{p4, p5, p6}}}

	out, err := json.Marshal(m)
	if err != nil {
		t.Errorf("JSON MultiPolygon Test failed, error in JSON serialisation: %s", err)
	}
	var mout MultiPolygon
	err = json.Unmarshal(out, &mout)
	if err != nil {
		t.Errorf("JSON MultiPolygon Test failed, error in JSON deserialisation: %s", err)
	}

	if !mout.Equals(m) {
		t.Errorf("WKT MultiPolygon Test failed, expected: %+v, got: %+v", m, mout)
	}
}

func TestMultiPolygonBSON(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	p4 := Point{X: 8.0, Y: 9.5}
	p5 := Point{X: 6.0, Y: 9.5}
	p6 := Point{X: 8.0, Y: 5.5}
	m := MultiPolygon{Polygon{LinearRing{p1, p2, p3}}, Polygon{LinearRing{p4, p5, p6}}}

	out, err := bson.Marshal(m)
	if err != nil {
		t.Errorf("BSON MultiPolygon Test failed, error in JSON serialisation: %s", err)
	}
	var mout MultiPolygon
	err = bson.Unmarshal(out, &mout)
	if err != nil {
		t.Errorf("JSON MultiPolygon Test failed, error in JSON deserialisation: %s", err)
	}

	if !mout.Equals(m) {
		t.Errorf("BSON MultiPolygon Test failed, expected: %+v, got: %+v", m, mout)
	}
}

func TestMultiPolygonWKT(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	p4 := Point{X: 8.0, Y: 9.5}
	p5 := Point{X: 6.0, Y: 9.5}
	p6 := Point{X: 8.0, Y: 5.5}
	m := MultiPolygon{Polygon{LinearRing{p1, p2, p3}}, Polygon{LinearRing{p4, p5, p6}}}

	wktMultiPolygon := m.MarshalWKT()

	var mout MultiPolygon
	err := mout.UnmarshalWKT(wktMultiPolygon)
	if err != nil {
		t.Errorf("WKT MultiPolygon Test failed, error in WKT deserialisation: %s", err)
	}

	if !mout.Equals(m) {
		t.Errorf("WKT MultiPolygon Test failed, expected: %+v, got: %+v", m, mout)
	}
}

func TestMultiPolygonWKB(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	p4 := Point{X: 8.0, Y: 9.5}
	p5 := Point{X: 6.0, Y: 9.5}
	p6 := Point{X: 8.0, Y: 5.5}
	m := MultiPolygon{Polygon{LinearRing{p1, p2, p3}}, Polygon{LinearRing{p4, p5, p6}}}

	wkbMultiPolygon := m.MarshalWKB(1)

	var mout MultiPolygon
	err := mout.UnmarshalWKB(wkbMultiPolygon)
	if err != nil {
		t.Errorf("WKB MultiPolygon Test failed, error in WKB deserialisation: %s", err)
	}

	if !mout.Equals(m) {
		t.Errorf("WKT MultiPolygon Test failed, expected: %+v, got: %+v", m, mout)
	}
}

func TestFeaturePoint(t *testing.T) {
	p := Point{X: 4.0, Y: 9.5}
	f := Feature{Type:"Feature",Geometry:&p}

	out, err := json.Marshal(f)
	if err != nil {
		t.Errorf("GeoJSON Feature Point Test failed, error in JSON serialisation: %s", err)
	}
	var fout Feature
	err = json.Unmarshal(out, &fout)
	if err != nil {
		t.Errorf("GeoJSON Feature Point Test failed, error in JSON deserialisation: %s", err)
	}

	out2, err := json.Marshal(fout)
	if string(out2) != string(out) {
		t.Errorf("GeoJSON Feature Point Test failed, expected: %+v, got: %+v", f, fout)
	}
}

func TestFeatureLineString(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	ls := LineString{p1, p2, p3}
	f := Feature{Type:"Feature",Geometry:&ls}

	out, err := json.Marshal(f)
	if err != nil {
		t.Errorf("GeoJSON Feature Line String Test failed, error in JSON serialisation: %s", err)
	}
	var fout Feature
	err = json.Unmarshal(out, &fout)
	if err != nil {
		t.Errorf("GeoJSON Feature Line String Test failed, error in JSON deserialisation: %s", err)
	}

	out2, err := json.Marshal(fout)
	if string(out2) != string(out) {
		t.Errorf("GeoJSON Feature Line String Test failed, expected: %+v, got: %+v", f, fout)
	}
}

func TestFeaturePolygonGeoJSON(t *testing.T) {
	p1 := Point{X: 4.0, Y: 9.5}
	p2 := Point{X: 2.0, Y: 9.5}
	p3 := Point{X: 4.0, Y: 5.5}
	p := Polygon{LinearRing{p1, p2, p3}}
	f := Feature{Type:"Feature",Geometry:&p}

	out, err := json.Marshal(f)
	if err != nil {
		t.Errorf("GeoJSON Feature Polygon Test failed, error in JSON serialisation: %s", err)
	}
	var fout Feature
	err = json.Unmarshal(out, &fout)
	if err != nil {
		t.Errorf("GeoJSON Feature Polygon Test failed, error in JSON deserialisation: %s", err)
	}

	out2, err := json.Marshal(fout)
	if string(out2) != string(out) {
		t.Errorf("GeoJSON Feature Polygon Test failed, expected: %+v, got: %+v", f, fout)
	}
}
