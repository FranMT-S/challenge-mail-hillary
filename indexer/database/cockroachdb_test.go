package database

import (
	"fmt"
	"testing"
	"time"

	"indexer/models"

	"github.com/stretchr/testify/assert"
)

func getConn() (*Connection, error) {
	return NewConnection(
		"localhost", // host
		"defaultdb", // default DB
		"root",      // user
		"",          // password (empty for --insecure)
		"26257",     // port
		false,       // SSL disabled with --insecure
	)

}

// IntegrationTestFlow runs all the integration tests
func TestIntegration_Flow(t *testing.T) {

	TestConnection(t)
	TestCreateSchema(t)
	TestIsSchemaCreated(t)
	TestSendMails(t)
	conn, _ := getConn()
	_, _ = conn.Open()
	defer conn.Close()
	defer conn.DB.Exec(fmt.Sprintf(`DROP SCHEMA IF EXISTS %s CASCADE;`, DBSchemaNameTest))
}

// TestConnection tests the connection to the database
func TestConnection(t *testing.T) {

	t.Run("Test connection", func(t *testing.T) {
		conn, err := getConn()

		assert.NoError(t, err)

		_, err = conn.Open()
		assert.NoError(t, err)
		defer conn.Close()

		assert.NoError(t, conn.Ping())
	})
}

func TestCreateSchema(t *testing.T) {
	conn, err := getConn()
	assert.NoError(t, err)

	_, err = conn.Open()
	assert.NoError(t, err)

	defer conn.Close()

	// if exists drop it
	assert.NoError(t, conn.CreateSchemaIfNotExist("test_query"))
	exists, err := conn.IsSchemaCreated("test_query")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestIsSchemaCreated(t *testing.T) {
	t.Run("Must be return true if the table exists", func(t *testing.T) {
		conn, err := getConn()
		assert.NoError(t, err)

		_, err = conn.Open()
		assert.NoError(t, err)
		defer conn.Close()

		exists, err := conn.IsSchemaCreated(DBSchemaNameTest)
		assert.NoError(t, err)
		assert.True(t, exists)

	})
}

func TestSendMails(t *testing.T) {
	conn, err := getConn()
	assert.NoError(t, err)

	_, err = conn.Open()
	assert.NoError(t, err)

	truncateQuery := fmt.Sprintf(`
	delete from %s.emails_search;
	delete from %s.emails; 
	`, DBSchemaNameTest, DBSchemaNameTest)
	_, err = conn.DB.Exec(truncateQuery)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	defer conn.Close()

	defer conn.DB.Exec(truncateQuery)

	t.Run("Must be able to send an empty batch", func(t *testing.T) {
		rows, err := conn.SendMails(DBSchemaNameTest, []models.Email{})
		assert.NoError(t, err)
		assert.Equal(t, int64(0), rows)
	})

	testMails := []models.Email{
		{
			ID:      1,
			Date:    time.Now().UTC(),
			Subject: "this is a testing email",
			From:    "Hillary Clinton <hillary@clinton.com>",
			To:      "Bill Clinton <bill@clinton.com>",
			Content: "this is a testing email content",
		},
	}

	t.Run("Must be 1 row in the table", func(t *testing.T) {
		rows, err := conn.SendMails(DBSchemaNameTest, testMails)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), rows)
		var total int64
		conn.DB.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM %s.emails;`, DBSchemaNameTest)).Scan(&total)
		assert.Equal(t, int64(1), total)
		var searchVectorTotal int64
		conn.DB.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM %s.emails_search;`, DBSchemaNameTest)).Scan(&searchVectorTotal)
		assert.Equal(t, int64(1), searchVectorTotal)
	})

	testMails = []models.Email{
		{
			ID:      2,
			Date:    time.Now().UTC(),
			Subject: "this is a secret message",
			From:    "Hillary Clinton <hillary@clinton.com>",
			To:      "Bill Clinton <bill@clinton.com>",
			Content: "this is a secret message content",
		},
		{
			ID:      3,
			Date:    time.Now().UTC(),
			Subject: "this is a private message",
			From:    "Hillary Clinton <hillary@clinton.com>",
			To:      "Bill Clinton <bill@clinton.com>",
			Content: "this is a private message content",
		},
		{
			ID:      3,
			Date:    time.Now().UTC(),
			Subject: "this is a private message",
			From:    "Hillary Clinton <hillary@clinton.com>",
			To:      "Bill Clinton <bill@clinton.com>",
			Content: "this is a private message content",
		},
		{
			ID:      4,
			Date:    time.Now().UTC(),
			Subject: "this is a confidential message",
			From:    "Hillary Clinton <hillary@clinton.com>",
			To:      "Bill Clinton <bill@clinton.com>",
			Content: "this is a confidential message content",
		},
	}

	t.Run("Must be 4 rows in the table and insert 3 rows", func(t *testing.T) {
		var count int64
		rowsInserted, err := conn.SendMails(DBSchemaNameTest, testMails)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), rowsInserted)

		err = conn.DB.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM %s.emails;`, DBSchemaNameTest)).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, int64(4), count)
		var searchVectorTotal int64
		conn.DB.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM %s.emails_search;`, DBSchemaNameTest)).Scan(&searchVectorTotal)
		assert.Equal(t, int64(4), searchVectorTotal)
	})

}
