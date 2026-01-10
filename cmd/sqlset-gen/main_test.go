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
			Data: []byte(`--SQL: GetUserById
SELECT 1;
--end

--SQL: CreateUser
INSERT...
--end`),
		},
		"posts.sql": &fstest.MapFile{
			Data: []byte(`--SQL: GetPostById
SELECT 1;
--end`),
		},
	}

	sqlSet, err := sqlset.New(testFS)
	require.NoError(t, err)

	generated, err := GenerateConstants(sqlSet, "queries")
	require.NoError(t, err)

	// минимальные проверки
	require.Contains(t, generated, `UsersGetUserById = "users.GetUserById"`)
	require.Contains(t, generated, `UsersCreateUser = "users.CreateUser"`)
	require.Contains(t, generated, `PostsGetPostById = "posts.GetPostById"`)

	// if err := os.WriteFile("tmp_consts.go", []byte(generated), 0644); err != nil {
	// 	log.Fatal(err)
	// }
}
