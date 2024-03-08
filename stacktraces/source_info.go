// Copright Kirk Rader 2024

package stacktraces

// Struct returned by stacktraces.FunctionInfo().
type SourceInfo struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}
