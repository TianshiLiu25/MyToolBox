package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DependencyGraph 用于表示依赖关系图
type DependencyGraph map[string][]string

// collectDependencies 递归收集指定路径下的 .so 文件的依赖关系
func collectDependencies(soPath string, graph DependencyGraph) error {
	err := filepath.Walk(soPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".so") {
			dependencies, err := getSODependencies(path)
			if err != nil {
				return err
			}


			graph[filepath.Base(path)] = dependencies
		}

		return nil
	})

	return err
}

// getSODependencies 使用 readelf 获取指定 .so 文件的依赖关系
func getSODependencies(soPath string) ([]string, error) {
	cmd := exec.Command("readelf", "-d", soPath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	dependencies := make([]string, 0)

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "NEEDED") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				dependency := strings.Trim(fields[4], "[]")
				dependencies = append(dependencies, dependency)
				// fmt.Printf("%v -> %v\n", soPath, dependency)
			}
		}
	}

	return dependencies, nil
}

// findDependencyPath 查找两个 .so 文件之间的依赖路径
func findDependencyPath(graph DependencyGraph, depender, dependee string, visited map[string]bool, path []string) bool {
	if depender == dependee {
		path = append(path, dependee)
		fmt.Println(strings.Join(path, " -> "))
		return true
	}

	visited[depender] = true
	path = append(path, depender)

	for _, dependency := range graph[depender] {
		if !visited[dependency] {
			if findDependencyPath(graph, dependency, dependee, visited, path) {
				return true
			}
		}
	}

	return false
}

// showDependencyTree 以树形结构展示所有依赖
func showDependencyTree(graph DependencyGraph, root string, level int) {
	fmt.Printf("%s%s\n", strings.Repeat("  ", level), root)

	for _, dependency := range graph[root] {
		showDependencyTree(graph, dependency, level+1)
	}
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run script.go --search-path <so_path> (--depender <depender> --dependee <dependee> | --show-dependence-of <dependency>)")
		os.Exit(1)
	}

	// 解析命令行参数
	var soPath, depender, dependee, showDependencyOf string
	var showDependencyTreeMode bool

	for i := 1; i < len(os.Args)-1; i += 2 {
		switch os.Args[i] {
		case "--search-path":
			soPath = os.Args[i+1]
		case "--depender":
			depender = os.Args[i+1]
		case "--dependee":
			dependee = os.Args[i+1]
		case "--show-dependence-of":
			showDependencyOf = os.Args[i+1]
			showDependencyTreeMode = true
		}
	}

	if soPath == "" || (depender == "" && dependee == "" && showDependencyOf == "") {
		fmt.Println("Invalid arguments.")
		os.Exit(1)
	}

	graph := make(DependencyGraph)

	err := collectDependencies(soPath, graph)
	if err != nil {
		fmt.Println("Error collecting dependencies:", err)
		os.Exit(1)
	}

	if showDependencyTreeMode {
		showDependencyTree(graph, showDependencyOf, 0)
	} else {
		visited := make(map[string]bool)
		path := make([]string, 0)

		if !findDependencyPath(graph, depender, dependee, visited, path) {
			fmt.Printf("No dependency path found from %s to %s\n", depender, dependee)
		}
	}
}
