package sqlBuild

import (
	"fmt"
	"github.com/mrzhangs520/go-tiger/dbManager"
)

type JoinTableType struct {
	joinTableName string
	condition     string
	join          string
}

type MysqlType struct {
	filed         string
	sql           string
	limit         string
	tableName     string
	joinTableList []JoinTableType
	whereList     []string
}

func (m *MysqlType) SetPage(page, pageSize uint) *MysqlType {
	if 0 == page {
		page = 1
	}
	start := (page - 1) * pageSize
	m.limit = fmt.Sprintf("limit %d, %d", start, pageSize)
	return m
}

func (m *MysqlType) SetTableName(tableName string) *MysqlType {
	m.tableName = tableName
	return m
}

func (m *MysqlType) SetFiled(filed string) *MysqlType {
	m.filed = filed
	return m
}

func (m *MysqlType) SetJoinTable(joinTableName, condition, join string) *MysqlType {
	m.joinTableList = append(m.joinTableList, JoinTableType{joinTableName, condition, join})
	return m
}

func (m *MysqlType) SetCreateTime(startDate, endDate string) *MysqlType {
	condition := fmt.Sprintf("a.create_time BETWEEN '%s' and '%s'", startDate, endDate)
	m.Where(condition)
	return m
}

func (m *MysqlType) Where(condition string) {
	m.whereList = append(m.whereList, condition)
}

func (m *MysqlType) Get(data interface{}) int {
	m.build()
	// 定一个临时结构体用于获取总条数
	total := struct{ Total int }{}
	db := dbManager.GetInstance()
	db.Raw(m.sql).Scan(data)
	db.Raw("SELECT FOUND_ROWS() as total").Scan(&total)
	return total.Total
}

func (m *MysqlType) build() {
	// 构建条件
	if 0 < len(m.whereList) {
		m.sql = fmt.Sprintf("%swhere", m.sql)
	}
	for index, where := range m.whereList {
		if 0 == index {
			m.sql = fmt.Sprintf("%s %s", m.sql, where)
		} else {
			m.sql = fmt.Sprintf("%s and %s", m.sql, where)
		}
	}
	// 处理limit
	m.sql = fmt.Sprintf("%s %s", m.sql, m.limit)

	// 处理sql开始部分
	m.sql = fmt.Sprintf("select SQL_CALC_FOUND_ROWS %s from %s a %s", m.filed, m.tableName, m.sql)
}
