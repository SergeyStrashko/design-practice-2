package archive_bin

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
)

var (
	pctx = blueprint.NewPackageContext("github.com/SergeyStrashko/design-practice-2/build/archive_bin")

	archiveBin = pctx.StaticRule("archive_bin", blueprint.RuleParams{
		Command:     "cd $workDir && zip $outputPath -j $inputFile",
		Description: "make archive from $inputFile",
	}, "workDir", "outputPath", "inputFile")
)

type ArchiveModule struct {
	blueprint.SimpleName
	properties struct {
		Src        string
		SrcsExclude []string
	}
}

func (gb *ArchiveModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	archiveName := name + ".zip"
	config := bood.ExtractConfig(ctx)
	outputPath := path.Join(config.BaseOutputDir, "archives", archiveName)
	var input = path.Join(config.BaseOutputDir, "bin", gb.properties.Src)

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as zip archive", name),
		Rule:        archiveBin,
		Outputs:     []string{outputPath},
		Implicits:   nil,
		Args: map[string]string{
			"workDir":    ctx.ModuleDir(),
			"outputPath": outputPath,
			"inputFile":  input,
		},
	})
}

func SimpleZipFactory() (blueprint.Module, []interface{}) {
	mType := &ArchiveModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
