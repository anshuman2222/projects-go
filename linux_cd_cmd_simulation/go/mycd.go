package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func remove_index_from_slice_of_strings(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func panic_recovery(err error) {
	if err := recover(); err != nil { //catch
		fmt.Println("mycd: ", err)
		os.Exit(1)
	}
}

func is_alphanum(word string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(word)
}

func remove_empty_index_in_slice(slice []string, from_index int) []string {
	for i := from_index; i < len(slice); {
		if strings.Compare(slice[i], "") == 0 {
			slice = remove_index_from_slice_of_strings(slice, i)
		} else {
			break
		}
	}
	return slice
}

func validate_and_convert_dir_path_to_slice(dir_path string) []string {

	dir_path_list := strings.Split(dir_path, "/")

	if strings.Compare(dir_path_list[0], "") == 0 {
		dir_path_list = remove_empty_index_in_slice(dir_path_list, 1)
	}

	for i, dir_name := range dir_path_list {
		// Validate if dir name is alphanumeric or not
		if !(is_alphanum(dir_name) ||
			(strings.Compare(dir_name, "..") == 0) ||
			(strings.Compare(dir_name, ".") == 0)) {
			panic(dir_name + ": No such file or directory")
		}
		// Remove extra backslashes if persists (In strings.Split with "/"
		// converts backslashes to empty string)
		if strings.Compare(dir_name, "") == 0 {
			if i == 0 {
				// ignore first backslash if exists
				continue
			}
			dir_path_list = remove_empty_index_in_slice(dir_path_list, i)
		}
	}

	return dir_path_list
}

func back_to_the_child_path(cwd []string) []string {
	return append(cwd[:len(cwd)-1], cwd[len(cwd):]...)
}

func prepare_new_dir(cwd_dir_path_list, new_dir_path_list []string) string {

	new_dir := "/"

	if strings.Compare(new_dir_path_list[0], "") == 0 {
		// If target path starts with "/". Checking with empty string because In strings.Split with "/"
		// "/" converts backslash to empty string in slice of strings
		if len(new_dir_path_list) == 1 && strings.Compare(new_dir_path_list[0], "") == 0 {
			// If target has only backslash "/"
			return new_dir
		}
		cwd_dir_path_list = nil
		cwd_dir_path_list = append(cwd_dir_path_list, "/")
	}

	for i := 0; i < len(new_dir_path_list); i++ {
		if strings.Compare(new_dir_path_list[i], "") == 0 ||
			strings.Compare(new_dir_path_list[i], ".") == 0 {
			continue
		}
		if strings.Compare(new_dir_path_list[i], "..") == 0 {
			if len(cwd_dir_path_list) > 1 {
				cwd_dir_path_list = back_to_the_child_path(cwd_dir_path_list)
			}
		} else {
			cwd_dir_path_list = append(cwd_dir_path_list, new_dir_path_list[i])
		}
	}
	if len(cwd_dir_path_list) == 1 && strings.Compare(cwd_dir_path_list[0], "") == 0 {
		return new_dir
	} else {
		new_dir = strings.Join(cwd_dir_path_list, "/")
		if (len(new_dir) > 1) &&
			strings.Compare(string(new_dir[0]), "/") == 0 &&
			strings.Compare(string(new_dir[1]), "/") == 0 {
			new_dir = new_dir[1:]
		}
		return new_dir
	}
}

func main() {
	defer panic_recovery(nil)

	if len(os.Args) < 3 {
		panic("too few arguments")
	} else if len(os.Args) > 3 {
		panic("too many arguments")
	}

	cwd := os.Args[1]
	to_dirctory := os.Args[2]

	cwd_dir_path_list := validate_and_convert_dir_path_to_slice(cwd)
	new_dir_path_list := validate_and_convert_dir_path_to_slice(to_dirctory)

	new_dir := prepare_new_dir(cwd_dir_path_list, new_dir_path_list)

	fmt.Println(new_dir)
}
