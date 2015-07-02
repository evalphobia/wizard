package errors

import(
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
	return Err{ Code: 10000, Info: "cannot find db, name=" + fmt.Sprint(name)}
}

func NewErrNoSession(name interface{}) Err {
	return Err{ Code: 20001, Info: "cannot find session, name=" + fmt.Sprint(name)}
}

func NewErrDuplicateTx() Err {
	return Err{ Code: 20002, Info: "transaction already exists" }
}

func NewErrWrongTx() Err {
	return Err{ Code: 20003, Info: "something wrong with the transaction" }
}
func NewErrCommitAll(es []error) Err {
	messages := []string{"commit all error: "}
	for _, err := range es {
		messages = append(messages, err.Error())
	}
	return Err{ Code: 20004, Info: strings.Join(messages, " ") }
}
