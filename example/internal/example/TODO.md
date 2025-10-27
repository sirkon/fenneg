# TODO

Move to length guard + copy semantics for the sake of performance:

| Benchmark                  | Iterations  | Time per op (ns/op) |
| -------------------------- | ----------- | ------------------- |
| **PerfcheckAppend/encode** | 342,545,304 | 3.372               |
| **PerfcheckAppend/decode** | 372,060,190 | 3.220               |
| **PerfcheckCopy/encode**   | 405,621,120 | 2.966               |
| **PerfcheckCopy/decode**   | 612,631,831 | 1.963               |

Could ignore encoding, but the encoding difference is too much faster to pass.