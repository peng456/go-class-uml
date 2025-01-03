package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	// "../../mermaidgo"
	"github.com/peng456/goclassuml/mermaidgo"

	goplantuml "github.com/peng456/goclassuml/parser"
)

const versionFlag = "0.1"

// RenderingOptionSlice will implements the sort interface
type RenderingOptionSlice []goplantuml.RenderingOption

// Len is the number of elements in the collection.
func (as RenderingOptionSlice) Len() int {
	return len(as)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (as RenderingOptionSlice) Less(i, j int) bool {
	return as[i] < as[j]
}

// Swap swaps the elements with indexes i and j.
func (as RenderingOptionSlice) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
}

func main() {

	// 判断 入参是否有 --version
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-version" || os.Args[1] == "version") {
		fmt.Println("version:", versionFlag)
		return
	}
	recursive := flag.Bool("recursive", false, "walk all directories recursively")
	ignore := flag.String("ignore", "", "comma separated list of folders to ignore")
	showAggregations := flag.Bool("show-aggregations", false, "renders public aggregations even when -hide-connections is used (do not render by default)")
	hideFields := flag.Bool("hide-fields", false, "hides fields")
	hideMethods := flag.Bool("hide-methods", true, "hides methods")
	hideConnections := flag.Bool("hide-connections", false, "hides all connections in the diagram")
	showCompositions := flag.Bool("show-compositions", false, "Shows compositions even when -hide-connections is used")
	showImplementations := flag.Bool("show-implementations", false, "Shows implementations even when -hide-connections is used")
	showAliases := flag.Bool("show-aliases", false, "Shows aliases even when -hide-connections is used")
	showConnectionLabels := flag.Bool("show-connection-labels", true, "Shows labels in the connections to identify the connections types (e.g. extends, implements, aggregates, alias of")
	title := flag.String("title", "", "Title of the generated diagram")
	notes := flag.String("notes", "", "Comma separated list of notes to be added to the diagram")
	outputType := flag.Int("ot", 1, "output type.  1: app  2: string, default 1")
	saveType := flag.String("st", "svg", "save diagram type : svg or png ")
	scale := flag.Int("scale", 0, "0 or 2.0")
	showOptionsAsNote := flag.Bool("show-options-as-note", false, "Show a note in the diagram with the none evident options ran with this CLI")
	aggregatePrivateMembers := flag.Bool("aggregate-private-members", false, "Show aggregations for private members. Ignored if -show-aggregations is not used.")
	hidePrivateMembers := flag.Bool("hide-private-members", false, "Hide private fields and methods")
	flag.Parse()

	renderingOptions := map[goplantuml.RenderingOption]interface{}{
		goplantuml.RenderConnectionLabels:  *showConnectionLabels,
		goplantuml.RenderFields:            !*hideFields,
		goplantuml.RenderMethods:           !*hideMethods,
		goplantuml.RenderAggregations:      *showAggregations,
		goplantuml.RenderAliases:           *showAliases,
		goplantuml.RenderTitle:             *title,
		goplantuml.AggregatePrivateMembers: *aggregatePrivateMembers,
		goplantuml.RenderPrivateMembers:    !*hidePrivateMembers,
	}
	if *hideConnections {
		renderingOptions[goplantuml.RenderAliases] = *showAliases
		renderingOptions[goplantuml.RenderCompositions] = *showCompositions
		renderingOptions[goplantuml.RenderImplementations] = *showImplementations

	}
	noteList := []string{}
	if *showOptionsAsNote {
		legend, err := getLegend(renderingOptions)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		noteList = append(noteList, legend)
	}
	if *notes != "" {
		noteList = append(noteList, "", "<b><u>Notes</u></b>")
	}
	split := strings.Split(*notes, ",")
	for _, note := range split {
		trimmed := strings.TrimSpace(note)
		if trimmed != "" {
			noteList = append(noteList, trimmed)
		}
	}
	renderingOptions[goplantuml.RenderNotes] = strings.Join(noteList, "\n")
	dirs, err := getDirectories()

	if err != nil {
		fmt.Println("usage:\ngoClassDiagrams <DIR>\nDIR Must be a valid directory")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	ignoredDirectories, err := getIgnoredDirectories(*ignore)
	if err != nil {

		fmt.Println("usage:\ngoClassDiagrams [-ignore=<DIRLIST>]\nDIRLIST Must be a valid comma separated list of existing directories")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	result, err := goplantuml.NewClassDiagram(dirs, ignoredDirectories, *recursive)
	result.SetRenderingOptions(renderingOptions)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	rendered := result.Render()

	ctx := context.Background()
	// mermaid.initialize({
	//   maxTextSize: 90000
	// });
	re, _ := mermaidgo.NewRenderEngine(ctx)
	defer re.Cancel()

	// 输出方式 0 字符串直接输出 1 文件（默认打开）

	var writer io.Writer
	if *outputType == 2 {
		writer = os.Stdout

		fmt.Fprint(writer, rendered)
		return
	}

	content := rendered
	// 根据 参数 获取 name
	// fileName := "go-class-mermaid.mm"
	fileName, _ := getSaveFileName()
	tempDir := os.TempDir()
	// fmt.Println("User temp directory:", tempDir)
	outFileName := tempDir + fileName
	fileSavePath := ""
	switch *saveType {
	case "svg":
		svg_content, _ := re.Render(content)
		fileSavePath = outFileName + ".svg"
		os.WriteFile(fileSavePath, []byte(svg_content), 0644)

	case "png":
		if *scale == (int)(0) {
			// get the result as PNG bytes
			png_in_bytes, _, _ := re.RenderAsPng(content)
			fileSavePath = outFileName + ".png"

			os.WriteFile(fileSavePath, png_in_bytes, 0644)

		} else {
			fileSavePath = outFileName + "_scaled.png"

			scaled_png_in_bytes, _, _ := re.RenderAsScaledPng(content, 2.0)
			os.WriteFile(fileSavePath, scaled_png_in_bytes, 0644)
		}
	}

	// 是否打开文件 默认 vscode 或者调用系统 选择对应的软件打开
	errOpenFile := openFileWithVSCode(fileSavePath)
	if errOpenFile != nil {
		fmt.Println("Error:", errOpenFile)
	} else {
		fmt.Println("File opened successfully")
	}

}

func openFileWithVSCode(filePath string) error {
	// 检查 VS Code 是否安装
	_, err := exec.LookPath("code")
	if err != nil {
		// 如果没有安装 VS Code，则使用系统默认程序打开文件
		return openFile(filePath)
	}

	// 使用 VS Code 打开文件
	cmd := exec.Command("code", filePath)
	return cmd.Start()
}

func openFile(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// macOS
		// 检查 VS Code 是否安装
		openApp := "open"
		_, err := exec.LookPath("code")
		if err == nil {
			openApp = "code"
		}
		cmd = exec.Command(openApp, filePath)
	case "linux":
		// Linux
		cmd = exec.Command("xdg-open", filePath)
	case "windows":
		// Windows
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", filePath)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

func getDirectories() ([]string, error) {

	args := flag.Args()
	if len(args) < 1 {
		return nil, errors.New("DIR missing")
	}
	dirs := []string{}
	for _, dir := range args {
		fi, err := os.Stat(dir)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("could not find directory %s", dir)
		}

		if !fi.Mode().IsDir() {
			// 判断是否是 go 文件
			if strings.HasSuffix(fi.Name(), ".go") {
				dirs = append(dirs, dir)

				break
			}
			return nil, fmt.Errorf("%s is not a directory", dir)

		}
		dirAbs, err := filepath.Abs(dir)
		if err != nil {
			return nil, fmt.Errorf("could not find directory %s", dir)
		}
		dirs = append(dirs, dirAbs)
	}
	return dirs, nil
}

func getSaveFileName() (string, error) {

	// return "go-class-mermaid.mm", nil
	args := flag.Args()

	name := ""
	for _, dir := range args {
		fi, _ := os.Stat(dir)
		if fi.Mode().IsDir() {
			// 字符串 先 spilt 然后获取最后一个
			// name = strings.Split(strings.Trim(dir, "/"), "/").pop()
			nameArr := strings.Split(strings.Trim(dir, "/"), "/")
			name = nameArr[len(nameArr)-1]
		} else {
			name = strings.TrimSuffix(fi.Name(), ".go")
		}

	}
	return name, nil
}

func getIgnoredDirectories(list string) ([]string, error) {
	result := []string{}
	list = strings.TrimSpace(list)
	if list == "" {
		return result, nil
	}
	split := strings.Split(list, ",")
	for _, dir := range split {
		dirAbs, err := filepath.Abs(strings.TrimSpace(dir))
		if err != nil {
			return nil, fmt.Errorf("could not find directory %s", dir)
		}
		result = append(result, dirAbs)
	}
	return result, nil
}

func getLegend(ro map[goplantuml.RenderingOption]interface{}) (string, error) {
	result := "<u><b>Legend</b></u>\n"
	orderedOptions := RenderingOptionSlice{}
	for o := range ro {
		orderedOptions = append(orderedOptions, o)
	}
	sort.Sort(orderedOptions)
	for _, option := range orderedOptions {
		val := ro[option]
		switch option {
		case goplantuml.RenderAggregations:
			result = fmt.Sprintf("%sRender Aggregations: %t\n", result, val.(bool))
		case goplantuml.RenderAliases:
			result = fmt.Sprintf("%sRender Connections: %t\n", result, val.(bool))
		case goplantuml.RenderCompositions:
			result = fmt.Sprintf("%sRender Compositions: %t\n", result, val.(bool))
		case goplantuml.RenderFields:
			result = fmt.Sprintf("%sRender Fields: %t\n", result, val.(bool))
		case goplantuml.RenderImplementations:
			result = fmt.Sprintf("%sRender Implementations: %t\n", result, val.(bool))
		case goplantuml.RenderMethods:
			result = fmt.Sprintf("%sRender Methods: %t\n", result, val.(bool))
		case goplantuml.AggregatePrivateMembers:
			result = fmt.Sprintf("%sPrivate Aggregations: %t\n", result, val.(bool))
		}
	}
	return strings.TrimSpace(result), nil
}
