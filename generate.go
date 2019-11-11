package wrappy

type GenerateWrappersOptions struct {
	StructWhitelist map[string]struct{}
}

func GenerateWrappers(inputDir string, outputDir string, options GenerateWrappersOptions) {
	structs := ParseWithWhitelist(inputDir, options.StructWhitelist)
	files := Remodel(structs, inputDir, outputDir)
	WriteCode(files)
}
