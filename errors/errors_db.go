package errors

type NotFoundDbRecordError struct {
	id string
}

type DbRecordLockError struct {
	msg string
}

func (d *DbRecordLockError) Error() string {
	return d.msg
}

func (n *NotFoundDbRecordError) Error() string {
	return "没有找到数据记录:" + n.id
}

func ErrDbNotFoundRecord(id string) error {
	return &NotFoundDbRecordError{id: id}
}

func ErrDbRecordLock() error {
	return &DbRecordLockError{msg: "数据记录已经锁定。"}
}
