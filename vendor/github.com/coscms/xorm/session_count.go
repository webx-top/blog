package xorm

// Count counts the records. bean's non-empty fields
// are conditions.
func (session *Session) Count(bean interface{}) (int64, error) {
	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	var sqlStr string
	var args []interface{}
	if session.Statement.RawSQL == "" {
		sqlStr, args = session.Statement.genCountSQL(bean)
	} else {
		sqlStr = session.Statement.RawSQL
		args = session.Statement.RawParams
	}

	session.queryPreprocess(&sqlStr, args...)

	var err error
	var total int64
	if session.IsAutoCommit {
		err = session.DB().QueryRow(sqlStr, args...).Scan(&total)
	} else {
		err = session.Tx.QueryRow(sqlStr, args...).Scan(&total)
	}
	if err != nil {
		return 0, err
	}

	return total, nil
}

// Sum call sum some column. bean's non-empty fields are conditions.
func (session *Session) Sum(bean interface{}, columnName string) (float64, error) {
	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	var sqlStr string
	var args []interface{}
	if len(session.Statement.RawSQL) == 0 {
		sqlStr, args = session.Statement.genSumSQL(bean, columnName)
	} else {
		sqlStr = session.Statement.RawSQL
		args = session.Statement.RawParams
	}

	session.queryPreprocess(&sqlStr, args...)

	var err error
	var res float64
	if session.IsAutoCommit {
		err = session.DB().QueryRow(sqlStr, args...).Scan(&res)
	} else {
		err = session.Tx.QueryRow(sqlStr, args...).Scan(&res)
	}
	if err != nil {
		return 0, err
	}

	return res, nil
}

// Sums call sum some columns. bean's non-empty fields are conditions.
func (session *Session) Sums(bean interface{}, columnNames ...string) ([]float64, error) {
	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	var sqlStr string
	var args []interface{}
	if len(session.Statement.RawSQL) == 0 {
		sqlStr, args = session.Statement.genSumSQL(bean, columnNames...)
	} else {
		sqlStr = session.Statement.RawSQL
		args = session.Statement.RawParams
	}

	session.queryPreprocess(&sqlStr, args...)

	var err error
	var res = make([]float64, len(columnNames), len(columnNames))
	if session.IsAutoCommit {
		err = session.DB().QueryRow(sqlStr, args...).ScanSlice(&res)
	} else {
		err = session.Tx.QueryRow(sqlStr, args...).ScanSlice(&res)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

// SumsInt sum specify columns and return as []int64 instead of []float64
func (session *Session) SumsInt(bean interface{}, columnNames ...string) ([]int64, error) {
	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	var sqlStr string
	var args []interface{}
	if len(session.Statement.RawSQL) == 0 {
		sqlStr, args = session.Statement.genSumSQL(bean, columnNames...)
	} else {
		sqlStr = session.Statement.RawSQL
		args = session.Statement.RawParams
	}

	session.queryPreprocess(&sqlStr, args...)

	var err error
	var res = make([]int64, 0, len(columnNames))
	if session.IsAutoCommit {
		err = session.DB().QueryRow(sqlStr, args...).ScanSlice(&res)
	} else {
		err = session.Tx.QueryRow(sqlStr, args...).ScanSlice(&res)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
