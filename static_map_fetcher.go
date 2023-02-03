package sm

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

// StaticMapFetcher .
type StaticMapFetcher struct {
	provider  StaticMapProvider
	cache     TileCache
	userAgent string
}

// StaticMap defines a staticMapFile
type StaticMap struct {
	Img                 image.Image
	X, Y                float64
	Zoom, Width, Height int
}

// NewStaticMapFetcher creates a new NewStaticMapFetcher struct
func NewStaticMapFetcher(provider StaticMapProvider, cache TileCache, online bool) *StaticMapFetcher {
	t := new(StaticMapFetcher)
	t.provider = provider
	t.cache = cache
	t.userAgent = "Mozilla/5.0+(compatible; go-staticmaps/0.1; https://github.com/flopp/go-staticmaps)"

	return t
}

// SetUserAgent sets the HTTP user agent string used when downloading map tiles
func (t *StaticMapFetcher) SetUserAgent(a string) {
	t.userAgent = a
}

func cacheStaticMapFileName(cache TileCache, providerName string, zoom int, x, y float64) string {
	return path.Join(
		cache.Path(),
		providerName,
		strconv.Itoa(zoom),
		fmt.Sprintf("%f-%f", x, y),
	)
}

// Fetch download (or retrieves from the cache) a tile image for the specified zoom level and tile coordinates
func (t *StaticMapFetcher) Fetch(m *StaticMap) error {
	if t.cache != nil {
		fileName := cacheStaticMapFileName(t.cache, t.provider.Name(), m.Zoom, m.X, m.Y)
		cachedImg, err := t.loadCache(fileName)
		if err == nil {
			m.Img = cachedImg
			return nil
		}
	}

	url := t.provider.GetURL(m.Zoom, m.X, m.Y, m.Width, m.Height)
	data, err := t.download(url)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if t.cache != nil {
		fileName := cacheStaticMapFileName(t.cache, t.provider.Name(), m.Zoom, m.X, m.Y)
		if err := t.storeCache(fileName, data); err != nil {
			log.Printf("Failed to store map tile as '%s': %s", fileName, err)
		}
	}

	m.Img = img
	return nil
}

func (t *StaticMapFetcher) download(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", t.userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		// Great! Nothing to do.

	case http.StatusNotFound:
		return nil, errTileNotFound

	default:
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (t *StaticMapFetcher) loadCache(fileName string) (image.Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (t *StaticMapFetcher) createCacheDir(path string) error {
	src, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, t.cache.Perm())
		}
		return err
	}
	if src.IsDir() {
		return nil
	}

	return fmt.Errorf("file exists but is not a directory: %s", path)
}

func (t *StaticMapFetcher) storeCache(fileName string, data []byte) error {
	dir, _ := filepath.Split(fileName)

	if err := t.createCacheDir(dir); err != nil {
		return err
	}

	// Create file using the configured directory create permission with the
	// 'x' bit removed.
	file, err := os.OpenFile(
		fileName,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		t.cache.Perm()&0666,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
