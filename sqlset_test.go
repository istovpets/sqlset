package sqlset_test

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/istovpets/sqlset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/valid_multi/*.sql
var testdataValidMulti embed.FS

//go:embed testdata/valid_single/*.sql
var testdataValidSingle embed.FS

//go:embed testdata/invalid/meta1.sql
var testdataInvalidMeta1 embed.FS

//go:embed testdata/invalid/meta2.sql
var testdataInvalidMeta2 embed.FS

//go:embed testdata/invalid/syntax1.sql
var testdataInvalidSyntax1 embed.FS

//go:embed testdata/invalid/syntax2.sql
var testdataInvalidSyntax2 embed.FS

//go:embed testdata/invalid/long-lines.sql
var testdataInvalidLongLines embed.FS

//nolint:funlen,lll
func TestSQLSet(t *testing.T) {
	sqlSet, err := sqlset.New(testdataValidMulti)
	require.NoError(t, err)
	require.NotNil(t, sqlSet)

	var sets sqlset.SQLSetsProvider = sqlSet
	var queries sqlset.SQLQueriesProvider = sqlSet

	queryTests := []struct {
		setID         string
		queryID       string
		expectedQuery string
		expectedErr   error
	}{
		{
			setID:         "test-id-override-1",
			queryID:       "GetData1",
			expectedQuery: "SELECT '515bbf3c-93c5-476a-8dbc-4a6db4fe3c0c' AS id, 'Igor' AS name, 'en' AS language, 'igor@example.com' AS email, ARRAY['token1','token2'] AS tokens;",
			expectedErr:   nil,
		},
		{
			setID:         "test-id-override-1",
			queryID:       "GetData2",
			expectedQuery: "SELECT 'ef84af8f-bb55-4f74-9d7c-3db30e740d20' AS id, 'Alexey' AS name, 'en' AS language, 'alex@example.com' AS email, '{}'::varchar[] as tokens;",
			expectedErr:   nil,
		},
		{
			setID:         "test-id-override-1",
			queryID:       "GetData3",
			expectedQuery: "SELECT 'e192f9e5-5e5c-4bba-b13e-0f9de32ec6bd' AS id, 'Denis' AS name, 'en' AS language, 'denis@example.com' AS email, ARRAY['token3','token4'] AS tokens;",
			expectedErr:   nil,
		},
		{
			setID:       "test-id-override-1",
			queryID:     "unknown",
			expectedErr: sqlset.ErrNotFound,
		},
	}

	for _, test := range queryTests {
		t.Run("Get "+test.setID+":"+test.queryID, func(t *testing.T) {
			t.Parallel()

			query, err := queries.Get(test.setID, test.queryID)

			if test.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, test.expectedErr)
			}

			assert.Equal(t, test.expectedQuery, query)
		})

		t.Run("MustGet "+test.setID+":"+test.queryID, func(t *testing.T) {
			t.Parallel()

			var query string

			fn := func() {
				query = queries.MustGet(test.setID, test.queryID)
			}

			if test.expectedErr == nil {
				assert.NotPanics(t, fn)
			} else {
				assert.Panics(t, fn)
			}

			assert.Equal(t, test.expectedQuery, query)
		})
	}

	t.Run("GetSetsMetas", func(t *testing.T) {
		t.Parallel()

		metas := sets.GetSetsMetas()

		require.Len(t, metas, 2)
		assert.Contains(t, metas, sqlset.QuerySetMeta{
			ID:          "test-id-override-1",
			Name:        "Test 1",
			Description: "Test description 1",
		})
		assert.Contains(t, metas, sqlset.QuerySetMeta{
			ID:          "test2",
			Name:        "test2",
			Description: "Test description 2",
		})
	})

	t.Run("GetQueryIDs", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name          string
			setID         string
			expectedIDs   []string
			expectedErr   error
			expectSuccess bool
		}{
			{
				name:          "get query IDs for test-id-override-1",
				setID:         "test-id-override-1",
				expectedIDs:   []string{"GetData1", "GetData2", "GetData3"},
				expectSuccess: true,
			},
			{
				name:          "get query IDs for test2",
				setID:         "test2",
				expectedIDs:   []string{"query1", "query2"},
				expectSuccess: true,
			},
			{
				name:          "get query IDs for non-existent set",
				setID:         "nonexistent",
				expectedErr:   sqlset.ErrNotFound,
				expectSuccess: false,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ids, err := sets.GetQueryIDs(test.setID)

				if test.expectSuccess {
					require.NoError(t, err)
					assert.Equal(t, test.expectedIDs, ids)
				} else {
					require.ErrorIs(t, err, test.expectedErr)
					assert.Nil(t, ids)
				}
			})
		}
	})
}

func TestSQLSet_Get_SingleArgument(t *testing.T) {
	t.Parallel()

	sqlSetValid, err := sqlset.New(testdataValidMulti)
	require.NoError(t, err)

	sqlSetValidSingle, err := sqlset.New(testdataValidSingle)
	require.NoError(t, err)

	tests := []struct {
		name          string
		sqlSet        sqlset.SQLQueriesProvider
		args          []string
		expectedQuery string
		expectedErr   error
	}{
		{
			name:   "get with setID.queryID",
			sqlSet: sqlSetValid,
			args:   []string{"test-id-override-1.GetData1"},
			expectedQuery: "SELECT '515bbf3c-93c5-476a-8dbc-4a6db4fe3c0c' AS id, 'Igor' AS name, 'en' AS language, " +
				"'igor@example.com' AS email, ARRAY['token1','token2'] AS tokens;",
		},
		{
			name:        "get with unknown setID in setID.queryID",
			sqlSet:      sqlSetValid,
			args:        []string{"unknown.GetData1"},
			expectedErr: sqlset.ErrNotFound,
		},
		{
			name:        "get with unknown queryID in setID.queryID",
			sqlSet:      sqlSetValid,
			args:        []string{"test-id-override-1.unknown"},
			expectedErr: sqlset.ErrNotFound,
		},
		{
			name:   "get with queryID from single set",
			sqlSet: sqlSetValidSingle,
			args:   []string{"GetData1"},
			expectedQuery: "SELECT '515bbf3c-93c5-476a-8dbc-4a6db4fe3c0c' AS id, 'Igor' AS name, 'en' AS language, " +
				"'igor@example.com' AS email, ARRAY['token1','token2'] AS tokens;",
			expectedErr: nil,
		},
		{
			name:        "get with queryID from multiple sets",
			sqlSet:      sqlSetValid,
			args:        []string{"GetData1"},
			expectedErr: sqlset.ErrRequiredArgMissing,
		},
		{
			name:        "get with no arguments",
			sqlSet:      sqlSetValid,
			args:        []string{},
			expectedErr: sqlset.ErrInvalidArgCount,
		},
		{
			name:        "get with too many arguments",
			sqlSet:      sqlSetValid,
			args:        []string{"a", "b", "c"},
			expectedErr: sqlset.ErrInvalidArgCount,
		},
		{
			name:        "get with empty argument",
			sqlSet:      sqlSetValid,
			args:        []string{""},
			expectedErr: sqlset.ErrArgumentEmpty,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			query, err := test.sqlSet.Get(test.args...)

			if test.expectedErr != nil {
				require.ErrorIs(t, err, test.expectedErr)
				assert.Empty(t, query)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedQuery, query)
			}
		})
	}
}

func TestNew_WhenInvalid_ExpectError(t *testing.T) {
	tests := []struct {
		name        string
		fs          fs.FS
		expectedErr error
	}{
		{
			name:        "invalid meta 1",
			fs:          testdataInvalidMeta1,
			expectedErr: sqlset.ErrInvalidSyntax,
		},
		{
			name:        "invalid meta 2",
			fs:          testdataInvalidMeta2,
			expectedErr: sqlset.ErrInvalidSyntax,
		},
		{
			name:        "invalid syntax 1",
			fs:          testdataInvalidSyntax1,
			expectedErr: sqlset.ErrInvalidSyntax,
		},
		{
			name:        "invalid syntax 2",
			fs:          testdataInvalidSyntax2,
			expectedErr: sqlset.ErrInvalidSyntax,
		},
		{
			name:        "long lines",
			fs:          testdataInvalidLongLines,
			expectedErr: sqlset.ErrMaxLineLenExceeded,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			set, err := sqlset.New(test.fs)

			//nolint:testifylint
			assert.ErrorIs(t, err, test.expectedErr)
			assert.Nil(t, set)
		})
	}
}
