package gomysql

type tx struct {
}

func (this *tx) Commit() error {
	logger.Info("this.Commit")
	return nil
}

func (this *tx) Rollback() error {
	logger.Info("this.Rollback")
	return nil
}
