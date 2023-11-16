package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/yujuncen/filehisto/histogram"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

var (
	files   = pflag.String("path", ".", "The path for collecting.")
	absTime = pflag.BoolP("abs-time", "t", false, "Show absolute datetime. (instead of timeago.)")
	bucket  = pflag.Int("buckets", 15, "How many buckets should we use?")
)

func readFileStat(p string) ([]fs.FileInfo, error) {
	out := make([]fs.FileInfo, 0)
	err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			stat, err := os.Stat(path)
			if err != nil {
				return err
			}
			out = append(out, stat)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

type Files []fs.FileInfo

// Len returns the number of x, y pairs.
func (fs Files) Len() int {
	return len(fs)
}

// XY returns an x, y pair.
func (fs Files) XY(i int) (x float64, y float64) {
	stat := fs[i]
	x = float64(stat.ModTime().Unix())
	y = float64(stat.Size())
	return
}

func plot2(info []fs.FileInfo) error {
	fs := Files(info)
	hist := histogram.Hist(*bucket, fs)
	now := time.Now()
	total := 0
	for _, bucket := range hist.Buckets {
		total += bucket.Count
	}
	fmt.Printf("statistics for %s, total %d files\n", color.GreenString(*files), hist.Count)
	histogram.Fprintf(os.Stdout, hist, histogram.Linear(*bucket), func(v float64) string {
		if *absTime {
			return color.YellowString(time.Unix(int64(v), 0).Format(time.RFC3339))
		}
		return color.YellowString(fmt.Sprintf("%15s", now.Sub(time.Unix(int64(v), 0)).Truncate(time.Second))) + " ago"
	})
	return nil
}

func main() {
	pflag.Parse()
	info, err := readFileStat(*files)
	if err != nil {
		panic(err)
	}
	plot2(info)
}
