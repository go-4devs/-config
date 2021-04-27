package ini_test

import (
	"testing"
	"time"

	"gitoa.ru/go-4devs/config/provider/ini"
	"gitoa.ru/go-4devs/config/test"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	file := test.NewINI()

	read := []test.Read{
		test.NewRead("project/PROJECT_BOARD_BASIC_KANBAN_TYPE", "To Do, In Progress, Done"),
		test.NewRead("repository.editor/PREVIEWABLE_FILE_MODES", "markdown"),
		test.NewRead("server/LOCAL_ROOT_URL", "http://0.0.0.0:3000/"),
		test.NewRead("server/LFS_HTTP_AUTH_EXPIRY", 20*time.Minute),
		test.NewRead("repository.pull-request/DEFAULT_MERGE_MESSAGE_SIZE", 5120),
		test.NewRead("ui/SHOW_USER_EMAIL", true),
		test.NewRead("cors/ENABLED", false),
	}

	prov := ini.New(file)
	test.Run(t, prov, read)
}
