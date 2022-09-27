package main

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Counter struct {
	Go_count                          int     // Count how many time the term "go" appear in source code
	Send_count                        int     // Count how many time a send  "chan <- val" appear in the source code
	Rcv_count                         int     // Count how many time a rcv "val <- chan" appear in the source code
	Chan_count                        int     // the number of channel overall
	Go_in_for_count                   int     // Count how many times.a goroutine is started in a for loop
	Range_over_chan_count             int     // Count the number of range over a chan
	Go_in_constant_for_count          int     // Goroutine launched in a for loop where the looping is controled by a constant
	Array_of_channels_count           int     // How many unknown length arrays are made chan of
	Sync_Chan_count                   int     // Count how many chan are created in the source code "make(chan type)"
	Known_chan_depth_count            int     // How many make(chan int, n) where n is either a constant or a hard coded number
	Unknown_chan_depth_count          int     // How many make(chan int, n) where n is completely dynamic
	Make_chan_in_for_count            int     // How many time a channel is created in a for loop
	Make_chan_in_constant_for_count   int     // How many time a channel is created in a constant for loop
	Constant_chan_array_count         int     // How many array of channels of constant size
	Chan_slice_count                  int     // How many dynamic array of channels
	Chan_map_count                    int     // how many map of channels
	Close_chan_count                  int     // How many close(chan)
	Select_count                      int     // how many select
	Default_select_count              int     // how many select with a default
	Assign_chan_in_for_count          int     // How many chan are assigned another chan in a for loop
	Assign_chan_in_constant_for_count int     // How many chan are assigned another chan in a for loop
	Chan_of_chans_count               int     // How many channel of channels
	Receive_chan_count                int     // how many receive chan
	Send_chan_count                   int     // how many send only chan
	Param_chan_count                  int     // How many times a chan is used as a param without specifying receives only or write only
	Waitgroup_count                   int     // How many waitgroup declaration are contained
	Known_add_count                   int     // How many known bound of add(n) where n is a constant
	Unknown_add_count                 int     // How many unknown bound of add(n) where n is not a constant
	Done_count                        int     // How many wg.Done()
	Mutex_count                       int     // How many mutex declaration where found
	Unlock_count                      int     // How many unlock in the code
	Lock_count                        int     // How many lock in the code
	IsPackage                         bool    // Return if the counter represent the counter for just a file or the whole package
	Package_name                      string  // The name of the package
	Package_path                      string  // path of the package
	Project_name                      string  // The name of the whole project
	Line_number                       int     // The number of lines in the counter
	Num_of_packages_with_features     int     // The number of package that contains at least one feature
	Has_feature                       bool    // Is there any features in this package ?
	Undefined_over_defined_chans      float64 // percent of undefined chan over defined (chan / chan<-, <-chan)
	Known_over_unknown_chan           float64 // percent of known chan size over unknown
	Features                          []*Feature
	filename                          string // the name of the file
	// RQ3
	Timeout_count        int // How many timeouts in total
	Timeout_select_count int // How many timeouts in selects
	// RQ4
	Select_Sync_R_Excl_Timeout_count                  int // How many selects only have sync. receives, excluding timeout
	Select_Sync_R_Incl_Timeout_count                  int // How many selects only have sync. receives, including timeout
	Select_Sync_R_Default_Excl_Timeout_count          int // How many selects only have sync. receives and defaults, excluding timeout
	Select_Sync_R_Default_Incl_Timeout_count          int // How many selects only have sync. receives and defaults, including timeout
	Select_Async_R_Excl_Timeout_count                 int // How many selects only have async. receives, excluding timeout
	Select_Async_R_Incl_Timeout_count                 int // How many selects only have async. receives, including timeout
	Select_Async_R_Default_Excl_Timeout_count         int // How many selects only have async. receives and defaults, excluding timeout
	Select_Async_R_Default_Incl_Timeout_count         int // How many selects only have async. receives and defaults, including timeout
	Select_Sync_S_count                               int // How many selects only have sync. receives
	Select_Sync_S_Default_count                       int // How many selects only have sync. receives
	Select_Sync_S_Timeout_count                       int // How many selects only have sync. receives
	Select_Sync_S_Default_Timeout_count               int // How many selects only have sync. receives
	Select_Async_S_count                              int // How many selects only have async. receives
	Select_Async_S_Default_count                      int // How many selects only have async. receives
	Select_Async_S_Timeout_count                      int // How many selects only have async. receives
	Select_Async_S_Default_Timeout_count              int // How many selects only have sync. receives
	Select_Sync_S_Sync_R_Excl_Timeout_count           int // How many selects only have sync. send and sync. receive, excluding timeout
	Select_Sync_S_Sync_R_Incl_Timeout_count           int // How many selects only have sync. send and sync. receive, including timeout
	Select_Sync_S_Sync_R_Default_Excl_Timeout_count   int // How many selects only have sync. send and sync. receive and defaults, excluding timeout
	Select_Sync_S_Sync_R_Default_Incl_Timeout_count   int // How many selects only have sync. send and sync. receive and defaults, including timeout
	Select_Async_S_Sync_R_Excl_Timeout_count          int // How many selects only have async. send and sync. receive, excluding timeout
	Select_Async_S_Sync_R_Incl_Timeout_count          int // How many selects only have async. send and sync. receive, including timeout
	Select_Async_S_Sync_R_Default_Excl_Timeout_count  int // How many selects only have async. send and sync. receive and defaults, excluding timeout
	Select_Async_S_Sync_R_Default_Incl_Timeout_count  int // How many selects only have async. send and sync. receive and defaults, including timeout
	Select_Sync_S_Async_R_Excl_Timeout_count          int // How many selects only have sync. send and async. receive, excluding timeout
	Select_Sync_S_Async_R_Incl_Timeout_count          int // How many selects only have sync. send and async. receive, including timeout
	Select_Sync_S_Async_R_Default_Excl_Timeout_count  int // How many selects only have sync. send and async. receive and defaults, excluding timeout
	Select_Sync_S_Async_R_Default_Incl_Timeout_count  int // How many selects only have sync. send and async. receive and defaults, including timeout
	Select_Async_S_Async_R_Excl_Timeout_count         int // How many selects only have async. send and async. receive, excluding timeout
	Select_Async_S_Async_R_Incl_Timeout_count         int // How many selects only have async. send and async. receive, including timeout
	Select_Async_S_Async_R_Default_Excl_Timeout_count int // How many selects only have async. send and async. receive and defaults, excluding timeout
	Select_Async_S_Async_R_Default_Incl_Timeout_count int // How many selects only have async. send and async. receive and defaults, including timeout

}

