package repos

import (
	"fmt"
	"strings"
)

func convertToSetStrs(args []string) string {
	strs := make([]string, 0)
	for i, v := range args {
		strs = append(strs, fmt.Sprintf("%s=$%d", v, i+1))
	}
	return strings.Join(strs, ", ")
}
