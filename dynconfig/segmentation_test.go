package dynconfig_test

import (
	"testing"

	"github.com/ExtraWhy/internal-libs/dynconfig"
	"github.com/stretchr/testify/assert"
)

const testConfigData = `{
	"level1_a": {
		"level2_a": {
			"key1": "aa_value1",
			"key2": "aa_value2"
		},
		"*": {
			"key1": "default_a_Value1",
			"key2": "default_a_Value2"
		}
	},
	"level1_b": {
		"level2_b": {
			"key1": "bb_value1",
			"key2": "bb_value2"
		},
		"default": {
			"key1": "default_b_Value1",
			"key2": "default_b_Value2"
		}
	},
	"level1_c": {
		"key1": "c_value1",
		"key2": "c_value2"
	},
	"level1_d": {
		"key1": "d_value1",
		"key2": "d_value2",
		"default": {
			"key1": "default_d_value1",
			"key2": "default_d_value2"
		}
	},
	"default": {
		"key1": "default_Value1",
		"key2": "default_Value2"
	}
}`

func TestSegmentationConfig_DefaultStar(t *testing.T) {
	seg := dynconfig.NewSegmentation(dynconfig.NewEntry(testConfigData))

	value, err := seg.Get("level1_a", "level2_a", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "aa_value1", value)

	value, err = seg.Get("level1_a", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "default_a_Value1", value)

	value, err = seg.Get("level1_a", "someOtherLevel", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "default_a_Value1", value)

	value, err = seg.Get("level1_a", "level2_a", "somethingElse", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "aa_value1", value)
}

func TestSegmentationConfig_Get_DefaultDefault(t *testing.T) {
	seg := dynconfig.NewSegmentation(dynconfig.NewEntry(testConfigData))

	value, err := seg.Get("level1_b", "level2_b", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "bb_value1", value)

	value, err = seg.Get("level1_b", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "default_b_Value1", value)

	value, err = seg.Get("level1_b", "someOtherLevel", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "default_b_Value1", value)

	value, err = seg.Get("level1_b", "level2_b", "somethingElse", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "bb_value1", value)
}

func TestSegmentationConfig_Get_NoImmediateDefault(t *testing.T) {
	seg := dynconfig.NewSegmentation(dynconfig.NewEntry(testConfigData))

	value, err := seg.Get("level1_c", "level2_c", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "c_value1", value)

	value, err = seg.Get("level1_c", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "c_value1", value)

	value, err = seg.Get("level1_c", "someOtherLevel", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "c_value1", value)

	value, err = seg.Get("level1_c", "level2_c", "somethingElse", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "c_value1", value)
}

func TestSegmentationConfig_Get_PrecedenceOfDyrectOverDefaultKey(t *testing.T) {
	seg := dynconfig.NewSegmentation(dynconfig.NewEntry(testConfigData))

	value, err := seg.Get("level1_d", "level2_d", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "d_value1", value)

	value, err = seg.Get("level1_d", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "d_value1", value)

	value, err = seg.Get("level1_d", "someOtherLevel", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "d_value1", value)

	value, err = seg.Get("level1_d", "level2_d", "somethingElse", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "d_value1", value)
}

func TestSegmentationConfig_Get_FallbackToRootDefault(t *testing.T) {
	seg := dynconfig.NewSegmentation(dynconfig.NewEntry(testConfigData))

	value, err := seg.Get("level1_X", "level2_X", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "default_Value1", value)

	value, err = seg.Get("level1_X", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "default_Value1", value)

	value, err = seg.Get("level1_X", "level1_XYZ", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "default_Value1", value)

	value, err = seg.Get("level1_X", "level2_Y", "someOtherLevel", "key1").AsString()
	assert.Nil(t, err)
	assert.Equal(t, "default_Value1", value)
}
