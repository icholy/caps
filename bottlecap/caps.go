package bottlecap

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math/rand"
	"path/filepath"

	"github.com/icholy/caps/imageutil"
)

// Cap stores a cap's image and average color
type Cap struct {
	Name  string
	Image image.Image
	Color color.Color
}

// String returns a string representation of the cap image
func (c *Cap) String() string {
	return fmt.Sprintf("Cap(%s)", c.Name)
}

// BestMatch returns the best matching cap based on its average color
func BestMatch(c color.Color, caps []Cap) Cap {
	var (
		bestindex int
		bestdist  float64
	)
	for i, cap := range caps {
		dist := imageutil.ColorDistance(c, cap.Color)
		if i == 0 || dist < bestdist {
			bestdist = dist
			bestindex = i
		}
	}
	return caps[bestindex]
}

// Read a cap image file
func Read(name string, radius float64) (Cap, error) {
	m, err := imageutil.Read(name)
	if err != nil {
		return Cap{}, err
	}
	size := radius * 2
	return Cap{
		Name:  name,
		Image: imageutil.Resize(m, size, size),
		Color: imageutil.AverageColor(m),
	}, nil
}

// MustRead reads a cap image file and panics if it fails
func MustRead(name string, radius float64) Cap {
	cap, err := Read(name, radius)
	if err != nil {
		panic(err)
	}
	return cap
}

// ReadDir reads all the caps in a directory
func ReadDir(dir string, radius float64) ([]Cap, error) {
	var caps []Cap
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		cap, err := Read(filepath.Join(dir, file.Name()), radius)
		if err != nil {
			return nil, err
		}
		caps = append(caps, cap)
	}
	return caps, nil
}

// Shuffle the bottlecaps
func Shuffle(caps []Cap) {
	rand.Shuffle(len(caps), func(i, j int) {
		caps[i], caps[j] = caps[j], caps[i]
	})
}

// MustReadDir reads all the caps in a directory and panics on failure
func MustReadDir(dir string, radius float64) []Cap {
	caps, err := ReadDir(dir, radius)
	if err != nil {
		panic(err)
	}
	return caps
}

// AssignBestColors assigns colors from the pallet to the colors.
// The order of the caps affects the color assignment.
func AssignBestColors(caps []Cap, pallet []color.Color) {
	if len(caps) != len(pallet) {
		panic("number of caps doesn't match number of colors")
	}
	pallet = imageutil.CloneColors(pallet)
	for i, cap := range caps {
		// find the best match for the cap
		var (
			bestindex int
			bestdist  float64
		)
		for i, c := range pallet {
			dist := imageutil.ColorDistance(c, cap.Color)
			if i == 0 || dist < bestdist {
				bestdist = dist
				bestindex = i
			}
		}
		caps[i].Color = pallet[bestindex]
		// remove the used color
		pallet[bestindex] = pallet[len(pallet)-1]
		pallet = pallet[:len(pallet)-1]
	}
}

// AssignBestColors2 assigns colors from the pallet to the colors.
// The order of the color affects the assignment.
func AssignBestColors2(caps []Cap, pallet []color.Color) {
	if len(caps) != len(pallet) {
		panic("number of caps doesn't match number of colors")
	}
	// used cap indexes
	used := map[int]bool{}
	for _, c := range pallet {
		// find the best cap for the color
		var (
			bestindex int
			bestdist  float64
		)
		for i, cap := range caps {
			if used[i] {
				continue
			}
			dist := imageutil.ColorDistance(c, cap.Color)
			if i == 0 || dist < bestdist {
				bestdist = dist
				bestindex = i
			}
		}
		caps[bestindex].Color = c
		used[bestindex] = true
	}
}
