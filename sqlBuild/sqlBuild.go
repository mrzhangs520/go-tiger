package sqlBuild

import (
	"fmt"
	"github.com/mrzhangs520/go-tiger/dbManager"
	"gorm.io/gorm"
)

type MysqlType struct {
	tx            *gorm.DB
	orderByString string
	groupByFiled  string
	filed         string
	sql           string
	limit         string
	tableName     string
	joinTableList []string
	whereList     []string
}

func (m *MysqlType) OpenTx(tx *gorm.DB) *MysqlType {
	m.tx = tx
	return m
}

func (m *MysqlType) db() *gorm.DB {
	if nil != m.tx {
		return m.tx
	}
	return dbManager.GetInstance().Session(&gorm.Session{})
}

func (m *MysqlType) SetPage(page, pageSize int) *MysqlType {
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

func (m *MysqlType) SetJoinTable(joinTable string) *MysqlType {
	m.joinTableList = append(m.joinTableList, joinTable)
	return m
}

func (m *MysqlType) SetGroupBy(groupByFiled string) *MysqlType {
	m.groupByFiled = groupByFiled
	return m
}

func (m *MysqlType) SetOrderBy(orderByString string) *MysqlType {
	m.orderByString = orderByString
	return m
}

func (m *MysqlType) Where(condition string) {
	m.whereList = append(m.whereList, condition)
}

func (m *MysqlType) Get(data interface{}) int {
	m.build()
	// 定一个临时结构体用于获取总条数
	total := struct{ Total int }{}

	db := m.db()
	db.Raw(m.sql).Scan(data)

	db.Raw("SELECT FOUND_ROWS() as total").Scan(&total)

	return total.Total
}

func (m *MysqlType) build() {
	// 处理join语句
	for _, joinTable := range m.joinTableList {
		m.sql = fmt.Sprintf("%s %s", m.sql, joinTable)
	}

	// 处理where条件
	if 0 < len(m.whereList) {
		m.sql = fmt.Sprintf("%s where", m.sql)
	}
	for index, where := range m.whereList {
		if 0 == index {
			m.sql = fmt.Sprintf("%s %s", m.sql, where)
		} else {
			m.sql = fmt.Sprintf("%s and %s", m.sql, where)
		}
	}

	// 处理groupBy
	if "" != m.groupByFiled {
		m.sql = fmt.Sprintf("%s group by %s", m.sql, m.groupByFiled)
	}

	// 处理orderBy
	if "" != m.orderByString {
		m.sql = fmt.Sprintf("%s order by %s", m.sql, m.orderByString)
	}

	if "" != m.limit {
		// 处理limit
		m.sql = fmt.Sprintf("%s %s", m.sql, m.limit)
	}

	// 处理sql开始部分
	m.sql = fmt.Sprintf("select SQL_CALC_FOUND_ROWS %s from %s a %s", m.filed, m.tableName, m.sql)
}
