package codegen

const (
	TargetDarwinArm64  = "arm64-apple-macosx15.0.0"
	TargetDarwinAmd64  = "x86_64-apple-darwin"
	TargetLinuxAmd64   = "x86_64-pc-linux-gnu"
	TargetLinuxArm64   = "aarch64-linux-gnu"
	TargetWindowsAmd64 = "x86_64-w64-mingw32"
)

var targetTripleMap = map[string]map[string]string{
	"darwin": {
		"arm64": TargetDarwinArm64,
		"amd64": TargetDarwinAmd64,
	},
	"linux": {
		"amd64": TargetLinuxAmd64,
		"arm64": TargetLinuxArm64,
	},
	"windows": {
		"amd64": TargetWindowsAmd64,
	},
}

// getTargetTriple returns the LLVM target triple based on environment variables OS and ARCH.
func getTargetTriple(os, arch string) string {
	if os != "" && arch != "" {
		if archMap, ok := targetTripleMap[os]; ok {
			if triple, ok := archMap[arch]; ok {
				return triple
			}
		}
	}
	return "" // let LLVM frontend decide
}
