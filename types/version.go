package types

import (
	"fmt"
	"strconv"
	"strings"
)

// 版本号
type Version struct {
	Major    int `json:"major"`
	Minor    int `json:"minor"`
	Build    int `json:"build"`
	Revision int `json:"revision"`
}

// 解析版本号
// version: major.minor.build.revision
func (s *Version) Parse(version string) bool {
	numbers := strings.Split(version, ".")
	if len(numbers) != 4 {
		return false
	}

	// major
	n, err := strconv.Atoi(numbers[0])
	if err != nil {
		return false
	}
	s.Major = n

	// minor
	s.Minor, err = strconv.Atoi(numbers[1])
	if err != nil {
		return false
	}

	// build
	s.Build, err = strconv.Atoi(numbers[2])
	if err != nil {
		return false
	}

	// revision
	s.Revision, err = strconv.Atoi(numbers[3])
	if err != nil {
		return false
	}

	return true
}

// 转字符串(major.minor.build.revision)
func (s *Version) Full() string {
	return fmt.Sprintf("%d.%d.%d.%d", s.Major, s.Minor, s.Build, s.Revision)
}

// 主版本号转字符串(major.minor)
func (s *Version) Main() string {
	return fmt.Sprintf("%d.%d", s.Major, s.Minor)
}

// 比较大小
// 0：s = version
// 1: s > version
// -1: s < version
func (s *Version) Compare(version *Version) int {
	if s.Major > version.Major {
		return 1
	} else if s.Major < version.Major {
		return -1
	}

	if s.Minor > version.Minor {
		return 1
	} else if s.Minor < version.Minor {
		return -1
	}

	if s.Build > version.Build {
		return 1
	} else if s.Build < version.Build {
		return -1
	}

	if s.Revision > version.Revision {
		return 1
	} else if s.Revision < version.Revision {
		return -1
	}

	return 0
}
