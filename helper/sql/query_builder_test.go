package sql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

/**
 * sometimes the where-clause will be exchanged index element it's, so test case will be failed, but it not a problem
 *
 */

func TestQueryBuilder(t *testing.T) {
	whereMap := make(map[string]interface{})
	whereMap["ac.id"] = "re.id"
	whereMap["name"] = "abc"
	orderBy := make(map[string]string)
	orderBy["id"] = ""
	orderBy["name"] = "DESC"
	builder := SelectQueryBuilder{
		SelectFields: []string{"id", "name", "age"},
		From:         []string{"account ac", "receive re"},
		Limit:        1,
		Offset:       10,
		OrderBy:      orderBy,
		WhereClause: WhereClauseBuilder{
			Pair:              whereMap,
			NativeWhereClause: "name = 'abc'",
		},
	}
	sql, _ := builder.BuildSelectQuery();
	assert.Equal(t, "SELECT id,name,age FROM account ac,receive re WHERE ac.id = re.id AND name = 'abc' AND name = 'abc' ORDER BY id ,name DESC LIMIT 1 OFFSET 10", sql, "done")
}

func TestInsertBuilder(t *testing.T) {
	val1 := make([]interface{}, 3)
	val1[0] = 1
	val1[1] = "name1"
	val1[2] = "name2"

	val2 := make([]interface{}, 3)
	val2[0] = 2
	val2[1] = "name2"
	val2[2] = "name2"
	values := make([][]interface{}, 2)
	values[0] = val1
	values[1] = val2
	builder := InsertBuilder{
		TableName: "account",
		Keys:      []string{"id", "name", "age"},
		Values:    values,
	}
	sql, _:= builder.BuildInsertQuery()
	assert.Equal(t, "INSERT INTO account ( id,name,age )  VALUES (1,'name1','name2'),(2,'name2','name2')", sql, "done")
}

func TestUpdateBuilder(t *testing.T) {
	whereMap := make(map[string]interface{})
	whereMap["id"] = 1
	whereMap["name"] = "abc"
	setMap := make(map[string]interface{})
	setMap["name"] = "xyz"
	setMap["age"] = 18
	builder := UpdateBuilder{
		TableName: "account",
		WhereClause: WhereClauseBuilder{
			Pair:              whereMap,
			NativeWhereClause: "",
		},
		Values: setMap,
	}

	sql, _ := builder.BuildUpdateQuery()
	assert.Equal(t, "UPDATE account SET name = 'xyz',age = 18 WHERE id = 1 AND name = 'abc'", sql, "done")
}

func TestDeleteBuilder(t *testing.T) {
	whereMap := make(map[string]interface{})
	whereMap["id"] = 1
	builder := DeleteBuilder{
		TableName: "account",
		WhereClause: WhereClauseBuilder{
			Pair:              whereMap,
			NativeWhereClause: "",
		},
	}
	sql, _ := builder.BuildDeleteQuery()
	assert.Equal(t, "DELETE FROM account WHERE id = 1", sql, "done")
}

type MockBo struct {
	Id int `json:"id" primary:"id"`
	Token string `json:"token" primary:"token"`
	Name string
}

func TestGetPrimaryKeysValues(t *testing.T){
	bo := MockBo{
		Id:    1,
		Token: "abc",
		Name: "name",
	}

	result, err := getPrimaryKeysValues(bo, map[string]interface{}{})
	fmt.Println(result)
	assert.Nil(t, err, "done")
}

//func TestComparePerformance(t *testing.T){
//	s := time.Now().Nanosecond()
//	fmt.Println("SELECT id,name,age FROM account ac,receive re WHERE name = 'abc' AND ac.id = re.id AND name = 'abc' ORDER BY id ,name DESC LIMIT 1 OFFSET 10")
//	fmt.Printf("builder duration: %v", (time.Now().Nanosecond()- s))
//}

func TestMain(m *testing.M) {
	r := m.Run()
	//destroy()
	os.Exit(r)
}
