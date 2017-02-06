/*
 * go-mysqlstack
 * xelabs.org
 *
 * Copyright (c) XeLabs
 * GPL License
 *
 */

package driver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	querypb "github.com/XeLabs/go-mysqlstack/sqlparser/depends/query"
	"github.com/XeLabs/go-mysqlstack/sqlparser/depends/sqltypes"
	"github.com/XeLabs/go-mysqlstack/xlog"
)

func TestRows(t *testing.T) {
	result1 := &sqltypes.Result{
		Fields: []*querypb.Field{
			{
				Name: "id",
				Type: querypb.Type_INT32,
			},
			{
				Name: "name",
				Type: querypb.Type_VARCHAR,
			},
		},
		Rows: [][]sqltypes.Value{
			{
				sqltypes.MakeTrusted(querypb.Type_INT32, []byte("10")),
				sqltypes.MakeTrusted(querypb.Type_VARCHAR, []byte("nice name")),
			},
			{
				sqltypes.MakeTrusted(querypb.Type_INT32, []byte("20")),
				sqltypes.NULL,
			},
		},
	}
	result2 := &sqltypes.Result{
		RowsAffected: 123,
		InsertID:     123456789,
	}

	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	th := NewTestHandler(log)
	port, svr, err := MockMysqlServer(log, th)
	assert.Nil(t, err)
	defer svr.Close()
	address := fmt.Sprintf(":%v", port)

	// query
	{
		client, err := NewConn("mock", "mock", address, "test")
		assert.Nil(t, err)
		defer client.Close()

		th.SetCond(&Cond{Query: "SELECT2", Result: result2})
		rows, err := client.Query("SELECT2")
		assert.Nil(t, err)

		assert.Equal(t, uint64(123), rows.RowsAffected())
		assert.Equal(t, uint64(123456789), rows.LastInsertID())
	}

	// query
	{
		client, err := NewConn("mock", "mock", address, "test")
		assert.Nil(t, err)
		defer client.Close()

		th.SetCond(&Cond{Query: "SELECT1", Result: result1})
		rows, err := client.Query("SELECT1")
		assert.Nil(t, err)
		assert.Equal(t, result1.Fields, rows.Fields())
		for rows.Next() {
			_ = rows.Datas()
			_, _ = rows.RowValues()
		}
	}
}
