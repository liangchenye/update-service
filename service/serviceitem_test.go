package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestItemisValid(t *testing.T) {
	cases := []struct {
		i        UpdateServiceItem
		expected bool
	}{
		{i: UpdateServiceItem{FullName: "fullname", SHAS: []string{"sha0"}}, expected: true},
		{i: UpdateServiceItem{FullName: "fullname", SHAS: []string{}}, expected: false},
		{i: UpdateServiceItem{FullName: "", SHAS: []string{"sha0"}}, expected: false},
	}

	for _, c := range cases {
		ok, _ := c.i.isValid()
		assert.Equal(t, c.expected, ok, "Fail to check `isValid`")
	}
}

func TestItemEqual(t *testing.T) {
	cases := []struct {
		a        UpdateServiceItem
		b        UpdateServiceItem
		expected bool
	}{
		{a: UpdateServiceItem{FullName: "fullname"}, b: UpdateServiceItem{FullName: "fullname"}, expected: true},
		{a: UpdateServiceItem{FullName: "fullname"}, b: UpdateServiceItem{FullName: "fullname1"}, expected: false},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.a.Equal(c.b), "Fail to check `Equal`")
	}
}

func TestItemGetSHAS(t *testing.T) {
	item := UpdateServiceItem{FullName: "fullname", SHAS: []string{"sha0"}}
	shas := item.GetSHAS()
	assert.Equal(t, 1, len(shas), "Fail to get correct shas len")
	assert.Equal(t, "sha0", shas[0], "Fail to get correct shas value")
}

func TestItemSetGet(t *testing.T) {
	testItem, _ := NewUpdateServiceItem("fn", []string{"sha0"})

	// set/get expired
	testNewExpired := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	testItem.SetExpired(testNewExpired)
	assert.Equal(t, testNewExpired, testItem.GetExpired(), "Fail to set/get expired time")
	assert.Equal(t, true, testItem.IsExpired(), "Fail to get the expired status")
	testItem.SetExpired(time.Now().Add(time.Hour * 1))
	assert.Equal(t, false, testItem.IsExpired(), "Fail to get the expired status")

	// set/get created
	testNewCreated := testItem.GetCreated().Add(time.Hour * 2)
	testItem.SetCreated(testNewCreated)
	assert.Equal(t, testNewCreated, testItem.GetCreated(), "Fail to set/get created time")
}
