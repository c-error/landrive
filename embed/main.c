#include <windows.h>
#include <stdio.h>
#include <wchar.h>
#include <locale.h>

// #include <time.h>

#define PRESET_START_DIR "<cl><nm><a>FO:</a><b>"
#define PRESET_START_FILE "<cl><nm><a>FI:</a><b>"
#define PRESET_MIDDLE "</b></nm><a>"
#define PRESET_END_DIR "~</a></cl>"
#define PRESET_END_FILE "</a></cl>"
#define PRESET_FOOTER "</shell>"

static const char SPECIAL[] = "!_-+~=@#$%^&(){}[];,'";
#define SPECIAL_SIZE (sizeof(SPECIAL) - 1)

#define SHELL_TOP "<shell><top><div><a>Total: %d</a><a>Folder: %d</a><a>File: %d</a></div><a>%s</a></top>"

typedef struct {
    wchar_t name[MAX_PATH];
    BOOL is_directory;
    int special_priority;
    char size[10];
} FileEntry;



static const char lower_table[256] = {
	['A'] = 'a', ['B'] = 'b', ['C'] = 'c', ['D'] = 'd', ['E'] = 'e',
	['F'] = 'f', ['G'] = 'g', ['H'] = 'h', ['I'] = 'i', ['J'] = 'j',
	['K'] = 'k', ['L'] = 'l', ['M'] = 'm', ['N'] = 'n', ['O'] = 'o',
	['P'] = 'p', ['Q'] = 'q', ['R'] = 'r', ['S'] = 's', ['T'] = 't',
	['U'] = 'u', ['V'] = 'v', ['W'] = 'w', ['X'] = 'x', ['Y'] = 'y',
	['Z'] = 'z'
};

int get_special_priority(wchar_t first_char) {
    // Convert to UTF-8 and check if it's in SPECIAL array
    char utf8_char[5] = {0};
    WideCharToMultiByte(CP_UTF8, 0, &first_char, 1, utf8_char, sizeof(utf8_char), NULL, NULL);
    
    for (int i = 0; i < SPECIAL_SIZE; i++) {
        if (utf8_char[0] == SPECIAL[i]) {
            return i;
        }
    }
    return -1;
}

int natural_wcscmp(const wchar_t *s1, const wchar_t *s2) {
	while (*s1 && *s2) {
		if (iswdigit(*s1) && iswdigit(*s2)) {
			long num1 = wcstol(s1, (wchar_t**)&s1, 10);
			long num2 = wcstol(s2, (wchar_t**)&s2, 10);
			if (num1 != num2) return num1 - num2;
		} else {
			int diff = tolower(*s1) - tolower(*s2);
			if (diff != 0) return diff;
			s1++;
			s2++;
		}
	}
	return *s1 - *s2;
}

int compare_entries(const void* a, const void* b) {
	
	const FileEntry* entry_a = (const FileEntry*)a;
	const FileEntry* entry_b = (const FileEntry*)b;
	
	const BOOL a_is_dir = entry_a->is_directory;
	const int a_priority = entry_a->special_priority;
	const int b_priority = entry_b->special_priority;
	
	if (a_is_dir != entry_b->is_directory) return a_is_dir ? -1 : 1;
	
	if (a_priority >= 0 || b_priority >= 0) {
		if (a_priority >= 0 && b_priority >= 0) return a_priority - b_priority;
		return (a_priority >= 0) ? -1 : 1;
	}
	return natural_wcscmp(entry_a->name, entry_b->name);
}

char* get_file_size(size_t num) {

	char* buffer = malloc(10);
	if (!buffer) return NULL;
	
	if (num >= 0) {
		
		double file_size;
		char* mode;
		
		if (num < 1024) { snprintf(buffer, 10, "%d B", num); return buffer; }
		else if (num >= 1073741824) { mode = "GB"; file_size = num / 1073741824.0; }  // >1024 MB
		else if (num >= 1048576) { mode = "MB"; file_size = num / 1048576.0; }  // >1024 KB
		else if (num >= 1024) { mode = "KB"; file_size = num / 1024.0; }  // >1024 B
		
		if (file_size <= 1024) { snprintf(buffer, 10, "%.2f %s", file_size, mode); return buffer; }
		else  { snprintf(buffer, 10, "x.xx ~"); return buffer; };
	} else  { snprintf(buffer, 10, "x.xx ~"); return buffer; };

}

