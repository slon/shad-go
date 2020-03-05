package hogwarts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func impl(t *testing.T, prereqs map[string][]string, courseList []string) {
	learned := make(map[string]bool)
	for index, course := range courseList {
		for _, prereq := range prereqs[course] {
			if !learned[prereq] {
				t.Errorf("course %v (index %v) depends on course %v which is not taken yet", course, index, prereq)
			}
		}
		if learned[course] {
			t.Errorf("course %v (index %v) has been taken at least twice", course, index)
		}
		learned[course] = true
	}
	for course := range prereqs {
		if !learned[course] {
			t.Errorf("course %v is missing on the list", course)
		}
	}
}

func TestGetCourseList_computerScience(t *testing.T) {
	var computerScience = map[string][]string{
		"algorithms": {"data structures"},
		"calculus":   {"linear algebra"},
		"compilers": {
			"data structures",
			"formal languages",
			"computer organization",
		},
		"data structures":       {"discrete math"},
		"databases":             {"data structures"},
		"discrete math":         {"intro to programming"},
		"formal languages":      {"discrete math"},
		"networks":              {"operating systems"},
		"operating systems":     {"data structures", "computer organization"},
		"programming languages": {"data structures", "computer organization"},
	}
	impl(t, computerScience, GetCourseList(computerScience))
}

func TestGetCourseList_linearScience(t *testing.T) {
	var linearScience = map[string][]string{
		"1": {"0"},
		"2": {"1"},
		"3": {"2"},
		"4": {"3"},
		"5": {"4"},
		"6": {"5"},
		"7": {"6"},
		"8": {"7"},
		"9": {"8"},
	}
	impl(t, linearScience, GetCourseList(linearScience))
}

func TestGetCourseList_naiveScience(t *testing.T) {
	var naiveScience = map[string][]string{
		"здравый смысл":    {},
		"русский язык":     {"здравый смысл"},
		"литература":       {"здравый смысл"},
		"иностранный язык": {"здравый смысл"},
		"алгебра":          {"здравый смысл"},
		"геометрия":        {"здравый смысл"},
		"информатика":      {"здравый смысл"},
		"история":          {"здравый смысл"},
		"обществознание":   {"здравый смысл"},
		"география":        {"здравый смысл"},
		"биология":         {"здравый смысл"},
		"физика":           {"здравый смысл"},
		"химия":            {"здравый смысл"},
		"музыка":           {"здравый смысл"},
	}
	impl(t, naiveScience, GetCourseList(naiveScience))
}

func TestGetCourseList_weirdScience(t *testing.T) {
	var weirdScience = map[string][]string{
		"купи":   {"продай"},
		"продай": {"купи"},
	}
	require.Panics(t, func() {
		impl(t, weirdScience, GetCourseList(weirdScience))
	})
}

func TestGetCourseList_strangeScience(t *testing.T) {
	var strangeScience = map[string][]string{
		"1": {"0"},
		"2": {"1", "3"},
		"3": {"2"},
	}
	require.Panics(t, func() {
		impl(t, strangeScience, GetCourseList(strangeScience))
	})
}
