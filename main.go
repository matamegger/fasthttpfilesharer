package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

var rootelements []string

func main() {
	fmt.Println("Starting FastFileSharer")
	fmt.Println("Sharing content:")

	//Generating list of dirs and files to share
	// if non is given, the workingDir will be used
	var dirs []string
	var files []string

	//Command line arguments start at index 1 (= len>1)
	if len(os.Args) > 1 {
		for _, temp := range os.Args[1:] {
			//evaluating if proper file
			file, err := os.Stat(temp)
			if err != nil {
				continue
			}
			if file.IsDir() {
				dirs = append(dirs, temp)
			} else {
				files = append(files, temp)
			}
		}
	} else {
		workingDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		dirs = append(dirs, workingDir)
	}

	//adding handlers
	if len(files) > 0 || len(dirs) > 1 {
		http.HandleFunc("/", handler)
		addDirHandlerFromArray(&dirs)
		addFileHandlerFromArray(&files)
	} else {
		fs := http.FileServer(http.Dir(dirs[0]))
		http.Handle("/", fs)
		fmt.Println("/" + "\t" + dirs[0])
	}

	fmt.Println("\nOpening Server on Port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("\nCouldn't open Server on Port 80")
	}

}

func addFileHandlerFromArray(files *[]string) {
	for _, file := range *files {
		fs := http.FileServer(http.Dir(filepath.Dir(file) + "/"))
		fileBase := filepath.Base(file)
		rootelements = append(rootelements, fileBase)
		http.Handle("/"+fileBase, fs)
		fmt.Println("/" + fileBase + "\t" + file)
	}
}

func addDirHandlerFromArray(dirs *[]string) {
	for _, dir := range *dirs {
		fs := http.FileServer(http.Dir(dir))
		fileBase := filepath.Base(dir)
		rootelements = append(rootelements, fileBase)
		http.Handle("/"+fileBase+"/", http.StripPrefix("/"+fileBase+"/", fs))
		fmt.Println("/" + fileBase + "\t" + dir)
	}
}

func handler(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<pre>")
	for _, dir := range rootelements {
		fmt.Fprintf(w, "<a href=\"%[1]v\">%[1]v</a>\n", dir)
	}
	fmt.Fprint(w, "</pre>")

}
