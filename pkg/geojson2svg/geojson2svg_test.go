package geojson2svg_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/fapian/geojson2svg/pkg/geojson2svg"
)

const wantEmptySVG = `<svg width="400.000000" height="400.000000"></svg>`

func oneLine(s string) string {
	res := bytes.NewBufferString("")
	for _, l := range strings.Split(s, "\n") {
		fmt.Fprintf(res, strings.TrimSpace(l))
	}
	return res.String()
}

func empty(t *testing.T) {
	expected := `<svg width="400.000000" height="400.450000"></svg>`

	svg := geojson2svg.New()
	got := svg.Draw(400, 400.45)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func withAPoint(t *testing.T) {
	expected := oneLine(`
		<svg width="400.000000" height="400.000000">
			<circle cx="200.000000" cy="200.000000" r="1"/>
		</svg>
	`)

	svg := geojson2svg.New()
	if err := svg.AddGeometry(`{"type": "Point", "coordinates": [10.5,20]}`); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	got := svg.Draw(400, 400)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func withAMultiPoint(t *testing.T) {
	expected := oneLine(`
		<svg width="400.000000" height="400.000000">
			<circle cx="0.000000" cy="400.000000" r="1"/>
			<circle cx="95.238095" cy="0.000000" r="1"/>
		</svg>
	`)

	svg := geojson2svg.New()
	if err := svg.AddGeometry(`{"type": "MultiPoint", "coordinates": [[10.5,20], [20.5,62]]}`); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	got := svg.Draw(400, 400)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func withALineString(t *testing.T) {
	expected := oneLine(`
		<svg width="400.000000" height="400.000000">
			<path d="M0.000000 291.638796,400.000000 0.000000"/>
		</svg>
	`)

	svg := geojson2svg.New()
	if err := svg.AddGeometry(`{"type": "LineString", "coordinates": [[10.4,20.5], [40.3,42.3]]}`); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	got := svg.Draw(400, 400)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func withAMultiLineString(t *testing.T) {
	expected := oneLine(`
		<svg width="400.000000" height="400.000000">
			<path d="M0.000000 282.200647,387.055016 0.000000"/>
			<path d="M12.944984 269.255663,400.000000 12.944984"/>
		</svg>
	`)

	svg := geojson2svg.New()
	if err := svg.AddGeometry(`{"type": "MultiLineString", "coordinates": [[[10.4,20.5], [40.3,42.3]], [[11.4,21.5], [41.3,41.3]]]}`); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	got := svg.Draw(400, 400)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func withAPolygonWithoutHoles(t *testing.T) {
	expected := oneLine(`
		<svg width="400.000000" height="400.000000">
			<path d="M0.000000 271.651090,372.585670 0.000000,122.118380 400.000000,0.000000 271.651090 Z"/>
		</svg>
	`)

	svg := geojson2svg.New()
	if err := svg.AddGeometry(`{"type": "Polygon", "coordinates": [[[10.4,20.5], [40.3,42.3], [20.2, 10.2], [10.4,20.5]]]}`); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	got := svg.Draw(400, 400)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func withAPolygonWithHoles(t *testing.T) {
	expected := oneLine(`
		<svg width="400.000000" height="400.000000">
			<path d="M0.000000 400.000000,400.000000 400.000000,400.000000 0.000000,0.000000 0.000000,0.000000 400.000000 M80.000000 320.000000,320.000000 320.000000,320.000000 80.000000,80.000000 80.000000,80.000000 320.000000 Z"/>
		</svg>
	`)

	svg := geojson2svg.New()
	err := svg.AddGeometry(`{"type": "Polygon", "coordinates": [
		[[100.0,0.0], [101.0,0.0], [101.0,1.0], [100.0,1.0], [100.0,0.0]],
    [[100.2,0.2], [100.8,0.2], [100.8,0.8], [100.2,0.8], [100.2,0.2]]
	]}`)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	got := svg.Draw(400, 400)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func withAMultiPolygon(t *testing.T) {
	expected := oneLine(`
		<svg width="400.000000" height="400.000000">
			<path d="M0.000000 96.247241,132.008830 0.000000,43.267108 141.721854,0.000000 96.247241 Z"/>
			<path d="M395.584989 186.754967,400.000000 186.754967,400.000000 182.339956,395.584989 182.339956,395.584989 186.754967 M396.467991 185.871965,399.116998 185.871965,399.116998 183.222958,396.467991 183.222958,396.467991 185.871965 Z"/>
		</svg>
	`)

	svg := geojson2svg.New()
	err := svg.AddGeometry(`{"type": "MultiPolygon", "coordinates": [
		[
			[[10.4,20.5], [40.3,42.3], [20.2, 10.2], [10.4,20.5]]
		], [
			[[100.0,0.0], [101.0,0.0], [101.0,1.0], [100.0,1.0], [100.0,0.0]],
	    [[100.2,0.2], [100.8,0.2], [100.8,0.8], [100.2,0.8], [100.2,0.2]]
		]
	]}`)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	got := svg.Draw(400, 400)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func withAGeometryCollection(t *testing.T) {
	expected := oneLine(`
		<svg width="400.000000" height="400.000000">
			<path d="M0.000000 291.638796,400.000000 0.000000"/>
			<circle cx="1.337793" cy="298.327759" r="1"/>
		</svg>
	`)

	svg := geojson2svg.New()
	err := svg.AddGeometry(`{"type": "GeometryCollection", "geometries": [
		{"type": "LineString", "coordinates": [[10.4,20.5], [40.3,42.3]]},
		{"type": "Point", "coordinates": [10.5,20]}
	]}`)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	got := svg.Draw(400, 400)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func withMultipleGeometries(t *testing.T) {
	expected := oneLine(`
		<svg width="400.000000" height="400.000000">
			<path d="M0.000000 291.638796,400.000000 0.000000"/>
			<circle cx="1.337793" cy="298.327759" r="1"/>
		</svg>
	`)

	svg := geojson2svg.New()
	err := svg.AddGeometry(`{"type": "LineString", "coordinates": [[10.4,20.5], [40.3,42.3]]}`)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	err = svg.AddGeometry(`{"type": "Point", "coordinates": [10.5,20]}`)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	got := svg.Draw(400, 400)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func withAnInvalidGeometry(t *testing.T) {
	geometry := `"type": "Point", "coordinates": [10.5,20]}`
	expected := "invalid geometry: " + geometry

	svg := geojson2svg.New()
	if err := svg.AddGeometry(geometry); err == nil || expected != err.Error() {
		t.Errorf("expected '%s', got %v", expected, err)
	}
	got := svg.Draw(400, 400)
	if got != wantEmptySVG {
		t.Errorf(`expected %s, got %s`, wantEmptySVG, got)
	}
}

func withAFeature(t *testing.T) {
	expected := oneLine(`
		<svg width="400.000000" height="400.000000">
			<circle cx="200.000000" cy="200.000000" r="1"/>
		</svg>
	`)

	svg := geojson2svg.New()
	err := svg.AddFeature(`{"type": "Feature", "geometry": {
		"type": "Point",
		"coordinates": [10.5,20]
	}}`)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	got := svg.Draw(400, 400)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func withAnInvalidFeature(t *testing.T) {
	feature := `{"type": "Feature", "geometry": {
		type": "Point",
		"coordinates": [10.5,20]
	}}`
	expected := "invalid feature: " + feature

	svg := geojson2svg.New()
	if err := svg.AddFeature(feature); err == nil || err.Error() != expected {
		t.Errorf("expected %s, got %v", expected, err)
	}
	got := svg.Draw(400, 400)
	if got != wantEmptySVG {
		t.Errorf(`expected %s, got %s`, wantEmptySVG, got)
	}
}

func withAFeatureCollection(t *testing.T) {
	expected := oneLine(`
		<svg width="400.000000" height="400.000000">
			<circle cx="1.337793" cy="298.327759" r="1"/>
			<path d="M0.000000 291.638796,400.000000 0.000000"/>
		</svg>
	`)

	svg := geojson2svg.New()
	err := svg.AddFeatureCollection(`{"type": "FeatureCollection", "features": [
		{"type": "Feature", "geometry": {
			"type": "Point",
			"coordinates": [10.5,20]
		}},
		{"type": "Feature", "geometry": {
			"type": "LineString",
			"coordinates": [[10.4,20.5], [40.3,42.3]]
		}}
	]}`)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	got := svg.Draw(400, 400)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func withAnInvalidFeatureCollection(t *testing.T) {
	featureCollection := `{"type": "FeatureCollection", "features": [
		{"type": "Feature", "geometry": {
			"type": "Point",
			"coordinates": [10.5,20]
		}}
		{"type": "Feature", "geometry": {
			"type": "LineString",
			"coordinates": [[10.4,20.5], [40.3,42.3]]
		}}
	]}`
	expected := "invalid feature collection: " + featureCollection

	svg := geojson2svg.New()
	if err := svg.AddFeatureCollection(featureCollection); err == nil || err.Error() != expected {
		t.Errorf("expected %s, got %v", expected, err)
	}
	got := svg.Draw(400, 400)
	if got != wantEmptySVG {
		t.Errorf(`expected %s, got %s`, wantEmptySVG, got)
	}
}

func TestSVG(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{"empty svg", empty},
		{"svg with a point", withAPoint},
		{"svg with a multipoint", withAMultiPoint},
		{"svg with a linestring", withALineString},
		{"svg with a multilinestring", withAMultiLineString},
		{"svg with a polygon without holes", withAPolygonWithoutHoles},
		{"svg with a polygon with holes", withAPolygonWithHoles},
		{"svg with a multipolygon", withAMultiPolygon},
		{"svg with a geometry collection", withAGeometryCollection},
		{"svg with multiple geometries", withMultipleGeometries},
		{"svg with an invalid geometry", withAnInvalidGeometry},
		{"svg with a feature", withAFeature},
		{"svg with an invalid feature", withAnInvalidFeature},
		{"svg with a feature collection", withAFeatureCollection},
		{"svg with an invalid feature collection", withAnInvalidFeatureCollection},
	}

	for _, tc := range tcs {
		t.Run(tc.name, tc.test)
	}
}

func TestSVGAttributeOptions(t *testing.T) {
	tcs := []struct {
		name string
		fn   func(*testing.T)
	}{
		{"should add the passed attribute to the svg tag", withAttributeOption},
		{"should add the passed attributes to the svg tag", withAttributesOption},
		{"latest attribute wins", withAttributeMultipleTimesOption},
		{"no attributes are lost", withAttributesNothingIsLostOption},
	}

	for _, tc := range tcs {
		t.Run(tc.name, tc.fn)
	}
}

func withAttributeOption(t *testing.T) {
	want := `<svg width="200.000000" height="200.000000" class="a_class" id="the_id"></svg>`
	svg := geojson2svg.New()
	got := svg.Draw(200, 200,
		geojson2svg.WithAttribute("id", "the_id"),
		geojson2svg.WithAttribute("class", "a_class"))

	if got != want {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

func withAttributeMultipleTimesOption(t *testing.T) {
	want := `<svg width="200.000000" height="200.000000" class="a_class_2" id="the_id_2"></svg>`
	svg := geojson2svg.New()
	got := svg.Draw(200, 200,
		geojson2svg.WithAttribute("id", "the_id"),
		geojson2svg.WithAttribute("class", "a_class"),
		geojson2svg.WithAttribute("class", "a_class_2"),
		geojson2svg.WithAttribute("id", "the_id_2"))

	if got != want {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

func withAttributesOption(t *testing.T) {
	want := `<svg width="200.000000" height="200.000000" class="a_class" id="the_id"></svg>`

	attributes := map[string]string{
		"id":    "the_id",
		"class": "a_class",
	}

	svg := geojson2svg.New()
	got := svg.Draw(200, 200, geojson2svg.WithAttributes(attributes))

	if got != want {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

func withAttributesNothingIsLostOption(t *testing.T) {
	want := `<svg width="200.000000" height="200.000000" class="a_class_2" id="the_id"></svg>`

	attributesA := map[string]string{"id": "the_id", "class": "a_class"}
	attributesB := map[string]string{"class": "a_class_2"}

	svg := geojson2svg.New()
	got := svg.Draw(200, 200,
		geojson2svg.WithAttributes(attributesA),
		geojson2svg.WithAttributes(attributesB))

	if got != want {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

func TestSVGPaddingOption(t *testing.T) {
	tcs := []struct {
		name     string
		data     string
		padding  geojson2svg.Padding
		expected string
	}{
		{"without padding",
			"[[0,0], [0,400], [400,400], [400,0]]",
			geojson2svg.Padding{Top: 0, Right: 0, Bottom: 0, Left: 0},
			`<svg width="200.000000" height="200.000000"><path d="M0.000000 200.000000,0.000000 0.000000,200.000000 0.000000,200.000000 200.000000"/></svg>`},
		{"with padding",
			"[[0,0], [0,400], [400,400], [400,0]]",
			geojson2svg.Padding{Top: 5, Right: 5, Bottom: 5, Left: 5},
			`<svg width="200.000000" height="200.000000"><path d="M5.000000 195.000000,5.000000 5.000000,195.000000 5.000000,195.000000 195.000000"/></svg>`},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			svg := geojson2svg.New()
			err := svg.AddGeometry(fmt.Sprintf(`{"type": "LineString", "coordinates": %s}`, tc.data))
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
			padding := geojson2svg.WithPadding(tc.padding)
			got := svg.Draw(200, 200, padding)
			if got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestFeatureProperties(t *testing.T) {
	tcs := []struct {
		name      string
		feature   string
		usedProps []string
		expected  string
	}{
		{"no props (point)",
			`{"type": "Feature", "geometry": { "type": "Point", "coordinates": [10.5,20] }}`,
			nil,
			`<svg width="400.000000" height="400.000000"><circle cx="200.000000" cy="200.000000" r="1"/></svg>`},
		{"with class (point)",
			`{"type": "Feature", "properties": {"class": "class"}, "geometry": { "type": "Point", "coordinates": [10.5,20] }}`,
			nil,
			`<svg width="400.000000" height="400.000000"><circle cx="200.000000" cy="200.000000" r="1" class="class"/></svg>`},
		{"with class and unused (point)",
			`{"type": "Feature", "properties": {"class": "class", "style": "stroke:1"}, "geometry": { "type": "Point", "coordinates": [10.5,20] }}`,
			nil,
			`<svg width="400.000000" height="400.000000"><circle cx="200.000000" cy="200.000000" r="1" class="class"/></svg>`},
		{"with unused (point)",
			`{"type": "Feature", "properties": {"style": "stroke:1"}, "geometry": { "type": "Point", "coordinates": [10.5,20] }}`,
			nil,
			`<svg width="400.000000" height="400.000000"><circle cx="200.000000" cy="200.000000" r="1"/></svg>`},
		{"with added props (point)",
			`{"type": "Feature", "properties": {"style": "stroke:1"}, "geometry": { "type": "Point", "coordinates": [10.5,20] }}`,
			[]string{"style"},
			`<svg width="400.000000" height="400.000000"><circle cx="200.000000" cy="200.000000" r="1" style="stroke:1"/></svg>`},
		{"with class removed (point)",
			`{"type": "Feature", "properties": {"class": "class"}, "geometry": { "type": "Point", "coordinates": [10.5,20] }}`,
			[]string{},
			`<svg width="400.000000" height="400.000000"><circle cx="200.000000" cy="200.000000" r="1"/></svg>`},

		{"no props (linestring)",
			`{"type": "Feature", "geometry": { "type": "LineString", "coordinates": [[10.4,20.5], [40.3,42.3]] }}`,
			nil,
			`<svg width="400.000000" height="400.000000"><path d="M0.000000 291.638796,400.000000 0.000000"/></svg>`},
		{"with class (linestring)",
			`{"type": "Feature", "properties": {"class": "class"}, "geometry": { "type": "LineString", "coordinates": [[10.4,20.5], [40.3,42.3]] }}`,
			nil,
			`<svg width="400.000000" height="400.000000"><path d="M0.000000 291.638796,400.000000 0.000000" class="class"/></svg>`},
		{"with class and unused (linestring)",
			`{"type": "Feature", "properties": {"class": "class", "style": "stroke:1"}, "geometry": { "type": "LineString", "coordinates": [[10.4,20.5], [40.3,42.3]] }}`,
			nil,
			`<svg width="400.000000" height="400.000000"><path d="M0.000000 291.638796,400.000000 0.000000" class="class"/></svg>`},
		{"with unused (linestring)",
			`{"type": "Feature", "properties": {"style": "stroke:1"}, "geometry": { "type": "LineString", "coordinates": [[10.4,20.5], [40.3,42.3]] }}`,
			nil,
			`<svg width="400.000000" height="400.000000"><path d="M0.000000 291.638796,400.000000 0.000000"/></svg>`},
		{"with added props (linestring)",
			`{"type": "Feature", "properties": {"style": "stroke:1"}, "geometry": { "type": "LineString", "coordinates": [[10.4,20.5], [40.3,42.3]] }}`,
			[]string{"style"},
			`<svg width="400.000000" height="400.000000"><path d="M0.000000 291.638796,400.000000 0.000000" style="stroke:1"/></svg>`},
		{"with class removed (linestring)",
			`{"type": "Feature", "properties": {"class": "class"}, "geometry": { "type": "LineString", "coordinates": [[10.4,20.5], [40.3,42.3]] }}`,
			[]string{},
			`<svg width="400.000000" height="400.000000"><path d="M0.000000 291.638796,400.000000 0.000000"/></svg>`},

		{"no props (polygon)",
			`{"type": "Feature", "geometry": { "type": "Polygon", "coordinates": [[[10.4,20.5], [40.3,42.3], [20.2, 10.2], [10.4,20.5]]] }}`,
			nil,
			`<svg width="400.000000" height="400.000000"><path d="M0.000000 271.651090,372.585670 0.000000,122.118380 400.000000,0.000000 271.651090 Z"/></svg>`},
		{"with class (polygon)",
			`{"type": "Feature", "properties": {"class": "class"}, "geometry": { "type": "Polygon", "coordinates": [[[10.4,20.5], [40.3,42.3], [20.2, 10.2], [10.4,20.5]]] }}`,
			nil,
			`<svg width="400.000000" height="400.000000"><path d="M0.000000 271.651090,372.585670 0.000000,122.118380 400.000000,0.000000 271.651090 Z" class="class"/></svg>`},
		{"with class and unused (polygon)",
			`{"type": "Feature", "properties": {"class": "class", "style": "stroke:1"}, "geometry": { "type": "Polygon", "coordinates": [[[10.4,20.5], [40.3,42.3], [20.2, 10.2], [10.4,20.5]]] }}`,
			nil,
			`<svg width="400.000000" height="400.000000"><path d="M0.000000 271.651090,372.585670 0.000000,122.118380 400.000000,0.000000 271.651090 Z" class="class"/></svg>`},
		{"with unused (polygon)",
			`{"type": "Feature", "properties": {"style": "stroke:1"}, "geometry": { "type": "Polygon", "coordinates": [[[10.4,20.5], [40.3,42.3], [20.2, 10.2], [10.4,20.5]]] }}`,
			nil,
			`<svg width="400.000000" height="400.000000"><path d="M0.000000 271.651090,372.585670 0.000000,122.118380 400.000000,0.000000 271.651090 Z"/></svg>`},
		{"with added props (polygon)",
			`{"type": "Feature", "properties": {"style": "stroke:1"}, "geometry": { "type": "Polygon", "coordinates": [[[10.4,20.5], [40.3,42.3], [20.2, 10.2], [10.4,20.5]]] }}`,
			[]string{"style"},
			`<svg width="400.000000" height="400.000000"><path d="M0.000000 271.651090,372.585670 0.000000,122.118380 400.000000,0.000000 271.651090 Z" style="stroke:1"/></svg>`},
		{"with class removed (polygon)",
			`{"type": "Feature", "properties": {"class": "class"}, "geometry": { "type": "Polygon", "coordinates": [[[10.4,20.5], [40.3,42.3], [20.2, 10.2], [10.4,20.5]]] }}`,
			[]string{},
			`<svg width="400.000000" height="400.000000"><path d="M0.000000 271.651090,372.585670 0.000000,122.118380 400.000000,0.000000 271.651090 Z"/></svg>`},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			svg := geojson2svg.New()
			err := svg.AddFeature(tc.feature)

			if err != nil {
				tt.Errorf("unexpected error %v", err)
			}

			var got string
			if tc.usedProps != nil {
				got = svg.Draw(400, 400, geojson2svg.UseProperties(tc.usedProps))
			} else {
				got = svg.Draw(400, 400)
			}
			if got != tc.expected {
				tt.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestExample(t *testing.T) {
	exampleFile := path.Join("..", "..", "test", "example.json")
	geojson, err := ioutil.ReadFile(exampleFile)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	svgFile := path.Join("..", "..", "test", "example.svg")
	want, err := ioutil.ReadFile(svgFile)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	svg := geojson2svg.New()
	err = svg.AddFeatureCollection(string(geojson))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	got := svg.Draw(1000, 510,
		geojson2svg.WithAttribute("xmlns", "http://www.w3.org/2000/svg"),
		geojson2svg.UseProperties([]string{"style"}),
		geojson2svg.WithPadding(geojson2svg.Padding{
			Top:    10,
			Right:  10,
			Bottom: 10,
			Left:   10,
		}))
	if got != string(want) {
		t.Errorf("expected %s, got %s", string(want), got)
	}
}
