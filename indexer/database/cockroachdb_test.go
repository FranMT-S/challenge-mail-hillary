package database

import (
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
	TestCreateTable(t)
	TestIsTableCreated(t)
	TestSendMails(t)

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

func TestIsTableCreated(t *testing.T) {
	t.Run("Must be return true if the table exists", func(t *testing.T) {
		conn, err := getConn()
		assert.NoError(t, err)

		_, err = conn.Open()
		assert.NoError(t, err)
		defer conn.Close()

		exists, err := conn.IsTableCreated(DBTableNameTest)
		assert.NoError(t, err)
		assert.True(t, exists)

	})
}

func TestCreateTable(t *testing.T) {
	conn, err := getConn()
	assert.NoError(t, err)

	_, err = conn.Open()
	assert.NoError(t, err)

	defer conn.Close()

	// if exists drop it
	assert.NoError(t, conn.DropTable(DBTableNameTest))
	assert.NoError(t, conn.CreateTableIfNotExists(DBTableNameTest))
	exists, err := conn.IsTableCreated(DBTableNameTest)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestSendMails(t *testing.T) {
	conn, err := getConn()
	assert.NoError(t, err)

	_, err = conn.Open()
	assert.NoError(t, err)

	_, err = conn.DB.Exec("DELETE FROM " + DBTableNameTest + ";")
	assert.NoError(t, err)
	defer conn.Close()
	defer conn.DB.Exec("DELETE FROM " + DBTableNameTest + ";")

	// send empty batch
	t.Run("Must be able to send an empty batch", func(t *testing.T) {
		rows, err := conn.SendMails(DBTableNameTest, []models.Email{})
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
		rows, err := conn.SendMails(DBTableNameTest, testMails)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), rows)
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

	t.Run("Must be 4 rows in the table", func(t *testing.T) {
		var count int64
		rows, err := conn.SendMails(DBTableNameTest, testMails)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), rows)

		err = conn.DB.QueryRow("SELECT COUNT(*) FROM " + DBTableNameTest).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, int64(4), count)
	})

}
