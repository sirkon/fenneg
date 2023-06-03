package handlers

import "github.com/sirkon/fenneg/internal/renderer"

func dstExtend(r *renderer.Go, howmuch int) {
	r.L(`$dst = $dst[:len($dst)+$0]`, howmuch)
}
