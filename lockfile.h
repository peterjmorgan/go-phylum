#include <stddef.h>

typedef void const *lockfile_format;
typedef void *lockfile_parse_result;

typedef struct {
    const char* head;
    const size_t len;
} lockfile_string;

typedef struct {
    lockfile_string name;
    lockfile_string version;
} lockfile_dependency;

extern lockfile_format lockfile_format_from_str(const char* name);

extern lockfile_format lockfile_format_for_path(const char* path);

extern const char* lockfile_format_get_name(lockfile_format format);

extern int lockfile_format_is_path_lockfile(lockfile_format format, const char* path);

extern lockfile_parse_result lockfile_format_parse(lockfile_format format, const char* content, size_t content_length);

extern int lockfile_parse_result_is_ok(lockfile_parse_result result);

extern void lockfile_parse_result_get_error(lockfile_parse_result result, lockfile_string* error);

extern size_t lockfile_parse_result_get_dependencies_len(lockfile_parse_result result);

extern void lockfile_parse_result_get_dependency(lockfile_parse_result result, size_t index, lockfile_dependency* dependency);

extern void lockfile_parse_result_destroy(lockfile_parse_result result);
