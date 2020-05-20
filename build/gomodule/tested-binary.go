package gomodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
	"strings"
)

var (
	pctx = blueprint.NewPackageContext("github.com/SergeyStrashko/design-practice-2/build/gomodule")

	goTest = pctx.StaticRule("gotest", blueprint.RuleParams{
		Command:     "cd $workDir && go test -v $testPkg > $outReportPath",
		Description: "test $testPkg",
	}, "workDir", "testPkg", "outReportPath")

	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")
)


type BinaryModule struct {
	blueprint.SimpleName

	properties struct {
		Pkg string
		TestPkg string
		Srcs []string
		SrcsExclude []string
		VendorFirst bool
	}
}


func sliceIncludes(element string, slice []string) bool {
	includes := false
	for _, v := range slice {
		if element == v {
			includes = true
		}
	}
	return includes
}

func (gb *BinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	testReportName := name + ".txt"
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)
	outReportPath := path.Join(config.BaseOutputDir, "reports", testReportName)

	var inputs []string
	var testInputs []string
	inputErrors := false
	for _, src := range gb.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, gb.properties.SrcsExclude); err == nil {
			testInputs = append(testInputs, matches...)

			for _, input := range matches {
				if !strings.Contains(input, "_test.go") && !sliceIncludes(input, inputs) {
					inputs = append(inputs, input)
				}
			}
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErrors = true
		}
	}
	if inputErrors {
		return
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as Go binary", name),
		Rule:        goBuild,
		Outputs:     []string{outputPath},
		Implicits:   inputs,
		Args: map[string]string{
			"outputPath": outputPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        gb.properties.Pkg,
		},
	})

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as Go test report", testReportName),
		Rule:        goTest,
		Outputs:     []string{outReportPath},
		Implicits:   testInputs,
		Args: map[string]string{
			"outReportPath": outReportPath,
			"workDir":       ctx.ModuleDir(),
			"testPkg":       gb.properties.TestPkg,
		},
	})

}

func SimpleBinFactory() (blueprint.Module, []interface{}) {
	mType := &BinaryModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
