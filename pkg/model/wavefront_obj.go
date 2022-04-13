package model

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var spaceDelimiter = regexp.MustCompile(`\s+`)

type Vertex struct {
	X, Y, Z, W float32
}

func (v Vertex) ToVec3f() Vec3f {
	return Vec3f{float64(v.X), float64(v.Y), float64(v.Z)}
}

type VertexInfo [3]int

type Face []VertexInfo

type Model struct {
	vertices        []Vertex
	verticesNormal  []Vertex
	verticesTexture []Vertex
	faces           []Face
}

func ParseObjVertex(line string) (Vertex, error) {
	var data [4]float32

	// FIXME: Split with regex, to catch multi spaces delimiter
	s := spaceDelimiter.Split(line, -1)[1:] // ignore prefix
	for i := range s {
		f, err := strconv.ParseFloat(s[i], 32)
		if err != nil {
			return Vertex{}, errors.Wrapf(err, "invalid vertex value")
		}
		data[i] = float32(f)
	}

	return Vertex{data[0], data[1], data[2], 1}, nil
}

func ParseObjFaces(line string) (Face, error) {
	var face Face
	s := spaceDelimiter.Split(line, -1)[1:] //ignore prefix
	for i := 0; i < len(s); i++ {
		var vertexInfo VertexInfo
		split := strings.Split(s[i], "/")
		for j, v := range split {
			if len(v) == 0 {
				continue
			}

			val, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.Wrap(err, "invalid faces index value")
			}

			if val < 1 {
				return nil, errors.Wrap(err, "invalid value, index start at 1")
			}

			vertexInfo[j] = val
		}
		face = append(face, vertexInfo)
	}
	return face, nil
}

func NewModel(name string) (*Model, error) {
	f, err := os.OpenFile(name, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	model := &Model{}
	scanner := bufio.NewScanner(f)
	for line := 1; scanner.Scan(); line++ {
		text := strings.TrimSpace(scanner.Text())
		if len(text) == 0 {
			continue
		}

		switch prefix := strings.Split(text, " ")[0]; prefix {
		case "v", "vt", "vn":
			container := &model.vertices
			if prefix == "vt" {
				container = &model.verticesTexture
			} else if prefix == "vn" {
				container = &model.verticesNormal

			}

			val, err := ParseObjVertex(text)
			if err != nil {
				fmt.Printf("prefix: %v", prefix)
				return nil, errors.Wrapf(err, "failed parsing vertex at line %d", line)
			}

			(*container) = append(*container, val)
		case "f":
			val, err := ParseObjFaces(text)
			if err != nil {
				return nil, errors.Wrapf(err, "failed parsing faces at line %d", line)
			}
			model.faces = append(model.faces, val)
		default:
			continue
		}
	}

	return model, nil
}

func (m Model) VertexCount() int {
	return len(m.vertices)
}

func (m Model) TextureCount() int {
	return len(m.verticesTexture)
}

func (m Model) NormalCount() int {
	return len(m.verticesNormal)
}

func (m Model) FacesCount() int {
	return len(m.faces)
}

func (m Model) GetVertex(index int) (Vertex, error) {
	if index < 1 {
		return Vertex{}, errors.New("vertices index start at 1")
	}

	if index > len(m.vertices) {
		return Vertex{}, errors.New("index out of bound")
	}
	return m.vertices[index-1], nil

}

func (m Model) GetFace(index int) Face {
	return m.faces[index]
}
