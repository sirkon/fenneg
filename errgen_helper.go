package fenneg

import "github.com/sirkon/fenneg/internal/er"

// ReturnError returns a handler to simplify building error
// messages with structured context in a way that closely
// replicates the API of github.com/sirkon/errors error
// processing package.
//
// Example usage:
//    fenneg.ReturnError().Wrap("err", "$decode").Int("count-$0", 14).Rend(r, countIndex)
//
// The code rendered will look like – countIndex, say, is equal to 25:
//    return errors.Wrap(err, "decode BranchName.argName(argType)").Int("count-25", 14)
// There is a set of predefined values in the code rendering context, besides $decode.
// Here is the full list (2023-05-15):
//   - $decode -> "decode $branch.$dst($dstType)"
//	 - $recordTooSmall -> "record buffer is too small"
//	 - $malformedUvarint -> "malformed uvarint sequence"
//	 - $malformedVarint ->  "malformed varint sequence"
func ReturnError() *er.RType {
	return er.Return()
}

type (
	// Q is to be used for structured values to be string.
	// Like
	//     orlgen.ReturnError()...Str("key", fenneg.Q("value"))...
	// To be rendered like
	//     return errors...Str("key", "value")
	// It will be
	//     return errors...Str("key", value)
	// Without fenneg.Q
	Q = er.Q

	// L does the opposite job to Q for structured values:
	// When you want the builder key to be shown as a variable
	// you use fenneg.L:
	//     fenneg.ReturnError()...Str(fenneg.L("key"), fenneg.Q("value"))...
	// This will be rendered as
	//     return errors...Str(key, "value")
	L = er.L
)
