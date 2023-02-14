package main

import (
	"fmt"
	"log"
	"os"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

const googleStaticMapURL = "https://maps.googleapis.com/maps/api/staticmap?"

// GMapOption .
type GMapOption func(*GMapOptions)

// GMapOptions .
type GMapOptions struct {
	clientID  string
	signature string
	key       string
	styles    []string
}

// GoogleClientID .
func GMapClientID(clientID, signature string) GMapOption {
	return func(gmo *GMapOptions) {
		gmo.clientID = clientID
		gmo.signature = signature
	}
}

// GMapKey .
func GMapKey(key string) GMapOption {
	return func(gmo *GMapOptions) {
		gmo.key = key
	}
}

// GMapStyles .
func GMapStyles(styles ...string) GMapOption {
	return func(gmo *GMapOptions) {
		gmo.styles = append(gmo.styles, styles...)
	}
}

func GMapTileProvider(options ...GMapOption) sm.TileProvider {
	// default option
	opt := &GMapOptions{
		styles: make([]string, 0),
	}

	for _, option := range options {
		option(opt)
	}

	t := new(gMapProvider)
	t.name = "google-map"
	t.attribution = "Google Map (inc)"
	t.options = opt

	return t
}

// gMapProvider .
type gMapProvider struct {
	name        string
	attribution string
	options     *GMapOptions
}

func (p *gMapProvider) Name() string {
	return p.name
}

func (p *gMapProvider) Attribution() string {
	return p.attribution
}

func (p *gMapProvider) TileSize() int {
	// maximum google static width and height is 640
	// use moderate size and lets library merge tiles
	return 256
}

func (p *gMapProvider) GetURL(zoom int, x, y float64, width, height int) string {
	// construct google static map url
	var url string
	if p.options.key != "" {
		url = googleStaticMapURL + "key=" + p.options.key
	} else {
		url = googleStaticMapURL + "client-id=" + p.options.clientID + "&signature=" + p.options.signature
	}

	if width > 0 && height > 0 {
		url += fmt.Sprintf("&size=%dx%d", width, height)
	}

	if x > 0 && y > 0 {
		url += fmt.Sprintf("&center=%f,%f", y, x)
	}

	if zoom > 0 {
		url += fmt.Sprintf("&zoom=%d", zoom)
	}

	for _, style := range p.options.styles {
		url += "&style=" + style
	}

	log.Println("url", url)

	return url
}

func main() {
	ctx := sm.NewContext()
	ctx.SetCenter(s2.LatLngFromDegrees(1.3011624468555132, 103.85775516239742))
	ctx.SetZoom(19)
	ctx.SetSize(2048, 2048)
	ctx.SetTileProvider(
		GMapTileProvider(
			GMapKey(os.Getenv("GOOGLE_MAP_KEY")),
			GMapStyles(
				"feature:poi|visibility:off",
			),
		),
	)

	// no cache
	ctx.SetCache(nil)

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("google-map.png", img); err != nil {
		panic(err)
	}
}
