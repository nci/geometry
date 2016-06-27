package geometry

import (
	"encoding/json"
	"testing"
)

func TestPointJSON(t *testing.T) {
	p := Point{X: 4.0, Y: 9.5}

	out, err := json.Marshal(p)
	if err != nil {
		t.Errorf("JSON Point Test failed, error in JSON serialisation: %s", err)
	}
	var pout Point
	err = json.Unmarshal(out, &pout)
	if pout != p {
		t.Errorf("JSON Point Test failed, expected: %+v, got: %+v", p, pout)
	}
}

func TestPointWKT(t *testing.T) {
	p := Point{X: 4.0, Y: 9.5}

	wktPoint := p.MarshalWKT()

	var pout Point
	err := pout.UnmarshalWKT(wktPoint)
	if err != nil {
		t.Errorf("WKT Point Test failed, error in WKT deserialisation: %s", err)
	}
	if pout != p {
		t.Errorf("WKT Point Test failed, expected: %+v, got: %+v", p, pout)
	}
}

func TestPointWKB(t *testing.T) {
	p := Point{X: 4.0, Y: 9.5}

	wkbPoint := p.MarshalWKB(1)

	var pout Point
	err := pout.UnmarshalWKB(wkbPoint)
	if err != nil {
		t.Errorf("WKB Point Test failed, error in WKB deserialisation: %s", err)
	}
	if pout != p {
		t.Errorf("WKB Point Test failed, expected: %+v, got: %+v", p, pout)
	}
}

func TestLinearStringJSON(t *testing.T) {
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
	for i, point := range lsout {
		if point != ls[i] {
			t.Errorf("JSON LineString Test failed, expected: %+v, got: %+v", point, ls[i])
		}
	}
}

func TestLinearStringWKT(t *testing.T) {
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
	for i, point := range lsout {
		if point != ls[i] {
			t.Errorf("WKT LineString Test failed, expected: %+v, got: %+v", point, ls[i])
		}
	}
}

func TestLinearStringWKB(t *testing.T) {
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
	for i, point := range lsout {
		if point != ls[i] {
			t.Errorf("WKB LineString Test failed, expected: %+v, got: %+v", point, ls[i])
		}
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
	for i, ring := range pout {
		for j, point := range ring {
			if point != p[i][j] {
				t.Errorf("JSON Polygon Test failed, expected: %+v, got: %+v", point, p[i][j])
			}
		}
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
	for i, ring := range pout {
		for j, point := range ring {
			if point != p[i][j] {
				t.Errorf("WKT Polygon Test failed, expected: %+v, got: %+v", point, p[i][j])
			}
		}
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
	for i, ring := range pout {
		for j, point := range ring {
			if point != p[i][j] {
				t.Errorf("WKB Polygon Test failed, expected: %+v, got: %+v, loop: %d", point, p[i][j], j)
			}
		}
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
	for i, poly := range mout {
		for j, ring := range poly {
			for k, point := range ring {
				if point != m[i][j][k] {
					t.Errorf("WKT MultiPolygon Test failed, expected: %+v, got: %+v", point)
				}
			}
		}
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
	for i, poly := range mout {
		for j, ring := range poly {
			for k, point := range ring {
				if point != m[i][j][k] {
					t.Errorf("WKT MultiPolygon Test failed, expected: %+v, got: %+v", point)
				}
			}
		}
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
	for i, poly := range mout {
		for j, ring := range poly {
			for k, point := range ring {
				if point != m[i][j][k] {
					t.Errorf("WKB MultiPolygon Test failed, expected: %+v, got: %+v", point)
				}
			}
		}
	}
}
