package xfile

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/h2non/filetype"
	"github.com/mel2oo/go-dkit/hash"
)

var re = regexp.MustCompile(`^(?:file|image):(?P<Name>[^<]*?)(?:\.(?P<Suffix>[^.<]+))?<(?P<Extension>.+)>`)

type File struct {
	Prefix    string        // 文件类型: file/image
	Name      string        // 文件名
	ID        string        // 文件id
	Size      int64         // 文件大小
	Extension string        // 文件扩展名
	Reader    *bytes.Reader // 文件内容
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
		Size:      int64(len(data)),
		Reader:    bytes.NewReader(data),
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

func (f *File) MarshalJSON() ([]byte, error) {
	var content string
	if f.Reader != nil {
		data, err := io.ReadAll(f.Reader)
		if err != nil {
			return nil, err
		}

		content = base64.StdEncoding.EncodeToString(data)
	}

	return json.Marshal(map[string]any{
		"prefix":    f.Prefix,
		"name":      f.Name,
		"id":        f.ID,
		"size":      f.Size,
		"extension": f.Extension,
		"content":   content,
	})
}

func (f *File) UnmarshalJSON(data []byte) error {
	aux := struct {
		Prefix    string `json:"prefix"`
		Name      string `json:"name"`
		ID        string `json:"id"`
		Size      int64  `json:"size"`
		Extension string `json:"extension"`
		Content   string `json:"content"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	f.Prefix = aux.Prefix
	f.Name = aux.Name
	f.ID = aux.ID
	f.Size = aux.Size
	f.Extension = aux.Extension

	if aux.Content != "" {
		decoded, err := base64.StdEncoding.DecodeString(aux.Content)
		if err != nil {
			return err
		}
		f.Reader = bytes.NewReader(decoded)
	} else {
		f.Reader = nil
	}

	return nil
}

func (m *File) String() string {
	return fmt.Sprintf("%s:%s<%s>", m.Prefix, m.Name, m.ID)
}

func (m *File) SetName(name string) {
	m.Name = name
	m.Extension = strings.TrimLeft(filepath.Ext(m.Name), ".")
}

func (m *File) SetContent(br *bytes.Reader) {
	m.Reader = br
	m.Size = br.Size()
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
