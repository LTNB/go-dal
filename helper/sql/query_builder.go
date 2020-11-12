package sql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type SelectQueryBuilder struct {
	SelectFields []string
	From         []string
	WhereClause  WhereClauseBuilder
	JoinCondition JoinCondition
	OrderBy      map[string]string
	Limit        int
	Offset       int
}

type WhereClauseBuilder struct {
	Pair              map[string]interface{}
	NativeWhereClause string
}

type JoinCondition struct {
	Pair              map[string]string
	NativeClause string
}

type InsertBuilder struct {
	TableName string
	Keys      []string
	Values    [][]interface{}
}

type UpdateBuilder struct {
	TableName   string
	WhereClause WhereClauseBuilder
	JoinCondition JoinCondition
	Values      map[string]interface{}
}

type DeleteBuilder struct {
	TableName   string
	WhereClause WhereClauseBuilder
	JoinCondition JoinCondition
}

func (builder SelectQueryBuilder) BuildSelectQuery() (string, error) {
	if len(builder.From) == 0 {
		return "", errors.New("tables in FROM field not found")
	}
	result := "SELECT %v "
	if len(builder.SelectFields) == 0 {
		builder.SelectFields = append(builder.SelectFields, "*")
	}
	result = fmt.Sprintf(result, strings.Join(builder.SelectFields, ","))
	result = fmt.Sprintf(result+"FROM %v ", strings.Join(builder.From, ","))
	if len(builder.WhereClause.Pair) != 0 || builder.WhereClause.NativeWhereClause != "" {
		result = builder.WhereClause.buildWhereClause(result)
	}
	if len(builder.JoinCondition.Pair) != 0 || builder.JoinCondition.NativeClause != "" {
		result = builder.JoinCondition.buildJoinCondition(result)
	}

	if len(builder.OrderBy) != 0 {
		result = builder.buildOrderBy(result)
	}

	if builder.Limit != 0 {
		result = builder.buildLimit(result)
	}

	if builder.Offset != -1 {
		result = builder.buildOffset(result)
	}

	return result, nil
}

func (builder WhereClauseBuilder) buildWhereClause(result string) string {
	result = result + "WHERE "
	if len(builder.Pair) != 0 {
		sliceWhereClause := make([]string, 0)
		for k, v := range builder.Pair {
			if reflect.TypeOf(v).Kind() == reflect.String {
				v = fmt.Sprintf("'%v'", v)
			}
			sliceWhereClause = append(sliceWhereClause, fmt.Sprintf("%v = %v", k, v))
		}
		result = fmt.Sprintf(result+"%v", strings.Join(sliceWhereClause, " AND "))
	}
	if len(builder.NativeWhereClause) != 0 {
		if len(builder.Pair) != 0 {
			result = result + " AND %v"
		} else {
			result = result + " %v"
		}
		result = fmt.Sprintf(result, builder.NativeWhereClause)
	}
	return result
}

func (builder JoinCondition) buildJoinCondition(result string) string {
	result = result + " AND "
	if len(builder.Pair) != 0 {
		sliceWhereClause := make([]string, 0)
		for k, v := range builder.Pair {
			sliceWhereClause = append(sliceWhereClause, fmt.Sprintf("%v = %v", k, v))
		}
		result = fmt.Sprintf(result+"%v", strings.Join(sliceWhereClause, " AND "))
	}
	if len(builder.NativeClause) != 0 {
		if len(builder.Pair) != 0 {
			result = result + " AND %v"
		} else {
			result = result + " %v"
		}
		result = fmt.Sprintf(result, builder.NativeClause)
	}
	return result
}

func (builder SelectQueryBuilder) buildOrderBy(result string) string {
	result = result + " ORDER BY "
	sliceOrderBy := make([]string, 0)
	for k, v := range builder.OrderBy {
		sliceOrderBy = append(sliceOrderBy, fmt.Sprintf("%v %v", k, v))
	}
	result = fmt.Sprintf(result+"%v", strings.Join(sliceOrderBy, ","))
	return result
}

/**
 * if limit == 0 ==> return all records
 */
func (builder SelectQueryBuilder) buildLimit(result string) string {
	if builder.Limit == 0 {
		return result
	} else {
		return fmt.Sprintf(result+" LIMIT %v", builder.Limit)
	}
}

func (builder SelectQueryBuilder) buildOffset(result string) string {
	return fmt.Sprintf(result+" OFFSET %v", builder.Offset)
}

//=======insert========

func (insertBuilder InsertBuilder) BuildInsertQuery() (string, error) {
	if insertBuilder.TableName == "" || len(insertBuilder.Keys) == 0 || len(insertBuilder.Values) == 0 {
		return "", errors.New("no enough data for build query")
	}
	result := fmt.Sprintf("INSERT INTO %v ", insertBuilder.TableName)             //INSERT INTO table_name
	result = fmt.Sprintf(result+"( %v ) ", strings.Join(insertBuilder.Keys, ",")) // INSERT INTO table_name (...)
	values := make([]string, 0)
	for _, val := range insertBuilder.Values {
		ret := make([]string, len(val))
		for i := 0; i < len(val); i++ {
			if reflect.TypeOf(val[i]).Kind() == reflect.String {
				ret[i] = fmt.Sprintf("'%v'", val[i])
			} else {
				ret[i] = fmt.Sprintf("%v", val[i])
			}
		}
		values = append(values, fmt.Sprintf("(%v)", strings.Join(ret, ",")))
	}

	result = result + fmt.Sprintf(" VALUES %v", strings.Join(values, ",")) //INSERT INTO table_name (...) values (...), (...)
	return result, nil
}

//========update builder
func (builder UpdateBuilder) BuildUpdateQuery() (string, error) {
	if builder.TableName == "" || (builder.WhereClause.NativeWhereClause == "" && builder.WhereClause.Pair == nil) || len(builder.Values) == 0 {
		return "", errors.New("no enough data for build query")
	}
	result := fmt.Sprintf("UPDATE %v ", builder.TableName)
	sliceValue := make([]string, 0)
	for k, v := range builder.Values {
		if reflect.TypeOf(v).Kind() == reflect.String {
			sliceValue = append(sliceValue, fmt.Sprintf("%v = '%v'", k, v))
		} else {
			sliceValue = append(sliceValue, fmt.Sprintf("%v = %v", k, v))
		}
	}
	result = fmt.Sprintf(result+"SET %v ", strings.Join(sliceValue, ","))
	result = builder.WhereClause.buildWhereClause(result)

	if len(builder.JoinCondition.Pair) != 0 || builder.JoinCondition.NativeClause != "" {
		result = builder.JoinCondition.buildJoinCondition(result)
	}

	return result, nil
}

//======== delete builder
func (builder DeleteBuilder) BuildDeleteQuery() (string, error) {
	if builder.TableName == "" || (builder.WhereClause.NativeWhereClause == "" && builder.WhereClause.Pair == nil) {
		return "", errors.New("no enough data for build query")
	}
	result := fmt.Sprintf("DELETE FROM %v ", builder.TableName)
	result = builder.WhereClause.buildWhereClause(result)
	if len(builder.JoinCondition.Pair) != 0 || builder.JoinCondition.NativeClause != "" {
		result = builder.JoinCondition.buildJoinCondition(result)
	}
	return result, nil
}
