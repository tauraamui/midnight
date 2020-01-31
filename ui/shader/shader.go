package shader

import (
	"bytes"
	"fmt"

	"github.com/markbates/pkger"
)

type Shader struct {
	Uniforms            map[string]interface{}
	StrReplaceCallbacks []func(src string) string

	path  string
	src   string
	dirty bool
}

func New(srcPath string) *Shader {
	fi, err := pkger.Stat(srcPath)
	if err != nil {
		panic(fmt.Errorf("unable to stat shader src file: %w", err))
	}

	if fi.Size() == 0 {
		panic(fmt.Errorf("shader src file looks empty"))
	}

	shader := Shader{
		path:  srcPath,
		dirty: true,
	}
	// *shader.AmbientLightIntensity = INTENSITY_PER_MINUTE
	return &shader
}

func (s *Shader) Code() (string, error) {
	if len(s.src) > 0 {
		return s.src, nil
	}

	file, err := pkger.Open(s.path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var bytes bytes.Buffer
	_, err = bytes.ReadFrom(file)
	if err != nil {
		return "", err
	}

	s.src = bytes.String()
	for _, cb := range s.StrReplaceCallbacks {
		s.src = cb(s.src)
	}

	return s.src, nil
}

func (s *Shader) Dirty() bool {
	wasDirty := s.dirty
	s.dirty = false
	return wasDirty
}
