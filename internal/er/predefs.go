package er

import "github.com/sirkon/fenneg/internal/renderer"

var predefs = map[string]string{
	"decode":           "decode ${branch}.${dst}(${dstType})",
	"recordTooSmall":   "record buffer is too small",
	"malformedUvarint": "malformed uvarint sequence",
	"malformedVarint":  "malformed varint sequence",
}

func feedRenderer(r *renderer.Go) *renderer.Go {
	r = r.Scope()

	vals := []string{"branch", "dst", "dstType"}
	for _, val := range vals {
		if !r.InCtx(val) {
			r.Let(val, "")
		}
	}

	for k, v := range predefs {
		if k == "decode" && r.S("$branch") == "" {
			v = "decode ${dst}(${dstType})"
		}
		r.Let(k, r.S(v))
	}

	return r
}
