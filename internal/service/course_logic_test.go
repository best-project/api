package service

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestCalculateLevel(t *testing.T) {
	logic := NewCourseLogic(nil, 0.7)

	assert.Equal(t, logic.calculateLevel(1000), 6)
}