type PackageCounter struct {
	Counter           Counter    // The overall counter of the package an
	File_counters     []*Counter // the counters of each of the file in the package
	Featured_packages int
	Featured_files    int
	Num_files         int
}

type ProjectCounter struct {
	Counter          Counter // the overall counter for the project
	Package_counters []*PackageCounter
	Project_name     string
	Num_packages     int
}

type GlobalCounter struct {
	Counter          Counter // overall counter for all projects processed
	Project_counters []*ProjectCounter
	Num_projects     int
}

const (
	clone_dir  = "_cloned_projects"
	dir_mode   = os.ModeDir // os.ModePerm
	result_dir = "_results"
	temp_go    = "_temp.go"
)

var general_log_enabled bool = true
var debug_log_enabled bool = false
var warning_log_enabled bool = true
var failure_log_enabled bool = true

var keep_repos bool = false
var overwrite_results bool = false
var overwrite_repos bool = false
var separate_results bool = false
var single_run bool = false

var projects_to_process int = 0

var repo_memory_used int64 = 0
var repo_memory_capacity int64 = 0
var repo_memory_limited bool = false

// uses projects.txt by default
var projects_path string = ".\\projects.txt"

func main() {
	// goconcurrency.exe (projects_file_path, overwrite_results, separate_results, overwrite_repos, debug_log_enabled, warning_log_enabled, failure_log_enabled, dev, repo_capacity)

	// run with git bash in /analyser:
	// go build && ./gocurrency_tool.exe

	// run, keeping clones, deleting old results and splitting results
	// go build && ./gocurrency_tool.exe projects.txt true true true

	// set the type of printouts: defaults
	SetLoggers(general_log_enabled, debug_log_enabled, warning_log_enabled, failure_log_enabled)
	GeneralLog("The program takes the following optional parameters, but also has defaults:\n\n\tprojects_file_path, str: path to input dir\n\n\toverwrite_results, bool: true false.\nPermission to delete anything inside the \"%s\" directory to output results.\n\n\tseparate_results, bool: true false.\nSeparates the CSV output into 3 files.\n\n\toverwrite_repos, bool:true false.\nPermission to delete anything inside the \"%s\" directory to clone repos.\n\n\tdebug_log_enabled, bool: true false\n\n\twarning_log_enabled, bool: true false\n\n\tfailure_log_enabled, bool: true false\n\n\tdev, str:\nspecial keyword \"dev\" will make the program run once, without deleting the repo.\n\n\tmax_repo_memory, int:\nThe program will run more than once, getting as far as it can while staying under the capacity in MB provided here. 0 Means uncapped. The total memory used is checked after each repo has been downloaded. Keeping repos like this will break most of the analysis features. Purely for manual debug/analysis only.\n\n\n\n", result_dir, clone_dir)
	GeneralLog("Example git bash command:\n\ngo build && ./gocurrency_tool.exe projects.txt true true true false true true dev 1000\n\nThis would run the program for all user/repo in \"projects.txt\", overwriting everything, splitting the CSV into separate files, hide debug logs.\n\n\n")
	GeneralLog("Example git bash command for debugging/developing this program:\n\ngo build && ./gocurrency_tool.exe projects.txt true true true false true true dev 1000\n\nThis would run the program for the first project, overwriting everything, save repos with a capacity of ~1 GB splitting the CSV into separate files, and hides debug logs.\n\n\n")

	// goconcurrency.exe (projects_file_path, overwrite_results, separate_results, debug_log_enabled, warning_log_enabled, failure_log_enabled, dev, repo_capacity)
	if len(os.Args) > 1 {
		projects_path = os.Args[1]
		// check for test!
		// for running tests, then exit
		if projects_path == "test" {
			InitDirs()
			GeneralLog("Entering Test:\n\n")
			// force debug
			SetLoggers(true, true, true, true)
			DebugLog("Enabled all Loggers.\n\n")
			var new_counter PackageCounter = ParseDir("test", "tests", "")
			var test_counter Counter = HtmlOutputCounters([]*PackageCounter{&new_counter}, "test", "test", nil, "")

			GeneralLog("Test Counter, Features: %d, before PCP.\n", len(test_counter.Features)-1)
			test_counter = ParseConcurrencyPrimitives("tests", test_counter) // analyses occurences of Waitgroup,mutexes and operations on them

			GeneralLog("Test Counter, Features: %d, after PCP\n", len(test_counter.Features)-1)

			// for i, f := range test_counter.Features {
			// 	GeneralLog("Test: Feature %d:\n%v\n\tline: %v\n\ttype: %v, %v\n\taddit.: %v\n\n", i, f.F_filename, f.F_line_num, f.F_type, f.F_type_num, f.F_number)
			// }

			CsvOutputCounters("tests", []*PackageCounter{&new_counter}, "", test_counter)
			separate_results = true
			CsvOutputCounters("tests", []*PackageCounter{&new_counter}, "", test_counter)

			ExitLog(1, "Finished Tests\n")
		}
		if len(os.Args) > 2 {
			// project to process
			if _process_int, _process_err := strconv.ParseInt(os.Args[2], 10, 64); _process_err == nil {
				projects_to_process = int(_process_int)
			} else {
				WarningLog("Error with 2nd parameter: \"%+v\"\n", os.Args[2])
			}
			if len(os.Args) > 3 {
				// separate results
				if _result_over_bool, _result_over_err := strconv.ParseBool(os.Args[3]); _result_over_err == nil {
					overwrite_results = _result_over_bool
				} else {
					WarningLog("Error with 3rd parameter: \"%+v\"\n", os.Args[3])
				}
				if len(os.Args) > 4 {
					// separate results
					if _result_bool, _result_err := strconv.ParseBool(os.Args[4]); _result_err == nil {
						separate_results = _result_bool
					} else {
						WarningLog("Error with 4th parameter:\"%+v\"\n", os.Args[4])
					}
					if len(os.Args) > 5 {
						// overwrite repos
						if _repo_over_bool, _repo_over_err := strconv.ParseBool(os.Args[5]); _repo_over_err == nil {
							overwrite_repos = _repo_over_bool
						} else {
							WarningLog("Error with 5th parameter: \"%+v\"\n", os.Args[5])
						}
						// check for log declaration, must for all 3
						if len(os.Args) > 6 {
							if len(os.Args) >= 8 {
								// debugging
								if _debug_bool, _debug_err := strconv.ParseBool(os.Args[6]); _debug_err == nil {
									debug_log_enabled = _debug_bool
								} else {
									WarningLog("Error with 6th parameter: \"%+v\"\n", os.Args[6])
								}
								// warning
								if _warning_bool, _warning_err := strconv.ParseBool(os.Args[7]); _warning_err == nil {
									warning_log_enabled = _warning_bool
								} else {
									WarningLog("Error with 7th parameter: \"%+v\"\n", os.Args[7])
								}
								// failure
								if _failure_bool, _failure_err := strconv.ParseBool(os.Args[8]); _failure_err == nil {
									failure_log_enabled = _failure_bool
								} else {
									WarningLog("Error with 8th parameter: \"%+v\"\n", os.Args[8])
								}
								// update
								SetLoggers(general_log_enabled, debug_log_enabled, warning_log_enabled, failure_log_enabled)

								if len(os.Args) > 10 && os.Args[9] == "dev" {
									single_run = true
									keep_repos = true
									// repo capacity
									if _repo_cap_int, _repo_cap_err := strconv.ParseInt(os.Args[10], 10, 64); _repo_cap_err == nil && keep_repos {
										repo_memory_capacity = _repo_cap_int
										repo_memory_limited = true
										single_run = false
										keep_repos = true
									} else {
										WarningLog("Error with 10th parameter: \"%+v\"\n", os.Args[10])
									}
								}
							} else {
								// only partial logging parameters provided
								WarningLog("Not allowed to define logs partially.\n\n")
							}
						}
					}
				}
			}
		} else {
			GeneralLog("Default \"Logging\" options will be used.\n\n")
		}
	} else {
		GeneralLog("No \"projects.txt\" path provided. Will assume one is in the local directory.\n\n")
	}

	if keep_repos && repo_memory_capacity == 0 {
		repo_memory_limited = false
	}

	// print settings
	GeneralLog("Proceeding with the following parameters:\n\tProjects path: %s\n\tProjecs to process: %d, if 0 then all.\n\tOverwriting results: %t\n\tSeparating results: %t\n\tOverwriting repos: %t\n\tSingle run: %t (dev mode)\n\tKeeping repos: %t\n\tRepo Memory Limited: %t, capacity: %d MB\n\n", projects_path, projects_to_process, overwrite_results, separate_results, overwrite_repos, single_run, keep_repos, repo_memory_limited, repo_memory_capacity)

	GeneralLog("Output logging settings:\n\tGeneral: %t\n\tDebug: %t\n\tWarning: %t\n\tFailure: %t\n\n", general_log_enabled, debug_log_enabled, warning_log_enabled, failure_log_enabled)

	GeneralLog("Uneditable logging settings:\n\tVerbose: %t\n\tPanic: %t\n\tExit: %t\n\n", false, true, true)

	// try to read, if cant then exit
	data, e := os.ReadFile(projects_path)
	if e != nil {
		ExitLog(1, "Unable to open the file where the projects are stored...\n\tpath: %s\n\terror: %v\n", projects_path, e)
	}
	// array of user/repo
	proj_listings := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")
	var aborted_projects string

	total_project_count := len(proj_listings) - 1
	if projects_to_process == 0 {
		projects_to_process = total_project_count
	}
	// var project_counters []Counter

	// wipe results?
	if overwrite_results {
		o_r_err_ := os.RemoveAll(result_dir)
		if o_r_err_ == nil {
			GeneralLog("Wiped result dir\n")
		} else if !os.IsNotExist(o_r_err_) {
			PanicLog(o_r_err_, "Unable to wipe result dir\n")
		}
	}
	if overwrite_repos {
		o_r_err_ := os.RemoveAll(clone_dir)
		if o_r_err_ == nil {
			GeneralLog("Wiped clone dir\n")
		} else if !os.IsNotExist(o_r_err_) {
			PanicLog(o_r_err_, "Unable to wipe clone dir\n")
		}
	}

	// assertion
	if html_result_dir != filepath.Join(result_dir, "html") {
		WarningLog("Potential issue: HTML results output dir is not set to expected value...\n\texpected: results\\html\n\tactual: %s\n", html_result_dir)
	}
	if csv_result_dir != filepath.Join(result_dir, "csv") {
		WarningLog("Potential issue: CSV results output dir is not set to expected value...\n\texpected: results\\csv\n\tactual: %s\n", csv_result_dir)
	}

	InitDirs()

	//
	var index_data *IndexFileData = &IndexFileData{Indexes: []*IndexData{}}

	GeneralLog("Initialisation Complete\n\n\n")

	// go through each project:
	// 	clone repo
	//
	GeneralLog("Starting %d/%d projects.\n\n", projects_to_process, total_project_count)
	for _index, project_name := range proj_listings {
		if project_name != "" && _index < projects_to_process {
			GeneralLog("Project %d/%d: %s\n", _index+1, projects_to_process, project_name)

			proj_name := filepath.Base(string(project_name))
			var path_to_dir string
			var commit_hash string

			path_to_dir, commit_hash = CloneRepo(string(project_name))

			if _, _repo_dir_err := os.Stat(path_to_dir); _repo_dir_err != nil {
				// skip this repo as it failed
				FailureLog("Aborting project \"%s\". Error occured during Cloning of repo\n\tpath: %s\n\terror: %v\n\n\n", project_name, path_to_dir, _repo_dir_err)
				aborted_projects += fmt.Sprintf("gitclone fail: %s\n", project_name)
				continue
			}

			var packages []*PackageCounter

			walk_err := filepath.Walk(path_to_dir, func(path string, info os.FileInfo, walk_func_err error) error {
				// called on each file visited in the walk
				// if file cannot be reached and throws error
				if walk_func_err != nil {
					WarningLog("FileWalk of %s, file path could not be accessed.\n\tpath: %s\n\terror: %v\n", project_name, path_to_dir, walk_func_err)
					return walk_func_err
				}
				// check if folder that needs to be explored
				if info.IsDir() {
					// dir to be skipped
					if info.Name() == "vendor" || info.Name() == "tests" || info.Name() == "test" {
						DebugLog("Skipping dir in FileWalk of %s due to its name: \"%+v\"\n", project_name, info.Name())
						return filepath.SkipDir
					}
					// dir to be explored
					VerboseLog("FileWalk of %s, found dir to explore: %+v\n", project_name, info.Name())
					var new_counter PackageCounter = ParseDir(proj_name, path, path_to_dir)
					packages = append(packages, &new_counter)
					return nil
				}
				// just a file
				VerboseLog("FileWalk of %s, found file: %+v\n", project_name, info.Name())
				return nil
			})

			// if any errors on file walk
			if walk_err != nil {
				// skip analysis
				FailureLog("Aborting project \"%s\". Error occured during FileWalk of files in repo.\n\tpath: %s\n\terror: %s\n\n\n", project_name, path_to_dir, walk_err)
				aborted_projects += fmt.Sprintf("filewalk fail: %s\n", project_name)
				continue
			} else {
				DebugLog("FileWalk of \"%s\" yielded no errors.\n", project_name)
			}

			if keep_repos {
				GeneralLog("Skipping result output, as repos are being saved, and these would fail anyway. Saving time.\n\n")
			} else {

				// create html results
				var project_counter Counter = HtmlOutputCounters(packages, commit_hash, project_name, index_data, path_to_dir) // html

				// analysis
				project_counter = ParseConcurrencyPrimitives(path_to_dir, project_counter) // analyses occurences of Waitgroup,mutexes and operations on them

				// create csv results
				CsvOutputCounters(project_name, packages, path_to_dir, project_counter) // csvs
				// project_counters = append(project_counters, project_counter)

			}
			// remove repo as we go if enabled
			_temp_proj_dir := filepath.Join(clone_dir, ProjectName(project_name))
			_, repo_err := os.Stat(_temp_proj_dir)
			if repo_err == nil {
				_temp_dir_size, size_err := DirSize(_temp_proj_dir)
				repo_dir_size_mb := int64(math.Ceil(float64(_temp_dir_size) / float64(1000000)))
				repo_memory_used += repo_dir_size_mb
				if keep_repos {
					if size_err != nil {
						FailureLog("Failed to calculate \"%s\" total size...\n\terror: %v\n", _temp_proj_dir, size_err)
					} else {
						repo_memory_percent := int64(math.Ceil(float64(repo_memory_used) / float64(repo_memory_capacity) * 100))
						GeneralLog("Finished Project %d/%d: %s\n\trepo path: %s\n\trepo size: %d MB\n\tcurrent total: %d MB\n\n", _index+1, projects_to_process, project_name, _temp_proj_dir, repo_dir_size_mb, repo_memory_used)
						GeneralLog("Memory used: (%d MB/%d MB), %v %%\n\n\n", repo_memory_used, repo_memory_capacity, repo_memory_percent)
						if repo_memory_limited && repo_memory_used > repo_memory_capacity {
							GeneralLog("Memory capacity exceeded, skipping the remaining %d repos\n\n\n", projects_to_process-_index)
							break
						}
					}
				} else {
					DebugLog("Deleting repo for project at %s\n", _temp_proj_dir)
					if repo_rem_err := os.RemoveAll(_temp_proj_dir); repo_rem_err != nil {
						FailureLog("Error occured trying to delete repo dir: %s\n\terror: %v", _temp_proj_dir, repo_rem_err)
					}
					GeneralLog("Finished Project %d/%d: %s\n\n\n", _index+1, projects_to_process, project_name)
				}
			} else {
				WarningLog("Something has happened to the projects repo.\n")
				GeneralLog("Finished Project %d/%d: %s\n\n\n", _index+1, projects_to_process, project_name)
			}
			// for debugging when dev
			if single_run {
				DebugLog("Finishing after single run\n")
				break
			} else {
				DebugLog("Continuing through user/repo projects\n")
			}
		} else {
			if projects_to_process == 0 {
				WarningLog("Skipping %d, \"%s\": Unable to read project from \"projects.txt\"\n", _index, project_name)
			} else {
				GeneralLog("Skipping remaining %d projects due to limit of %d\n", total_project_count-_index, projects_to_process)
				break
			}
		}
	}
	createIndexFile(index_data) // index html

	// check if temp.go is still there
	if _info, temp_err := os.Stat(temp_go); temp_err == nil && !_info.IsDir() {
		DebugLog("File \"%s\" was not deleted, removing now.\n", temp_go)
		t_err := os.Remove(temp_go)
		if t_err == nil {
			DebugLog("Successfully deleted \"%s\"\n", temp_go)
		} else {
			FailureLog("Error: Unable to delete \"%s\", will need to be done manually...\n\t%v\n", temp_go, t_err)
		}
	} else {
		DebugLog("File \"%s\" was already deleted.\n", temp_go)
	}

	DebugLog("Total Logs: %d\n", total_logs)

	aborted_count := len(strings.Split(aborted_projects, "\n")) - 1
	GeneralLog("Total number of projects aborted: %d\n%s\n", aborted_count, aborted_projects)
	GeneralLog("Total number of projects succeeded: %d\n\n", total_project_count-aborted_count)

	if keep_repos {
		GeneralLog("Total size of \"%s\": %d MB\n", clone_dir, repo_memory_used)
	} else {
		GeneralLog("Total size of Repos downloaded, now deleted: %d MB\n", repo_memory_used)
	}

	GeneralLog("Finished. Results can be found...\n\tlocal path: %s\n\tglobal path: %s\n\n", result_dir, GenerateFullPath(result_dir))

}

