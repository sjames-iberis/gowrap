package generator

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMethod_Declaration(t *testing.T) {
	m := Method{Name: "method"}
	assert.Equal(t, "method() ()", m.Declaration())
}

func TestMethod_Signature(t *testing.T) {
	m := Method{
		Name:    "method",
		Params:  []Param{{Name: "args", Type: "...string"}},
		Results: []Param{{Name: "err", Type: "error"}},
	}
	assert.Equal(t, "(args ...string) (err error)", m.Signature())
}

func TestMethod_ReturnStruct(t *testing.T) {
	t.Run("with results", func(t *testing.T) {
		m := Method{
			Name:    "method",
			Results: []Param{{Name: "err"}},
		}
		assert.Equal(t, "return s.err", m.ReturnStruct("s"))
	})

	t.Run("no results", func(t *testing.T) {
		m := Method{
			Name: "method",
		}
		assert.Equal(t, "return", m.ReturnStruct("s"))
	})
}

func TestMethod_HasResults(t *testing.T) {
	t.Run("with results", func(t *testing.T) {
		m := Method{
			Name:    "method",
			Results: []Param{{Name: "err"}},
		}
		assert.True(t, m.HasResults())
	})

	t.Run("no results", func(t *testing.T) {
		m := Method{
			Name: "method",
		}
		assert.False(t, m.HasResults())
	})
}

func TestMethod_HasParams(t *testing.T) {
	t.Run("with params", func(t *testing.T) {
		m := Method{
			Name:   "method",
			Params: []Param{{}},
		}
		assert.True(t, m.HasParams())
	})

	t.Run("no params", func(t *testing.T) {
		m := Method{
			Name: "method",
		}
		assert.False(t, m.HasParams())
	})
}

func TestMethod_ResultsStruct(t *testing.T) {
	m := Method{
		Name:    "method",
		Results: []Param{{Name: "s", Type: "string"}},
	}
	assert.Equal(t, "struct{\ns string}", m.ResultsStruct())
}

func TestMethod_ResultsNames(t *testing.T) {
	m := Method{
		Name:    "method",
		Results: []Param{{Name: "s"}, {Name: "t"}},
	}
	assert.Equal(t, "s, t", m.ResultsNames())
}

func TestMethod_Pass(t *testing.T) {
	t.Run("no results", func(t *testing.T) {
		m := Method{
			Name:   "method",
			Params: []Param{{Name: "s"}, {Name: "t"}},
		}
		assert.Equal(t, "d.method(s, t)\nreturn", m.Pass("d."))
	})

	t.Run("with results", func(t *testing.T) {
		m := Method{
			Name:    "method",
			Params:  []Param{{Name: "s"}, {Name: "t"}},
			Results: []Param{{Name: "err"}, {Name: "error"}},
		}
		assert.Equal(t, "return d.method(s, t)", m.Pass("d."))
	})
}

func TestMethod_Call(t *testing.T) {
	m := Method{
		Name:   "method",
		Params: []Param{{Name: "s"}, {Name: "t"}},
	}
	assert.Equal(t, "method(s, t)", m.Call())
}

func TestMethod_ParamsMap(t *testing.T) {
	m := Method{
		Name:   "method",
		Params: []Param{{Name: "s", Type: "string"}},
	}
	assert.Equal(t, "map[string]interface{}{\n\"s\": s}", m.ParamsMap())
}

func TestMethod_ResultsMap(t *testing.T) {
	m := Method{
		Name:    "method",
		Results: []Param{{Name: "s", Type: "string"}},
	}
	assert.Equal(t, "map[string]interface{}{\n\"s\": s}", m.ResultsMap())
}

func TestMergeMaps(t *testing.T) {

	tests := []struct {
		name string
		target string
		source string
		want string
	}{
		{
			name:   "AddTopLevelKey",
			target: `{"Log":{"ctx":"XYZ"}}`,
			source: `{"Trace":{"id":"target.NetworkElementID"}}`,
			want:   `{"Log":{"ctx":"XYZ"}, "Trace":{"id":"target.NetworkElementID"}}`,
		},
		{
			name:   "AddSecondLevelKey",
			target: `{"Log":{"ctx":"XYZ"}}`,
			source: `{"Log":{"new":"123"}}`,
			want:   `{"Log":{"ctx":"XYZ", "new": "123"}}`,
		},
		{
			name:   "ReplaceValue",
			target: `{"Log":{"ctx":"XYZ", "new":"123"}}`,
			source: `{"Log":{"ctx":"ABC"}}`,
			want:   `{"Log":{"ctx":"ABC", "new": "123"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := mapFromJson(t, tt.target)
			source := mapFromJson(t, tt.source)
			want := mapFromJson(t, tt.want)

			mergeMaps(target, source)

			assert.Equal(t, want, target)
		})
	}
}

func mapFromJson(t *testing.T, source string) map[string] interface{} {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(source), &m)
	assert.NoError(t, err, "Unmarshal failed")
	return m
}