char* wide_to_utf8(const wchar_t *wide_str) {
    if (!wide_str) return NULL;
    
    int size = WideCharToMultiByte(CP_UTF8, 0, wide_str, -1, NULL, 0, NULL, NULL);
    if (size == 0) return NULL;
    
    char *utf8_str = malloc(size);
    if (!utf8_str) return NULL;
    
    WideCharToMultiByte(CP_UTF8, 0, wide_str, -1, utf8_str, size, NULL, NULL);
    return utf8_str;
}

wchar_t* to_wchar(const char *data) {

    int wlen = MultiByteToWideChar(CP_UTF8, 0, data, -1, NULL, 0);
    if (wlen == 0) return NULL;

    wchar_t* wstr = malloc(wlen * sizeof(wchar_t));
    if (!wstr) return NULL;

    if (MultiByteToWideChar(CP_UTF8, 0, data, -1, wstr, wlen) == 0) {
        free(wstr);
        return NULL;
    }

    return wstr;
}

void char_lower(char *str) {
	if (str == NULL) return;
	
	for (; *str; ++str) {
		*str = lower_table[(unsigned char)*str] ? lower_table[(unsigned char)*str] : *str;
	}
}

void wchar_lower(wchar_t *str) {
	if (str == NULL) return;
	
	for (; *str; ++str) {
		*str = lower_table[(unsigned char)*str] ? lower_table[(unsigned char)*str] : *str;
	}
}

BOOL is_name_match(const wchar_t *find, const wchar_t *name) {

	wchar_t *name_clone = malloc((wcslen(name) + 1) * sizeof(wchar_t));
	if (!name_clone) { free(name_clone); return FALSE; }

	wcscpy(name_clone, name);
	wchar_lower(name_clone);
	
	if (wcsstr(name_clone, find) != NULL) {
		
		free(name_clone);
		return TRUE;
	} else {

		free(name_clone);
		return FALSE;
	}
}

BOOL is_size_match(const char *find, const char *name) {

	char* size_clone = strdup(name);
	if (!size_clone) { free(size_clone); return FALSE; }
						
	char_lower(size_clone);
	if (strstr(size_clone, find) != NULL) {

		free(size_clone);
		return TRUE;
	} else {

		free(size_clone);
		return FALSE;
	}
	
	




	// if (wcsstr(name_clone, find) != NULL) {
		
	// 	free(name_clone);
	// 	return TRUE;
	// } else {

	// 	free(name_clone);
	// 	return FALSE;
	// }
}

