package k

import (
	"bufio"
	"io"
)

func CountLines(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	lines := 0
	// 如果还有下一行，true 继续
	// 如果遇到 error， 暂存 err (setErr("..."))，然后false 退出
	for sc.Scan() {
		lines++
	}

	return lines, sc.Err()
}