func createIndexFile(index_data *IndexFileData) {
	f, err := os.Create("index.html")
	if err != nil {
		PanicLog(err, "Unable to create \"index.html\" file.\n")
	} else {
		GeneralLog("successfuly created index.html\n")
	}
	tmpl := template.Must(template.ParseFiles("html\\index_layout.html"))
	tmpl.Execute(f, index_data) // write the index page
}

func GenerateWDPath() string {
	wd, err := os.Getwd()
	if err != nil {
		PanicLog(err, "Unable to work out the current working directory.\n")
	}
	return wd
}

func GenerateFullPath(local_path string) string {
	return filepath.Join(GenerateWDPath(), local_path)
}

func PathIsGlobal(path string) bool {
	if strings.Contains(path, GenerateWDPath()) {
		return true
	} else {
		return false
	}
}

func ProjectName(project_name string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(string(project_name), "/", "___"), "-", "_"), "\\", "")
}

func ProjectURL(project_name string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(string(project_name), "___", "/"), "_", "-"), "\\", "/")
}

func EnsureDir(_path string, build_parents bool) (string, int, bool, error) {
	// for consistency
	_path = strings.ReplaceAll(_path, "/", "\\")
	// check already exists
	if _, err1 := os.Stat(_path); err1 != nil {
		if os.IsNotExist(err1) {
			// try make dir
			os.Mkdir(_path, dir_mode)
			if _, err2 := os.Stat(_path); err2 != nil {
				if os.IsNotExist(err2) {
					// path not found, build parents
					if build_parents {
						outer_dir := filepath.Base(_path)
						parent_dir, depth, success, parent_err := EnsureDir(strings.TrimSuffix(_path, "\\"+outer_dir), true)
						// sanity check
						if _path != parent_dir+"\\"+outer_dir {
							WarningLog("EnsureDir,  %s: After creating the parent folders, the path returned doesn't match...\n\tparent path: %s\n", _path, filepath.Join(parent_dir, outer_dir))
						}
						// check successful recursion
						if success && parent_err == nil {
							// success
							os.Mkdir(_path, dir_mode)
							// final check
							if _, err3 := os.Stat(_path); err3 != nil {
								// check why error but doesnt exist
								FailureLog("EnsureDir,  %s: Unrecoverable error occured. Unable to create dir despite parent folders being created...\n\tdepth: %d\n\terror: %v\n", _path, depth, err3)
								return _path, depth + 1, false, err3
							} else {
								// success !!!
								DebugLog("EnsureDir,  %s: Success, created %d parents too.\n", _path, depth)
								return _path, depth + 1, true, nil
							}
						} else {
							// a parent failed, pass it on
							FailureLog("EnsureDir,  %s: Parent unable to be created, pass it on...\n\tparent path: %s\n\tbuild parents: %t\n\terror: %v\n", _path, parent_dir, build_parents, parent_err)
							return parent_dir, depth + 1, false, nil
						}
					}
					FailureLog("EnsureDir,  %s: Unable to create dir, and parents haven't been told to be created...\n\tbuild parents: %t\n\terror: %v\n", _path, build_parents, err2)
					return _path, 1, false, err2
				} else {
					FailureLog("EnsureDir,  %s: Unrecoverable error occured. Unable to create dir...\n\tbuild parents: %t\n\terror: %v\n", _path, build_parents, err2)
					return _path, 1, false, err2
				}
			} else {
				DebugLog("EnsureDir,  %s: Success.\n", _path)
				return _path, 1, true, nil
			}
		} else {
			FailureLog("EnsureDir,  %s: Unrecoverable error occured. Unable to check if dir exists...\n\tbuild parents: %t\n\terror: %v\n", _path, build_parents, err1)
			return _path, 1, false, err1
		}
	} else {
		// already exists, return
		DebugLog("EnsureDir,  %s: Dir already exists.\n", _path)
		return _path, 1, true, nil
	}
}