char* get_data_list(const WCHAR* directory, int is_filter, int filter_type, const char* filter_name, char* filter_size, const char* clint_ip) {
	
	// char* s = strdup("PRESET_START_DIR");

	// return strdup("PRESET_START_DIR");


	setlocale(LC_ALL, "");
    wchar_t searchPath[MAX_PATH];
    
	// clock_t start_time = clock();

    int folder_count = 0;
	
	size_t NOW_SIZE = 0;
	size_t entry = 0;
    size_t TOTAL_DATA_SIZE = 0;

	

	wchar_t* fiend_name = NULL;
	// char* fiend_size = NULL;

	size_t find_name_len = strlen(filter_name);
	size_t find_size_len = strlen(filter_size);

	if (find_name_len > 0) {
		fiend_name = to_wchar(filter_name);
		if (!fiend_name) return strdup("_E_Unknown filter name");
		wchar_lower(fiend_name);
	}
	if (find_size_len > 0) char_lower(filter_size);

	BOOL FILTER_COMBO = (find_name_len > 0 && find_size_len > 0) ? TRUE : FALSE;

	FileEntry* ent = malloc(8192 * sizeof(FileEntry));
	
    if (wcslen(directory) > 0 && directory[wcslen(directory)-1] == L'\\') {
        swprintf(searchPath, MAX_PATH, L"%ls*", directory);
    } else {
        swprintf(searchPath, MAX_PATH, L"%ls\\*", directory);
    }

	WIN32_FIND_DATAW fileData;
	HANDLE hFind = FindFirstFileExW(
		searchPath,
		FindExInfoBasic,
		&fileData,
		FindExSearchNameMatch,
		NULL,
		FIND_FIRST_EX_LARGE_FETCH
	);
	
	if (hFind == INVALID_HANDLE_VALUE) { FindClose(hFind); return strdup("_E_Forbidden folder path"); }
	do {
		
        if (fileData.cFileName[0] == L'.' && 
            (fileData.cFileName[1] == L'\0' || 
             (fileData.cFileName[1] == L'.' && fileData.cFileName[2] == L'\0'))) {
            continue;
        }
		
		LARGE_INTEGER fileSize;
		
		fileSize.LowPart = fileData.nFileSizeLow;
		fileSize.HighPart = fileData.nFileSizeHigh;
		
		if (entry >= NOW_SIZE) {
            NOW_SIZE += 8192;
			// size_t new_cap = NOW_SIZE + 8192;
			FileEntry* new_ent = realloc(ent, NOW_SIZE * sizeof(FileEntry));
			
			if (!new_ent) { FindClose(hFind); return strdup("_E_Internal server memory error"); }
			ent = new_ent;
		}
		
// printf("PRE: %s\n", fileSize.QuadPart);

		const char* data_size = get_file_size(fileSize.QuadPart);
		BOOL data_type = fileData.dwFileAttributes & FILE_ATTRIBUTE_DIRECTORY;
		
		if (is_filter == 1) {

			if (filter_type == 2 && data_type) goto JUMP_POINT;
			else if (filter_type == 1 && !data_type) goto JUMP_POINT;

			BOOL UNMATCH_NAME = (find_name_len > 0 && !is_name_match(fiend_name, fileData.cFileName)) ? TRUE : FALSE;
			BOOL UNMATCH_SIZE = (find_size_len > 0 && !is_size_match(filter_size, data_size)) ? TRUE : FALSE;

			if (UNMATCH_NAME || UNMATCH_SIZE) goto JUMP_POINT;
		}

        TOTAL_DATA_SIZE += wcslen(fileData.cFileName) * 4 + 10;
        
		wcsncpy(ent[entry].name, fileData.cFileName, MAX_PATH);
		ent[entry].name[MAX_PATH-1] = '\0';
		ent[entry].is_directory = data_type;
		ent[entry].special_priority = (fileData.cFileName[0] != '\0') ? get_special_priority(fileData.cFileName[0]) : -1;
		
		strncpy(ent[entry].size, data_size, 10);
		ent[entry].size[9] = '\0';

		if (data_type) folder_count++;
		entry++;
		
		JUMP_POINT:
// printf("END: %s\n", data_size);
		free((void*)data_size);
		
	} while (FindNextFileW(hFind, &fileData));
	
	if(find_name_len > 0) free(fiend_name);
	FindClose(hFind);

	qsort(ent, entry, sizeof(FileEntry), compare_entries);
	
    int file_count = entry - folder_count;

	size_t total_allocated = TOTAL_DATA_SIZE + (entry * (60+42+8)) + snprintf(NULL, 0, "%d%d%d%s", entry, folder_count, file_count, clint_ip) + 86 + 10;

	char* out = malloc(total_allocated);
	if (!out) return strdup("_E_Internal server buffer allocation error");

	char* current_pos = out;

	current_pos += sprintf(current_pos, SHELL_TOP, entry, folder_count, file_count, clint_ip);

    for (size_t i = 0; i < entry; i++) {
        char *utf8_name = wide_to_utf8(ent[i].name);
        if (!utf8_name) continue;

        if (ent[i].is_directory) {
            memcpy(current_pos, PRESET_START_DIR, 21);
            current_pos += 21;
            strcpy(current_pos, utf8_name);
            current_pos += strlen(utf8_name);
            memcpy(current_pos, PRESET_MIDDLE, 12);
            current_pos += 12;
            memcpy(current_pos, PRESET_END_DIR, 10);
            current_pos += 10;
        } else {
            memcpy(current_pos, PRESET_START_FILE, 21);
            current_pos += 21;
            strcpy(current_pos, utf8_name);
            current_pos += strlen(utf8_name);
            memcpy(current_pos, PRESET_MIDDLE, 12);
            current_pos += 12;
            strcpy(current_pos, ent[i].size);
            current_pos += strlen(ent[i].size);
            memcpy(current_pos, PRESET_END_FILE, 9);
            current_pos += 9;
        }
        free(utf8_name);
    }

	memcpy(current_pos, PRESET_FOOTER, 8);
	current_pos += 8;
	*current_pos = '\0';
	
	// printf("%s", out);
	// printf("\n DIR: %s ...\n\tTOTAL: %d | FOLDER: %d | FILE: %d\n", searchPath, entry, folder_count, (entry - folder_count));
	
	size_t bytes_used = current_pos - out;
	if (bytes_used > total_allocated) return strdup("_E_Internal server buffer overflow error");

	free(ent);
	return out;
}
