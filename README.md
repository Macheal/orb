orb [![CI](https://github.com/macheal/orb/workflows/CI/badge.svg)](https://github.com/macheal/orb/actions?query=workflow%3ACI+event%3Apush) [![codecov](https://codecov.io/gh/paulmach/orb/branch/master/graph/badge.svg?token=NuuTjLVpKW)](https://codecov.io/gh/paulmach/orb) [![Go Report Card](https://goreportcard.com/badge/github.com/macheal/orb)](https://goreportcard.com/report/github.com/macheal/orb) [![Go Reference](https://pkg.go.dev/badge/github.com/macheal/orb.svg)](https://pkg.go.dev/github.com/macheal/orb)
======

Package `orb` defines a set of types for working with 2d geo and planar/projected geometric data in Golang.
There are a set of sub-packages that use these types to do interesting things.
They each provider their own README with extra info.

## Interesting features

* **Simple types** - allow for natural operations using the `make`, `append`, `len`, `[s:e]` builtins.
* **GeoJSON** - support as part of the [`geojson`](geojson) sub-package.
* **Mapbox Vector Tile** - encoding and decoding as part of the [`encoding/mvt`](encoding/mvt) sub-package.
* **Direct to type from DB query results** - by scanning WKB data directly into types.
* **Rich set of sub-packages** - including [`clipping`](clip), [`simplifing`](simplify), [`quadtree`](quadtree) and more.

## Type definitions

	type Point [2]float64
	type MultiPoint []Point

	type LineString []Point
	type MultiLineString []LineString

	type Ring LineString
	type Polygon []Ring
	type MultiPolygon []Polygon

	type Collection []Geometry

	type Bound struct { Min, Max Point }

Defining the types as slices allows them to be accessed in an idiomatic way
using Go's built-in functions such at `make`, `append`, `len`
and with slice notation like `[s:e]`. For example:

	ls := make(orb.LineString, 0, 100)
	ls = append(ls, orb.Point{1, 1})
	point := ls[0]

### Shared `Geometry` interface

All of the base types implement the `orb.Geometry` interface defined as:

	type Geometry interface {
		GeoJSONType() string
		Dimensions() int // e.g. 0d, 1d, 2d
		Bound() Bound
	}

This interface is accepted by functions in the sub-packages which then act on the
base types correctly. For example:

	l := clip.Geometry(bound, geom)

will use the appropriate clipping algorithm depending on if the input is 1d or 2d,
e.g. a `orb.LineString` or a `orb.Polygon`.

Only a few methods are defined directly on these type, for example `Clone`, `Equal`, `GeoJSONType`.
Other operation that depend on geo vs. planar contexts are defined in the respective sub-package.
For example:

* Computing the geo distance between two point:

		p1 := orb.Point{-72.796408, -45.407131}
		p2 := orb.Point{-72.688541, -45.384987}

		geo.Distance(p1, p2)

* Compute the planar area and centroid of a polygon:

		poly := orb.Polygon{...}
		centroid, area := planar.CentroidArea(poly)

## GeoJSON

The [geojson](geojson) sub-package implements Marshalling and Unmarshalling of GeoJSON data.
Features are defined as:

	type Feature struct {
		ID         interface{}  `json:"id,omitempty"`
		Type       string       `json:"type"`
		Geometry   orb.Geometry `json:"geometry"`
		Properties Properties   `json:"properties"`
	}

Defining the geometry as an `orb.Geometry` interface along with sub-package functions
accepting geometries allows them to work together to create easy to follow code.
For example, clipping all the geometries in a collection:

	fc, err := geojson.UnmarshalFeatureCollection(data)
	for _, f := range fc {
		f.Geometry = clip.Geometry(bound, f.Geometry)
	}

## Mapbox Vector Tiles

The [encoding/mvt](encoding/mvt) sub-package implements Marshalling and
Unmarshalling [MVT](https://www.mapbox.com/vector-tiles/) data.
This package uses sets of `geojson.FeatureCollection` to define the layers,
keyed by the layer name. For example:

	collections := map[string]*geojson.FeatureCollection{}

	// Convert to a layers object and project to tile coordinates.
	layers := mvt.NewLayers(collections)
	layers.ProjectToTile(maptile.New(x, y, z))

    // In order to be used as source for MapboxGL geometries need to be clipped
    // to max allowed extent. (uncomment next line)
    // layers.Clip(mvt.MapboxGLDefaultExtentBound)

	// Simplify the geometry now that it's in tile coordinate space.
	layers.Simplify(simplify.DouglasPeucker(1.0))

	// Depending on use-case remove empty geometry, those too small to be
	// represented in this tile space.
	// In this case lines shorter than 1, and areas smaller than 2.
	layers.RemoveEmpty(1.0, 2.0)

	// encoding using the Mapbox Vector Tile protobuf encoding.
	data, err := layers.Marshal() // this data is NOT gzipped.

	// Sometimes MVT data is stored and transfered gzip compressed. In that case:
	data, err := layers.MarshalGzipped()

## Decoding WKB from a database query

Geometries are usually returned from databases in WKB format. The [encoding/wkb](encoding/wkb)
sub-package offers helpers to "scan" the data into the base types directly.
For example:

	row := db.QueryRow("SELECT ST_AsBinary(point_column) FROM postgis_table")

	var p orb.Point
	err := row.Scan(wkb.Scanner(&p))

	db.Exec("INSERT INTO table (point_column) VALUES (ST_GeomFromWKB(?))", wkb.Value(p))

Scanning directly from MySQL columns is supported. By default MySQL returns geometry
data as WKB but prefixed with a 4 byte SRID. To support this, if the data is not
valid WKB, the code will strip the first 4 bytes, the SRID, and try again.
This works for most use cases.

## List of sub-package utilities

* [`clip`](clip) - clipping geometry to a bounding box
* [`encoding/mvt`](encoding/mvt) - encoded and decoding from [Mapbox Vector Tiles](https://www.mapbox.com/vector-tiles/)
* [`encoding/wkb`](encoding/wkb) - well-known binary as well as helpers to decode from the database queries
* [`encoding/wkt`](encoding/wkt) - well-known text encoding
* [`geojson`](geojson) - working with geojson and the types in this package
* [`maptile`](maptile) - working with mercator map tiles
* [`project`](project) - project geometries between geo and planar contexts
* [`quadtree`](quadtree) - quadtree implementation using the types in this package
* [`resample`](resample) - resample points in a line string geometry
* [`simplify`](simplify) - linear geometry simplifications like Douglas-Peucker
