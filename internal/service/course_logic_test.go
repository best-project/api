package service

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestCalculateLevel(t *testing.T) {
	logic := NewCourseLogic(nil, 0.7)

	lvl, nextLvl := logic.calculateLevel(1000)

	assert.Equal(t, lvl, 7)
	assert.Equal(t, nextLvl, 204)
}
