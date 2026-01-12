// cmd/sqlset-gen/main_test.go
package main

import (
	"testing"

	"testing/fstest"

	"github.com/istovpets/sqlset"
	"github.com/stretchr/testify/require"
)

func TestGenerateConstants_Smoke(t *testing.T) {
	testFS := fstest.MapFS{
		"users.sql": &fstest.MapFile{
			Data: []byte(`--SQL: get_user_by_id
SELECT 1;
--end

--SQL: create_user
INSERT...
--end`),
		},
		"posts.sql": &fstest.MapFile{
			Data: []byte(`--SQL: get-post-by-id
SELECT 1;
--end`),
		},
	}

	sqlSet, err := sqlset.New(testFS)
	require.NoError(t, err)

	generated, err := GenerateFileContent(sqlSet, "queries", "test")
	require.NoError(t, err)

	// if err := os.WriteFile("tmp_consts.go", []byte(generated), 0644); err != nil {
	// 	log.Fatal(err)
	// }

	require.Contains(t, generated, `// posts.sql
var Posts = struct {
	GetPostById string
}{
	GetPostById: "posts.get-post-by-id",
}

// users.sql
var Users = struct {
	CreateUser string
	GetUserById string
}{
	CreateUser: "users.create_user",
	GetUserById: "users.get_user_by_id",
}`)
}
