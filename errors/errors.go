package errors

import (
	"fmt"
	"strings"
)

type Err struct {
	Code int
	Info string
}

func (e Err) Error() string {
	return e.Info
}

func NewErr(code int, msg string) Err {
	return Err{
		Code: code,
		Info: msg,
	}
}

func NewErrNilDB(name interface{}) Err {
	return Err{Code: 10000, Info: "cannot find db, name=" + fmt.Sprint(name)}
}

func NewErrAlreadyRegistared(name interface{}) Err {
	return Err{Code: 11001, Info: "already registered table name=" + fmt.Sprint(name)}
}

func NewErrSlotSizeMin(min int64) Err {
	return Err{Code: 11002, Info: fmt.Sprintf("minimun slot size must be positive interger, value=%d", min)}
}

func NewErrSlotSizeMax(max, slot int64) Err {
	return Err{Code: 11003, Info: fmt.Sprintf("maximum slot size is out of range, DefinedSize=%d GivenSize=%d", slot, max)}
}

func NewErrSlotMinOverlapped(size int64) Err {
	return Err{Code: 11004, Info: fmt.Sprintf("minimun slot size is overlapped, value=%s", size)}
}

func NewErrSlotMaxOverlapped(size int64) Err {
	return Err{Code: 11005, Info: fmt.Sprintf("maximun slot size is overlapped, value=%s", size)}
}

func NewErrNoSession(name interface{}) Err {
	return Err{Code: 20001, Info: "cannot find session, name=" + fmt.Sprint(name)}
}

func NewErrDuplicateTx() Err {
	return Err{Code: 20002, Info: "transaction already exists"}
}

func NewErrWrongTx() Err {
	return Err{Code: 20003, Info: "something wrong with the transaction"}
}
func NewErrCommitAll(es []error) Err {
	messages := []string{"commit all error: "}
	for _, err := range es {
		messages = append(messages, err.Error())
	}
	return Err{Code: 20004, Info: strings.Join(messages, " ")}
}

func NewErrRollbackAll(es []error) Err {
	messages := []string{"rollback all error: "}
	for _, err := range es {
		messages = append(messages, err.Error())
	}
	return Err{Code: 20005, Info: strings.Join(messages, " ")}
}

func NewErrAnotherTx(name interface{}) Err {
	return Err{Code: 20006, Info: "transaction already exists, db=" + fmt.Sprint(name)}
}

func NewErrParallelQuery(es []error) Err {
	messages := []string{"parallel query error: "}
	for _, err := range es {
		messages = append(messages, err.Error())
	}
	return Err{Code: 30001, Info: strings.Join(messages, " || ")}
}

func NewErrArgType(msg string) Err {
	return Err{Code: 30002, Info: msg}
}
