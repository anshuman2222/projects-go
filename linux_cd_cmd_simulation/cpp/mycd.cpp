#include <iostream>
#include <string>
#include <vector>

using namespace std;

/**
 * This function checks whether character is alphanumeric or not
 * @param ch character to be validated
 * @return Return true if character is alphanumeric else false
 */
bool is_alnum(char ch) {
    // Check whether character is alphanumeric or not.
    return !isalnum(ch);
}

/**
 * This function validates dirctory name
 * @param dir_name Directory name
 * @return Return true if dir_name has only alphanumeric characters else false
 */
bool validate_dir_path_name(string dir_name) {
    if (dir_name.compare("..") == 0 || dir_name.compare(".") == 0) {
        return true;
    }

    return find_if(dir_name.begin(), dir_name.end(), is_alnum) == dir_name.end();
}

/**
 * This function is used to get directory name from directory path and validates it.
 * @param dir_path Directory path
 * @param index Index is used as starting point in dir_path string to get next dirctory name in path
 * @param separator Separator is used to separate directory name
 * @return Returns directory name
 */
 string get_and_validate_directory_name(string dir_path, int &index, char separator) {
    // Check whether index is greater than or equal to the length of dir_path
    if(index >= dir_path.length()) {
        return dir_path;
    }
    string res = "";

    for(; index < dir_path.length(); index++) {
        if(dir_path[index] == '/') {
            break;
        } else {
            res.push_back(dir_path[index]);
        }
    }
    bool is_alphanum = validate_dir_path_name(res);

    if(is_alphanum) {
        return res;
    } else {
        string expt = ": No such file or directory: ";
        expt += res;
        throw expt;
    }
}

/**
 * This function is used to extract all directory names from directory path
 * @param dir_path Directory path (/abc/edf)
 * @return Returns list of direcrory names corresponding to directory path. 
 */
vector<string> convert_dir_path_to_list(string dir_path) {
    vector<string> res;
    string dir_name = "";

    for(int index = 0; index < dir_path.length(); ) {
        if(dir_path[index] == '/') {
            res.push_back("/");
            index++;
            while(index < dir_path.length() && dir_path[index] == '/') {
                index++;
            }
        } else {
            char separator = '/';
            dir_name = get_and_validate_directory_name(dir_path, index, separator);
            res.push_back(dir_name);
        }
    }
    return res;
}

/**
 * This function is used to prepare new dirctory path
 * @param cwd_dir_path Current working directory
 * @param new_dir_path Destination directory path
 * @return Returns new directory path
 */
vector<string> prepare_new_dir(vector<string> cwd_dir_path, vector<string> new_dir_path) {

	if (new_dir_path[0].compare("/") == 0) {
        cwd_dir_path.clear();
        cwd_dir_path.push_back("/");
	}

	for(int i = 0; i < new_dir_path.size(); i++) {
		if (new_dir_path[i].compare(".") == 0 || new_dir_path[i].compare("/") == 0) {
			continue;
		}
		if (new_dir_path[i].compare("..") == 0) {
			if (cwd_dir_path.size() > 1) {
				cwd_dir_path.pop_back();
                if(cwd_dir_path[cwd_dir_path.size() - 1] == "/") {
                    cwd_dir_path.pop_back();
                }
			}
		} else {
            if (cwd_dir_path[cwd_dir_path.size() - 1] != "/") {
                cwd_dir_path.push_back("/");
            }
            cwd_dir_path.push_back(new_dir_path[i]);
		}
    }
    if(cwd_dir_path.size() == 0) {
        cwd_dir_path.push_back("/");
    }
	return cwd_dir_path;
}

int main(int argc, char **argv) {
    try {
        // Validate number of passed parameters from terminal
        if(argc < 3 || argc > 3) {
            if(argc < 3) {
                cout<<"mycd: too few arguments";
            } else {
                cout<<"mycd: too many arguments";
            }
            return 1;
        }

	string cwd = argv[1];  // Current working directory path
	string to_dirctory = argv[2];  // To new directory path

	vector<string> cwd_dir_path_list = convert_dir_path_to_list(cwd);
	vector<string> new_dir_path_list = convert_dir_path_to_list(to_dirctory);

	vector<string> new_dir = prepare_new_dir(cwd_dir_path_list, new_dir_path_list);
        vector<string>::iterator it;

        for(it = new_dir.begin(); it != new_dir.end(); it++) {
	        cout<<*it;
        }
        cout<<"\n";

    } catch(string exp) { // Default exception
        cout<<"mycd: "<<exp<<"\n";
    }
    return 0;
}
