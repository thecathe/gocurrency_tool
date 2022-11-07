package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/thecathe/gocurrency_tool/analyser/log"
	"golang.org/x/tools/go/packages"
)

// parse a particular dir
func ParseDir(proj_name string, path_to_dir string, path_to_main_dir string) PackageCounter {

	var fileSet *token.FileSet = token.NewFileSet()
	var counter PackageCounter = PackageCounter{
		Counter: Counter{
			Go_count:     0,
			Send_count:   0,
			Rcv_count:    0,
			Chan_count:   0,
			IsPackage:    true,
			Project_name: proj_name},
		File_counters: []*Counter{}}

	f, err := parser.ParseDir(fileSet, path_to_dir, nil, parser.AllErrors)

	if proj_name == "test" {
		os.RemoveAll("_test_ast")
		EnsureDir("_test_ast", true)
		// write each file in set to file
		// save_file_ast := func(_f *token.File) bool {
		// 	file, err := os.Create(fmt.Sprintf("_test_ast\\%s.txt", filepath.Base(_f.Name())))
		// 	if err == nil {
		// 		// var f_writer *io.Writer = file.
		// 		var f_writer bytes.Buffer
		// 		var single_fileset *token.FileSet = token.NewFileSet()
		// 		single_fileset.AddFile(_f.Name(), _f.Base(), _f.Size())
		// 		ast.Fprint(&f_writer, single_fileset, f, ast.NotNilFilter)
		// 		file.WriteString(f_writer.String())
		// 	} else {
		// 		FailureLog("Parse, Dir: Test, AST Print: Error...\n\t%v\n", err)
		// 	}
		// 	file.Close()
		// 	GeneralLog("Finished writing to \"_test_ast\\\"%s\n", file.Name())
		// 	return true
		// }
		// fileSet.Iterate(save_file_ast)

		// ast.Print(fileSet, f)

		// write to file
		file, err := os.Create(fmt.Sprintf("_test_ast\\%s.yml", filepath.Base(path_to_dir)))
		if err == nil {
			var f_writer bytes.Buffer
			ast.Fprint(&f_writer, fileSet, f, ast.NotNilFilter)
			file.WriteString(f_writer.String())
			log.DebugLog("Parse, Dir: Test, AST written to file Successfully\n")
		} else {
			log.FailureLog("Parse, Dir: Test, AST Print: Error...\n\t%v\n", err)
		}
	}
	if err != nil {
		log.WarningLog("ParseDir: An error was found in package %s...\n\terror: %v\n", filepath.Base(path_to_dir), err)
	}

	if len(f) == 0 {
		return counter
	}

	for pack_name, pack := range f {

		var package_counter_chan chan Counter = make(chan Counter)
		counter.Counter.Package_name = strings.TrimPrefix(strings.TrimPrefix(path_to_dir, path_to_main_dir)+"/"+pack_name, "/")
		counter.Counter.Package_path = path_to_dir
		// Analyse each file
		log.GeneralLog("Parser, Dir: Spawning Goroutine to analyse AST of each file: %d\n", len(pack.Files)-1)
		for name, file := range pack.Files {
			filename := strings.TrimPrefix(strings.TrimPrefix(path_to_dir, path_to_main_dir)+"/"+filepath.Base(name), "/")
			// results sent accross package_counter_chan

			// for testing
			log.DebugLog("Parse, Dir: Spawning Goroutine: %s\n", filename)
			if filename == "tests/async-communication.go" {
				go AnalyseAst(fileSet, pack_name, filename, file, package_counter_chan, name) // launch a goroutine for each file
			}
		}

		// Receive the results of the analysis of each file
		// for range pack.Files {
		for i := 0; i < 1; i++ {

			var new_counter Counter = <-package_counter_chan

			new_counter.IsPackage = false
			new_counter.Project_name = proj_name
			if len(new_counter.Features) > 0 {
				new_counter.Has_feature = true
			}
			counter.Counter.Go_count += new_counter.Go_count
			counter.Counter.Send_count += new_counter.Send_count
			counter.Counter.Rcv_count += new_counter.Rcv_count
			counter.Counter.Chan_count += new_counter.Chan_count
			counter.Counter.Go_in_for_count += new_counter.Go_in_for_count
			counter.Counter.Range_over_chan_count += new_counter.Range_over_chan_count
			counter.Counter.Go_in_constant_for_count += new_counter.Go_in_constant_for_count
			counter.Counter.Array_of_channels_count += new_counter.Array_of_channels_count
			counter.Counter.Sync_Chan_count += new_counter.Sync_Chan_count
			counter.Counter.Known_chan_depth_count += new_counter.Known_chan_depth_count
			counter.Counter.Unknown_chan_depth_count += new_counter.Unknown_chan_depth_count
			counter.Counter.Make_chan_in_for_count += new_counter.Make_chan_in_for_count
			counter.Counter.Make_chan_in_constant_for_count += new_counter.Make_chan_in_constant_for_count
			counter.Counter.Constant_chan_array_count += new_counter.Constant_chan_array_count
			counter.Counter.Chan_slice_count += new_counter.Chan_slice_count
			counter.Counter.Chan_map_count += new_counter.Chan_map_count
			counter.Counter.Close_chan_count += new_counter.Close_chan_count
			counter.Counter.Select_count += new_counter.Select_count
			counter.Counter.Default_select_count += new_counter.Default_select_count
			counter.Counter.Assign_chan_in_for_count += new_counter.Assign_chan_in_for_count
			counter.Counter.Chan_of_chans_count += new_counter.Chan_of_chans_count
			counter.Counter.Send_chan_count += new_counter.Send_chan_count
			counter.Counter.Receive_chan_count += new_counter.Receive_chan_count
			counter.Counter.Param_chan_count += new_counter.Param_chan_count

			counter.File_counters = append(counter.File_counters, &new_counter)

		}

		log.GeneralLog("Parser, Dir: Retrieved AST analysis from all Goroutines\n")

	}

	return counter
}

