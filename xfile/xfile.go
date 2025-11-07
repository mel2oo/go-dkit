package xfile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/h2non/filetype"
	"github.com/mel2oo/go-dkit/hash"
)

var re = regexp.MustCompile(`^(?:file|image):(?P<Name>[^<]*?)(?:\.(?P<Suffix>[^.<]+))?<(?P<Extension>.+)>`)

type File struct {
	Prefix    string        `json:"prefix"`    // 文件类型: file/image
	Name      string        `json:"name"`      // 文件名
	ID        string        `json:"id"`        // 文件id
	Size      float64       `json:"size"`      // 文件大小
	Extension string        `json:"extension"` // 文件扩展名
	Content   *bytes.Reader `json:"content"`   // 文件内容
}

func New(name string, data []byte) *File {
	prefix := "file"
	_, _, err := image.DecodeConfig(bytes.NewBuffer(data))
	if err == nil {
		prefix = "image"
	}

	if len(name) == 0 {
		name = "default"
	}

	hash := hash.MD5String(data)

	if len(filepath.Ext(name)) == 0 {
		fty, err := filetype.Match(data)
		if err != nil {
			name += ".bin"
		} else {
			name += "." + fty.Extension
		}
	}

	return &File{
		Prefix:    prefix,
		ID:        hash,
		Name:      name,
		Size:      math.Round(float64(len(data))/1024*100) / 100,
		Content:   bytes.NewReader(data),
		Extension: strings.TrimLeft(filepath.Ext(name), "."),
	}
}

func LoadFromJson(data []byte) (*File, error) {
	xf := &File{}
	if err := json.Unmarshal(data, xf); err != nil {
		return nil, err
	}

	if len(xf.Prefix) == 0 || len(xf.Name) == 0 || len(xf.ID) == 0 {
		return nil, errors.New("invalid file json_string")
	}

	return xf, nil
}

func LoadFromRawID(str string) (*File, error) {
	matches := re.FindStringSubmatch(str)
	if len(matches) != 4 {
		return nil, errors.New("invalid file raw_id")
	}

	f := &File{
		Prefix:    matches[0],
		Name:      matches[1],
		Extension: matches[2],
		ID:        matches[3],
	}
	if len(f.Extension) > 0 {
		f.Name += "." + f.Extension
	}

	return f, nil
}

func (m *File) String() string {
	return fmt.Sprintf("%s:%s<%s>", m.Prefix, m.Name, m.ID)
}

func (m *File) Dump() ([]byte, error) {
	return json.Marshal(m)
}

func (m *File) SetName(name string) {
	m.Name = name
	m.Extension = strings.TrimLeft(filepath.Ext(m.Name), ".")
}

func (m *File) SetContent(br *bytes.Reader) {
	m.Content = br
}

func (m *File) SetNameFromHeader(header string) {
	if len(header) > 0 {
		re := regexp.MustCompile(`(?i).*filename\*?=(?:[^']+''|\"?)([^\";]*)\"?$`)
		matches := re.FindStringSubmatch(header)
		if len(matches) < 2 {
			m.Name = header
			return
		}
		filename := matches[1]
		decodedStr, err := url.QueryUnescape(matches[1])
		if err == nil {
			filename = decodedStr
		}
		if strings.Contains(filename, `''`) {
			parts := strings.SplitN(filename, `''`, 2)
			if len(parts) == 2 {
				filename = parts[1]
			}
		}
		m.Name = filename
		m.Extension = strings.TrimLeft(filepath.Ext(m.Name), ".")
	}
}

func (m *File) ToMap() map[string]any {
	val := map[string]any{
		"ID":        m.ID,
		"Name":      m.Name,
		"Size":      m.Size,
		"Extension": m.Extension,
	}

	return val
}
