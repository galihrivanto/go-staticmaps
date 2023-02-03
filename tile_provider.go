// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import "fmt"

// TileProvider .
type TileProvider interface {
	Name() string
	Attribution() string
}

// StaticMapProvider .
type StaticMapProvider interface {
	TileProvider
	GetURL(zoom int, x, y float64, width, height int) string
}

// MapTileProvider .
type MapTileProvider interface {
	TileProvider
	IgnoreNotFound() bool
	TileSize() int
	Shards() []string
	GetURL(shard string, zoom, x, y int) string
}

// tileProvider encapsulates all infos about a map tile provider service (name, url scheme, attribution, etc.)
type tileProvider struct {
	name           string
	attribution    string
	ignoreNotFound bool
	tileSize       int
	urlPattern     string // "%[1]s" => shard, "%[2]d" => zoom, "%[3]d" => x, "%[4]d" => y
	shards         []string
}

func (t *tileProvider) Name() string {
	return t.name
}

func (t *tileProvider) Attribution() string {
	return t.attribution
}

func (t *tileProvider) IgnoreNotFound() bool {
	return t.ignoreNotFound
}

func (t *tileProvider) TileSize() int {
	return t.tileSize
}

func (t *tileProvider) Shards() []string {
	return t.shards
}

func (t *tileProvider) GetURL(shard string, zoom, x, y int) string {
	return fmt.Sprintf(t.urlPattern, shard, zoom, x, y)
}

// NewTileProviderOpenStreetMaps creates a TileProvider struct for OSM's tile service
func NewTileProviderOpenStreetMaps() MapTileProvider {
	t := &tileProvider{
		name:        "osm",
		attribution: "Maps and Data (c) openstreetmap.org and contributors, ODbL",
		tileSize:    256,
		urlPattern:  "http://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b", "c"},
	}

	return t
}

func newTileProviderThunderforest(name string) MapTileProvider {
	t := &tileProvider{
		name:        fmt.Sprintf("thunderforest-%s", name),
		attribution: "Maps (c) Thundeforest; Data (c) OSM and contributors, ODbL",
		tileSize:    256,
		urlPattern:  "https://%[1]s.tile.thunderforest.com/" + name + "/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b", "c"},
	}
	return t
}

// NewTileProviderThunderforestLandscape creates a TileProvider struct for thundeforests's 'landscape' tile service
func NewTileProviderThunderforestLandscape() MapTileProvider {
	return newTileProviderThunderforest("landscape")
}

// NewTileProviderThunderforestOutdoors creates a TileProvider struct for thundeforests's 'outdoors' tile service
func NewTileProviderThunderforestOutdoors() MapTileProvider {
	return newTileProviderThunderforest("outdoors")
}

// NewTileProviderThunderforestTransport creates a TileProvider struct for thundeforests's 'transport' tile service
func NewTileProviderThunderforestTransport() MapTileProvider {
	return newTileProviderThunderforest("transport")
}

// NewTileProviderStamenToner creates a TileProvider struct for stamens' 'toner' tile service
func NewTileProviderStamenToner() MapTileProvider {
	t := &tileProvider{
		name:        "stamen-toner",
		attribution: "Maps (c) Stamen; Data (c) OSM and contributors, ODbL",
		tileSize:    256,
		urlPattern:  "http://%[1]s.tile.stamen.com/toner/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b", "c", "d"},
	}

	return t
}

// NewTileProviderStamenTerrain creates a TileProvider struct for stamens' 'terrain' tile service
func NewTileProviderStamenTerrain() MapTileProvider {
	t := &tileProvider{
		name:        "stamen-terrain",
		attribution: "Maps (c) Stamen; Data (c) OSM and contributors, ODbL",
		tileSize:    256,
		urlPattern:  "http://%[1]s.tile.stamen.com/terrain/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b", "c", "d"},
	}

	return t
}

// NewTileProviderOpenTopoMap creates a TileProvider struct for opentopomap's tile service
func NewTileProviderOpenTopoMap() MapTileProvider {
	t := &tileProvider{
		name:        "opentopomap",
		attribution: "Maps (c) OpenTopoMap [CC-BY-SA]; Data (c) OSM and contributors [ODbL]; Data (c) SRTM",
		tileSize:    256,
		urlPattern:  "http://%[1]s.tile.opentopomap.org/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b", "c"},
	}

	return t
}

// NewTileProviderWikimedia creates a TileProvider struct for Wikimedia's tile service
func NewTileProviderWikimedia() MapTileProvider {
	t := &tileProvider{
		name:        "wikimedia",
		attribution: "Map (c) Wikimedia; Data (c) OSM and contributors, ODbL.",
		tileSize:    256,
		urlPattern:  "https://maps.wikimedia.org/osm-intl/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{},
	}

	return t
}

// NewTileProviderOpenCycleMap creates a TileProvider struct for OpenCycleMap's tile service
func NewTileProviderOpenCycleMap() MapTileProvider {
	t := &tileProvider{
		name:        "cycle",
		attribution: "Maps and Data (c) openstreetmaps.org and contributors, ODbL",
		tileSize:    256,
		urlPattern:  "http://%[1]s.tile.opencyclemap.org/cycle/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b"},
	}

	return t
}

func newTileProviderCarto(name string) MapTileProvider {
	t := &tileProvider{
		name:        fmt.Sprintf("carto-%s", name),
		attribution: "Map (c) Carto [CC BY 3.0] Data (c) OSM and contributors, ODbL.",
		tileSize:    256,
		urlPattern:  "https://cartodb-basemaps-%[1]s.global.ssl.fastly.net/" + name + "_all/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b", "c", "d"},
	}

	return t
}

// NewTileProviderCartoLight creates a TileProvider struct for Carto's tile service (light variant)
func NewTileProviderCartoLight() MapTileProvider {
	return newTileProviderCarto("light")
}

// NewTileProviderCartoDark creates a TileProvider struct for Carto's tile service (dark variant)
func NewTileProviderCartoDark() MapTileProvider {
	return newTileProviderCarto("dark")
}

// NewTileProviderArcgisWorldImagery creates a TileProvider struct for Arcgis' WorldImagery tiles
func NewTileProviderArcgisWorldImagery() MapTileProvider {
	t := &tileProvider{
		name:        "arcgis-worldimagery",
		attribution: "Source: Esri, Maxar, GeoEye, Earthstar Geographics, CNES/Airbus DS, USDA, USGS, AeroGRID, IGN, and the GIS User Community",
		tileSize:    256,
		urlPattern:  "https://server.arcgisonline.com/arcgis/rest/services/World_Imagery/MapServer/tile/%[2]d/%[4]d/%[3]d",
		shards:      []string{},
	}

	return t
}

// GetTileProviders returns a map of all available TileProviders
func GetTileProviders() map[string]MapTileProvider {
	m := make(map[string]MapTileProvider)

	list := []MapTileProvider{
		NewTileProviderThunderforestLandscape(),
		NewTileProviderThunderforestOutdoors(),
		NewTileProviderThunderforestTransport(),
		NewTileProviderStamenToner(),
		NewTileProviderStamenTerrain(),
		NewTileProviderOpenTopoMap(),
		NewTileProviderOpenStreetMaps(),
		NewTileProviderOpenCycleMap(),
		NewTileProviderCartoLight(),
		NewTileProviderCartoDark(),
		NewTileProviderArcgisWorldImagery(),
		NewTileProviderWikimedia(),
	}

	for _, tp := range list {
		m[tp.Name()] = tp
	}

	return m
}