func GeneratePackageListFiles(path_to_dir string) string {
	git_cmd := exec.Command("ls")
	git_cmd.Dir = path_to_dir
	var git_out bytes.Buffer
	git_cmd.Stdout = &git_out
	err := git_cmd.Run()
	if err != nil {
		WarningLog("Main, GPLF: Error while running git ls-files: %v\n", err)
	}
	filenames := ""
	for _, name := range strings.Split(git_out.String(), "\n") {
		if strings.HasSuffix(name, ".go") {
			filenames += filepath.Join(path_to_dir, name) + "\n"
		}
	}
	// replace with forward for url output
	return strings.ReplaceAll(filenames, "\\", "/")
}

func ReadNumberOfLines(list_filenames string) int {
	var xargs_out bytes.Buffer
	var git_out bytes.Buffer
	filenames := strings.Split(list_filenames, "\n")
	// return if no files to count
	if len(filenames) == 0 {
		DebugLog("Main, RNoL: File list was empty.\n")
		return 0
	}

	git_out.Reset()
	for _, filename := range strings.Split(list_filenames, "\n") {
		if filename != "" {
			git_out.WriteString("\"" + filename + "\"\n")
		}
	}
	xargs_cmd := exec.Command("xargs", "cat")
	xargs_cmd.Stdin = &git_out
	xargs_cmd.Stdout = &xargs_out
	xargs_err := xargs_cmd.Run()
	if xargs_err != nil {
		FailureLog("Main, RNoF: Error running cat...\n\terror: %v\n", xargs_err)
	}

	f, _ := os.Create(temp_go)
	f.Write(xargs_out.Bytes())
	var wc_out bytes.Buffer
	wc_cmd := exec.Command("cloc", temp_go, "--csv")
	// wc_cmd.Stdin = &xargs_out
	wc_cmd.Stdout = &wc_out
	err3 := wc_cmd.Run()
	if err3 != nil {
		FailureLog("Main, RNoF: Error while running word count...\n\terror: %v\n", err3)
	}
	defer os.Remove(temp_go)
	defer f.Close()
	word_count := strings.Split(strings.TrimSpace(wc_out.String()), "\n")
	cloc_infos := strings.Split(strings.TrimSpace(word_count[len(word_count)-1]), ",")

	// if the command isnt blank
	if len(cloc_infos) >= 5 {
		// in the csv, from 0, index 4 is lines of code
		num, _ := strconv.Atoi(cloc_infos[4])
		return num
	} else {
		WarningLog("Main, RNoL: Command \"cloc\" may not be available on your system.\n\tCommand \"cloc\": Count Lines Of Code\n\tlink: https://github.com/AlDanial/cloc\n\tinstall: npm install -g cloc\n")
		return 0
	}
}

