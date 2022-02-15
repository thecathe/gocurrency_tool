package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const csv_result_dir = "_results\\csv"

func CsvOutputCounters(_project_name string, package_counters []*PackageCounter, path_to_dir string, project_counter Counter) {

	project_name := _project_name
	GeneralLog("%s: Entering CSV output stage.\n", project_name)

	var csv_path = GenerateCSVPath()

	if _, err := os.Stat(csv_path); os.IsNotExist(err) && err != nil {
		PanicLog(err, "Csv, OutputCounters: CSV output dir \"%s\" does not exist. Entering panic\n", csv_path)
	}

	// adjsut for separate result files
	if separate_results {
		project_name = _project_name + "_info"
	}

	var file_name = ProjectName(project_name) + ".csv"
	var file_path string = filepath.Join(csv_result_dir, file_name)
	f, err := os.Create(file_path)
	if err != nil {
		if os.IsExist(err) {
			ExitLog(0, "Csv, OutputCounters: Results already exist, either enable overwrite or move the current results.\n")
		} else {
			PanicLog(err, "Csv, OutputCounters: Unable to create file \"%s\", entering panic...\n", file_path)
		}
	}
	var num_of_packages int = 0

	num_featured_files := 0
	num_files := 0
	// num_featured_packages := 0
	DebugLog("Csv, OutputCounters: Detecting how many packages in %d features...\n", len(package_counters))
	for _, counter := range package_counters {
		if len(counter.File_counters) > 0 {
			num_of_packages++
		}
		num_featured_files += counter.Featured_files
		num_files += counter.Num_files
	}
	if len(package_counters) > 0 {
		// header
		f.WriteString("Data, Features, Total\n")
		f.WriteString(fmt.Sprintf("Line num, , %d\n", project_counter.Line_number))
		f.WriteString(fmt.Sprintf("packages num, %d, %d\n", project_counter.Num_of_packages_with_features, num_of_packages))
		f.WriteString(fmt.Sprintf("files num, %d, %d\n", num_featured_files, num_files))
		f.WriteString(fmt.Sprintf("Line num per featured file, , %.d\n", readNumberOfLinesPerFeaturedFile(package_counters)))

		VerboseLog("Line Num: %d\n", project_counter.Line_number)
		VerboseLog("Packages Num: %d, %d\n", project_counter.Num_of_packages_with_features, num_of_packages)
		VerboseLog("Files Num: %d, %d\n", num_featured_files, num_files)
		VerboseLog("Average Line Num: %.d\n", readNumberOfLinesPerFeaturedFile(package_counters))
	} else {
		WarningLog("Csv, OutputCounters: Package counters was empty...\n\tpackage counter length: %d\n", len(package_counters))
	}

	if separate_results {
		f.Close()
		GeneralLog("Csv, OutputCounters: Results file \"%s\" is finished.\n", file_name)
		project_name = _project_name + "_package_lines"
		file_name = ProjectName(project_name) + ".csv"

		file_path = filepath.Join(csv_result_dir, file_name)
		f, err = os.Create(file_path)
		if err != nil {
			if os.IsExist(err) {
				ExitLog(0, "Csv, OutputCounters: Results already exist, either enable overwrite or move the current results.\n")
			} else {
				PanicLog(err, "Csv, OutputCounters: Unable to create file \"%s\", entering panic...\n", file_path)
			}
		}
		f.WriteString("Package Name, Number of Lines\n")
	}

	DebugLog("Csv, OutputCounters: ReadingNumberOfLines in %d packages.\n", len(package_counters))
	for _, counter := range package_counters {
		if len(counter.File_counters) > 0 {
			number_of_lines := ReadNumberOfLines(GeneratePackageListFiles(counter.Counter.Package_path))
			f.WriteString(fmt.Sprintf("%s,%d\n", strings.ReplaceAll(counter.Counter.Package_name, "/", "\\"), number_of_lines))
		}
	}

	if separate_results {
		f.Close()
		GeneralLog("Csv, OutputCounters: Results file \"%s\" is finished.\n", file_name)
		project_name = _project_name + "_file_data"
		file_name = ProjectName(project_name) + ".csv"

		file_path = filepath.Join(csv_result_dir, file_name)
		f, err = os.Create(file_path)
		if err != nil {
			if os.IsExist(err) {
				ExitLog(0, "Csv, OutputCounters: Results already exist, either enable overwrite or move the current results.\n")
			} else {
				PanicLog(err, "Csv, OutputCounters: Unable to create file \"%s\", entering panic...\n", file_path)
			}
		}
	}

	f.WriteString("Filename, Line #, Feature #, Concurrent Type, Additional Feature Info, Package\n")

	DebugLog("Csv, OutputCounters: Writing CSV file.\n")
	for _, feature := range project_counter.Features {
		f.WriteString(fmt.Sprintf("%s,%v,%v,%v,%v,%v\n",
			strings.ReplaceAll(feature.F_filename, "/", "\\"),
			feature.F_line_num,
			feature.F_type_num,
			feature.F_type,
			feature.F_number,
			feature.F_package_name))
	}
	GeneralLog("Csv, OutputCounters: Results file \"%s\" is finished.\n", file_name)
	f.Close()
	GeneralLog("Csv, OutputCounters: Finished %s\n\n", project_name)
}

func readNumberOfLinesPerFeaturedFile(package_counters []*PackageCounter) int {
	var git_out bytes.Buffer
	var xargs_out bytes.Buffer

	var filenames []string

	for _, package_counter := range package_counters {
		for _, counter := range package_counter.File_counters {
			if len(counter.Features) > 0 {
				filenames = append(filenames, counter.filename)
			}
		}
	}
	for _, filename := range filenames {
		if filename != "" {
			git_out.WriteString("\"" + filename + "\"\n")
		}
	}
	xargs_cmd := exec.Command("xargs", "cat")
	xargs_cmd.Stdin = &git_out
	xargs_cmd.Stdout = &xargs_out
	xargs_cmd.Run()

	f, _ := os.Create(temp_go)
	f.Write(xargs_out.Bytes())
	var wc_out bytes.Buffer
	wc_cmd := exec.Command("cloc", temp_go, "--csv")
	// wc_cmd.Stdin = &xargs_out
	wc_cmd.Stdout = &wc_out
	err3 := wc_cmd.Run()
	if err3 != nil {
		FailureLog("Csv, RNoLpFF: Error while running wc: %v\n", err3)
	}
	defer os.Remove(temp_go)
	defer f.Close()
	word_count := strings.Split(strings.TrimSpace(wc_out.String()), "\n")
	cloc_infos := strings.Split(strings.TrimSpace(word_count[len(word_count)-1]), ",")

	if len(cloc_infos) >= 5 {
		num, _ := strconv.Atoi(cloc_infos[4])

		if len(filenames) == 0 {
			return 0
		}

		return num
	} else {
		return 0.0
	}

}

func GenerateCSVPath() string {
	return filepath.Join(GenerateWDPath(), csv_result_dir)
}
