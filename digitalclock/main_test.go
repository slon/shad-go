package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/tools/testtool"
)

const importPath = "gitlab.com/slon/shad-go/digitalclock"

var binCache testtool.BinCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var teardown testtool.CloseFunc
		binCache, teardown = testtool.NewBinCache()
		defer teardown()

		return m.Run()
	}())
}

func startServer(t *testing.T) (port string, stop func()) {
	binary, err := binCache.GetBinary(importPath)
	require.NoError(t, err)

	port, err = testtool.GetFreePort()
	require.NoError(t, err, "unable to get free port")

	cmd := exec.Command(binary, "-port", port)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr

	require.NoError(t, cmd.Start())

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	stop = func() {
		_ = cmd.Process.Kill()
		<-done
	}

	if err = testtool.WaitForPort(t, time.Second*30, port); err != nil {
		stop()
	}

	require.NoError(t, err)
	return
}

func readImage(fname string) (image.Image, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s: %w", fname, err)
	}
	defer func() { _ = f.Close() }()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("unable to decode %s: %w", fname, err)
	}

	return img, nil
}

func abs(a, b uint32) int64 {
	if a > b {
		return int64(a - b)
	}
	return int64(b - a)
}

func calcImgDiff(i1, i2 image.Image) float64 {
	w, h := i1.Bounds().Dx(), i1.Bounds().Dy()

	var sum int64
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r1, g1, b1, _ := i1.At(x, y).RGBA()
			r2, g2, b2, _ := i2.At(x, y).RGBA()
			sum += abs(r1, r2)
			sum += abs(g1, g2)
			sum += abs(b1, b2)
		}
	}

	return float64(sum) / (float64(w*h) * 0xffff * 3) * 100.0
}

func queryImage(t *testing.T, c *http.Client, url string) image.Image {
	t.Logf("querying: %s", url)

	resp, err := c.Get(url)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "image/png", resp.Header.Get("Content-Type"))

	img, err := png.Decode(resp.Body)
	require.NoError(t, err)

	return img
}

func TestDigitalClock_valid(t *testing.T) {
	port, stop := startServer(t)
	defer stop()

	files, err := ioutil.ReadDir("./testdata")
	require.NoError(t, err)

	c := &http.Client{Timeout: time.Second * 10}

	for _, f := range files {
		t.Run(f.Name(), func(t *testing.T) {
			parts := strings.SplitN(strings.TrimSuffix(f.Name(), ".png"), "_", 2)
			time := strings.Replace(parts[0], ".", ":", 2)

			v := url.Values{}
			v.Add("time", time)
			v.Add("k", parts[1])

			u := fmt.Sprintf("http://localhost:%s/?%s", port, v.Encode())
			img := queryImage(t, c, u)

			expected, err := readImage(path.Join("testdata", f.Name()))
			require.NoError(t, err)

			w, h := img.Bounds().Dx(), img.Bounds().Dy()
			ew, eh := expected.Bounds().Dx(), expected.Bounds().Dy()
			if w != ew || h != eh {
				t.Errorf("expected size %d x %d got %d x %d", ew, eh, w, h)
			}

			diff := calcImgDiff(img, expected)
			t.Logf("image diff: %.2f %%", diff)
			require.True(t, diff < 1.0,
				fmt.Sprintf("%s images are too different (%.2f %%)", f.Name(), diff))
		})
	}
}

func TestDigitalClock_invalid(t *testing.T) {
	port, stop := startServer(t)
	defer stop()

	c := &http.Client{Timeout: time.Second * 10}

	type TestCase struct {
		Time string `json:"time"`
		K    string `json:"k"`
	}

	testName := func(tc *TestCase) string {
		data, err := json.Marshal(tc)
		require.NoError(t, err)
		return string(data)
	}

	for _, tc := range []*TestCase{
		{K: "0"},
		{K: "31"},
		{K: "2.5"},
		{K: "f"},
		{Time: "15:04"},
		{Time: "15:04:0"},
		{Time: "3:04:05"},
		{Time: "15:4:05"},
		{Time: "24:00:00"},
		{Time: "00:00:60"},
		{Time: "f"},
	} {
		t.Run(testName(tc), func(t *testing.T) {
			v := url.Values{}
			if tc.Time != "" {
				v.Add("time", tc.Time)
			}
			if tc.K != "" {
				v.Add("k", tc.K)
			}

			u := fmt.Sprintf("http://localhost:%s/?%s", port, v.Encode())
			t.Logf("querying: %s", u)

			resp, err := c.Get(u)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestDigitalClock_now(t *testing.T) {
	port, stop := startServer(t)
	defer stop()

	c := &http.Client{Timeout: time.Second * 10}

	for _, k := range []int{1, 10, 29} {
		k := k
		t.Run(fmt.Sprintf("k=%d", k), func(t *testing.T) {
			v := url.Values{"k": []string{strconv.Itoa(k)}}
			u := fmt.Sprintf("http://localhost:%s/?%s", port, v.Encode())
			img := queryImage(t, c, u)

			w, h := img.Bounds().Dx(), img.Bounds().Dy()
			ew, eh := getExpectedWidth(k), getExpectedHeight(k)
			if w != ew || h != eh {
				t.Errorf("expected size %d x %d got %d x %d", ew, eh, w, h)
			}
		})
	}
}

func getSymbolWidth(s string) int {
	return len(strings.SplitN(s, "\n", 2)[0])
}

func getExpectedWidth(k int) int {
	return (getSymbolWidth(Zero)*6 + getSymbolWidth(Colon)*2) * k
}

func getExpectedHeight(k int) int {
	return len(strings.Split(Zero, "\n")) * k
}
