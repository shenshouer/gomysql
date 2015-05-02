package gomysql

type result struct {
}

func (this *result) LastInsertId() (int64, error) {
	logger.Info("result.LastInsertId")
	return 0, nil
}

func (this *result) RowsAffected() (int64, error) {
	logger.Info("result.RowsAffected")
	return 0, nil
}