// https://stackoverflow.com/questions/32482673/how-to-get-directory-total-size
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func InitDirs() {
	// create directories
	if _csv_path, _csv_depth, _csv_success, _csv_err := EnsureDir(csv_result_dir, true); !_csv_success && _csv_err != nil {
		WarningLog("Unable to make the dir needed to store the resulting CSV files...\n\tpath: %s\n\tdepth: %d\n\tsuccess: %t\n\terror: %v\n", _csv_path, _csv_depth, _csv_success, _csv_err)
	}
	if _html_path, _html_depth, _html_success, _html_err := EnsureDir(html_result_dir, true); !_html_success && _html_err != nil {
		WarningLog("Unable to make the dir needed to store the resulting HTML files...\n\tpath: %s\n\tdepth: %d\n\tsuccess: %t\n\terror: %v\n", _html_path, _html_depth, _html_success, _html_err)
	}
	if _clone_path, _clone_depth, _clone_success, _clone_err := EnsureDir(clone_dir, true); !_clone_success && _clone_err != nil {
		WarningLog("Unable to make the dir needed to store the files of the repos this tool will download...\n\tpath: %s\n\tdepth: %d\n\tsuccess: %t\n\terror: %v\n", _clone_path, _clone_depth, _clone_success, _clone_err)
	} else {
		if !keep_repos && projects_path != "test" {
			defer os.RemoveAll(clone_dir)
			GeneralLog("Marked the Clone Directory to be deleted after run.\n")
		}
	}
}