func ParseConcurrencyPrimitives(path_to_dir string, counter Counter) Counter {
	package_names := []string{}
	log.DebugLog("Parser, PCP: %s\n", path_to_dir)

	walk_err := filepath.Walk(path_to_dir, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			if file.Name() != "vendor" && file.Name() != "third_party" {
				path, _ = filepath.Abs(path)
				package_names = append(package_names, path)
			} else {
				return filepath.SkipDir
			}
		}
		return nil
	})
	log.GeneralLog("Parser, PCP: Found %d packages.\n", len(package_names))

	if walk_err != nil {
		log.FailureLog("Parser, PCP: Error occured during file walk...\n\terror: %v\n", walk_err)
	}

	var ast_map map[string]*packages.Package = make(map[string]*packages.Package)

	var cfg *packages.Config = &packages.Config{Mode: 991, Fset: &token.FileSet{}, Dir: path_to_dir, Tests: true}
	// var cfg *packages.Config = &packages.Config{Mode: packages., Fset: &token.FileSet{}, Dir: path_to_dir, Tests: true}

	package_names = append([]string{"."}, package_names...)
	loaded_packages, err := packages.Load(cfg, package_names...)

	if err != nil {
		log.FailureLog("Parser, PCP: Could not load: %s\n\terror: %v\n", path_to_dir, err)
		log.GeneralLog("Parser, PCP: Attempting to fix project, to load packages.\n")

		// THIS WILL COLLECT ALL MISSING PACKAGES
		init_cmd := exec.Command("go", "mod", "init")
		init_cmd.Dir = path_to_dir
		var init_out bytes.Buffer
		init_cmd.Stdout = &init_out
		var init_err_out bytes.Buffer
		init_cmd.Stderr = &init_err_out
		init_err := init_cmd.Run()

		if init_err != nil {
			log.VerboseLog("Parser, PCP: Attempted to run \"%v\" and it failed:\n\tpath: %s\n\terror: %v\n\tstdout: %v\n\tstderr:\n%v\n", init_cmd.Args, init_cmd.Dir, init_err, init_cmd.Stdout, init_cmd.Stderr)
			// guess module name
			init_cmd = exec.Command("go", "mod", "init", ProjectURL(filepath.Base(path_to_dir)))
			init_cmd.Dir = path_to_dir
			var init_out bytes.Buffer
			init_cmd.Stdout = &init_out
			var init_err_out bytes.Buffer
			init_cmd.Stderr = &init_err_out
			init_err = init_cmd.Run()

			if init_err != nil {
				log.FailureLog("Parser, PCP: Attempted to run \"%v\" and it failed:\n\tpath: %s\n\terror: %v\nThis could be an issue in the project itself or packages required, enable debug log to see the full error.\n", init_cmd.Args, init_cmd.Dir)
				log.DebugLog("Parser, PCP: Attempted to run \"%v\" and it failed:\n\tpath: %s\n\terror: %v\n\tstdout: %v\n\tstderr:\n%v\n", init_cmd.Args, init_cmd.Dir, init_err, init_cmd.Stdout, init_cmd.Stderr)
			} else {
				log.GeneralLog("Parser, PCP: Successfully ran \"%v\"\n\touput:\n%v\n", init_cmd.Args, init_cmd.Stdout)
			}
		} else {
			log.GeneralLog("Parser, PCP: Successfully ran \"%v\"\n\touput:\n%v\n", init_cmd.Args, init_cmd.Stdout)
		}

		tidy_cmd := exec.Command("go", "mod", "tidy")
		tidy_cmd.Dir = path_to_dir
		var tidy_out bytes.Buffer
		tidy_cmd.Stdout = &tidy_out
		var tidy_err_out bytes.Buffer
		tidy_cmd.Stderr = &tidy_err_out
		tidy_err := tidy_cmd.Run()

		if tidy_err != nil {
			log.FailureLog("Parser, PCP: Attempted to run \"%v\" and it failed:\n\tpath: %s\n\terror: \nThis could be an issue in the packages required, enable debug log to see the full error.\n", tidy_cmd.Args, tidy_cmd.Dir)
			log.DebugLog("Parser, PCP: Attempted to run \"%v\" and it failed:\n\tpath: %s\n\terror: %v\n\tstdout: %v\n\tstderr:\n%v\n", tidy_cmd.Args, tidy_cmd.Dir, tidy_err, tidy_cmd.Stdout, tidy_cmd.Stderr)
		} else {
			log.GeneralLog("Parser, PCP: Successfully ran \"%v\"\n\touput:\n%v\n", tidy_cmd.Args, tidy_cmd.Stdout)
			loaded_packages, err = packages.Load(cfg, package_names...)
			if err != nil {
				log.FailureLog("Parser, PCP: Load packages still failed:\n\tpath: %s\n\terror: %v\n\n", path_to_dir, err)
			} else {
				log.GeneralLog("Parser, PCP: Load packages recovered.\n")
			}
		}
	}

	for _, pack := range loaded_packages {
		ast_map[pack.Name] = pack
	}

	log.DebugLog("Parser, PCP: Analysing %d packages.", len(ast_map))
	for pack_name, node := range ast_map {
		// Analyse each package

		log.VerboseLog("Parser, PCP: %s, %d files.\n", pack_name, len(node.Syntax))

		var n_s_file_decl_count int = 0
		var n_s_file_func_decl_count int = 0

		for _, file := range node.Syntax {
			// each file in the package
			for _, decl := range file.Decls {
				// each declaration in the file
				switch decl := decl.(type) {
				case *ast.FuncDecl:
					// Analyse each function decleration
					if decl.Body != nil {
						counter = AnalyseConcurrencyPrimitives(pack_name, decl, counter, cfg.Fset, ast_map)
						n_s_file_func_decl_count++
					}
				}
				n_s_file_decl_count++
			}
		}
		log.VerboseLog("Parser, PCP: Finished %s, %d files.\n\tTotal Decl: %d\n\tFunction Decl: %d\n", pack_name, len(node.Syntax), n_s_file_decl_count, n_s_file_func_decl_count)

		// fmt.Print("\n\n\n\n")
	}

	log.DebugLog("Parser, PCP: Finished %s, %d packages\n\n\n", path_to_dir, len(ast_map))
	return counter
}
