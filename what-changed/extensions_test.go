// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package what_changed

import (
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestCompareExtensions(t *testing.T) {

	left := `x-test: 1`
	right := `x-test: 2`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lExt := low.ExtractExtensions(lNode.Content[0])
	rExt := low.ExtractExtensions(rNode.Content[0])

	extChanges := CompareExtensions(lExt, rExt)

	assert.Len(t, extChanges.Changes, 1)
	assert.Equal(t, Modified, extChanges.Changes[0].ChangeType)
	assert.Equal(t, "1", extChanges.Changes[0].Original)
	assert.Equal(t, "2", extChanges.Changes[0].New)
	assert.False(t, extChanges.Changes[0].Context.HasChanged())
}

func TestCompareExtensions_Moved(t *testing.T) {

	left := `pizza: pie
x-test: 1`

	right := `x-test: 1`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lExt := low.ExtractExtensions(lNode.Content[0])
	rExt := low.ExtractExtensions(rNode.Content[0])

	extChanges := CompareExtensions(lExt, rExt)

	assert.Len(t, extChanges.Changes, 1)
	assert.Equal(t, Moved, extChanges.Changes[0].ChangeType)
	assert.Equal(t, 2, extChanges.Changes[0].Context.OrigLine)
	assert.Equal(t, 1, extChanges.Changes[0].Context.NewLine)
	assert.True(t, extChanges.Changes[0].Context.HasChanged())
}

func TestCompareExtensions_ModifiedAndMoved(t *testing.T) {

	left := `pizza: pie
x-test: 1`

	right := `x-test: 2`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lExt := low.ExtractExtensions(lNode.Content[0])
	rExt := low.ExtractExtensions(rNode.Content[0])

	extChanges := CompareExtensions(lExt, rExt)

	assert.Len(t, extChanges.Changes, 1)
	assert.Equal(t, ModifiedAndMoved, extChanges.Changes[0].ChangeType)
	assert.Equal(t, 2, extChanges.Changes[0].Context.OrigLine)
	assert.Equal(t, 1, extChanges.Changes[0].Context.NewLine)
	assert.Equal(t, "1", extChanges.Changes[0].Original)
	assert.Equal(t, "2", extChanges.Changes[0].New)
	assert.True(t, extChanges.Changes[0].Context.HasChanged())
}

func TestCompareExtensions_Removed(t *testing.T) {

	left := `pizza: pie
x-test: 1`

	right := `pizza: pie`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lExt := low.ExtractExtensions(lNode.Content[0])
	rExt := low.ExtractExtensions(rNode.Content[0])

	extChanges := CompareExtensions(lExt, rExt)

	assert.Len(t, extChanges.Changes, 1)
	assert.Equal(t, PropertyRemoved, extChanges.Changes[0].ChangeType)
	assert.Equal(t, 2, extChanges.Changes[0].Context.OrigLine)
	assert.Equal(t, -1, extChanges.Changes[0].Context.NewLine)
	assert.Equal(t, "1", extChanges.Changes[0].Original)
	assert.True(t, extChanges.Changes[0].Context.HasChanged())
}

func TestCompareExtensions_Added(t *testing.T) {

	left := `pizza: pie`

	right := `pizza: pie
x-test: 1`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lExt := low.ExtractExtensions(lNode.Content[0])
	rExt := low.ExtractExtensions(rNode.Content[0])

	extChanges := CompareExtensions(lExt, rExt)

	assert.Len(t, extChanges.Changes, 1)
	assert.Equal(t, PropertyAdded, extChanges.Changes[0].ChangeType)
	assert.Equal(t, -1, extChanges.Changes[0].Context.OrigLine)
	assert.Equal(t, 2, extChanges.Changes[0].Context.NewLine)
	assert.Equal(t, "1", extChanges.Changes[0].New)
	assert.True(t, extChanges.Changes[0].Context.HasChanged())
}

func TestCompareExtensions_Identical(t *testing.T) {

	left := `x-test: 1`

	right := `x-test: 1`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lExt := low.ExtractExtensions(lNode.Content[0])
	rExt := low.ExtractExtensions(rNode.Content[0])

	extChanges := CompareExtensions(lExt, rExt)

	assert.Nil(t, extChanges)
}
