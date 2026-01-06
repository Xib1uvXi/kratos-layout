package jsonutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-kratos/kratos-layout/pkg/jsonutil"
)

type testModel struct {
	Text string `json:"text"`
}

func TestParseJSON(t *testing.T) {
	jsonStr := `{"text":"hahaha"}`

	var model testModel
	assert.NoError(t, jsonutil.ParseJSON(jsonStr, &model))
	assert.Equal(t, "hahaha", model.Text)
}

func TestStringifyJSON(t *testing.T) {
	var model = testModel{Text: "123"}

	str, err := jsonutil.StringifyJSON(model)
	assert.NoError(t, err)
	assert.Equal(t, `{"text":"123"}`, str)
}
