//go:build integration

package testcase

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Case struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Directory   string `yaml:"directory"`
	Error       string `yaml:"error"`

	root string
}

func Load(t *testing.T, path string) []Case {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var cases []Case
	require.NoError(t, yaml.Unmarshal(data, &cases))
	require.NotEmpty(t, cases)

	for i := range cases {
		cases[i].root = filepath.Join(filepath.Dir(path), "testcases", cases[i].Directory)
	}

	return cases
}

func (c Case) Path(name string) string {
	return filepath.Join(c.root, name)
}

func (c Case) ReadRequest(t *testing.T, request proto.Message) {
	t.Helper()

	readJSON(t, c.Path("request.json"), request)
}

func (c Case) AssertResponse(t *testing.T, response proto.Message, err error) {
	t.Helper()

	if c.Error != "" {
		require.Error(t, err)
		require.Equal(t, c.Error, status.Convert(err).Message())
		require.Nil(t, response)

		return
	}
	require.NoError(t, err)
	require.NotNil(t, response)

	expected := response.ProtoReflect().New().Interface()
	readJSON(t, c.Path("response.json"), expected)
	require.True(t, proto.Equal(expected, response), "expected: %s\nactual: %s",
		protojson.Format(expected), protojson.Format(response))
}

func readJSON(t *testing.T, path string, message proto.Message) {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.NoError(t, protojson.Unmarshal(data, message), "decode %s", path)
}
