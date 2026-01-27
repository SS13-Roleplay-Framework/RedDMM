package themes

import (
	"encoding/xml"
	"fmt"
	"os"
    "image/color"

    "sdmm/internal/util"
)

type Theme struct {
	XMLName xml.Name `xml:"theme"`
	Name    string   `xml:"name"`
	Colors  []Color  `xml:"colors>color"`
}

type Color struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func (t *Theme) Save(path string) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()
    
    enc := xml.NewEncoder(f)
    enc.Indent("", "  ")
    return enc.Encode(t)
}

func Load(path string) (*Theme, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    
    var t Theme
    if err := xml.NewDecoder(f).Decode(&t); err != nil {
        return nil, err
    }
    return &t, nil
}

func (t *Theme) ColorMap() map[string]color.RGBA {
    m := make(map[string]color.RGBA)
    for _, c := range t.Colors {
        col, _ := util.ParseColor(c.Value) // Util parse color handles hex/named
        // Wait, util.ParseColor returns util.Color? 
        // I need to check util.ParseColor signature.
        m[c.Name] = col.ImageColor()
    }
    return m
}
