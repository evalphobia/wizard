package xorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForceNewTransaction(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)

	s, err := orm.ForceNewTransaction(testUser{ID: 1})
	assert.Nil(err)
	assert.NotNil(s)

	assert.EqualValues(3, countUserBySession(s), "initial users count")

	s.Insert(&testUser{ID: 4})
	assert.EqualValues(4, countUserBySession(s), "users count after insert in the transaction")
	assert.EqualValues(3, countUserMaster(orm), "users count after insert not in the transaction")

	err = s.Rollback()
	assert.Nil(err)

	s.Init()
	assert.EqualValues(3, countUserBySession(s), "users count after rollback")
}

func TestTransaction(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)
	xsm := orm.XormSessionManager
	sl := xsm.getOrCreateSessionList(testID)

	assert.Len(sl.getTransactions(), 0)
	s, err := orm.Transaction(testID, testUser{ID: 1})
	assert.Nil(err)
	assert.NotNil(s)
	assert.Len(sl.getTransactions(), 1, "transaction is added")

	assert.EqualValues(3, countUserBySession(s), "initial users count")

	s.Insert(&testUser{ID: 4})
	assert.EqualValues(4, countUserBySession(s), "users count after insert in the transaction")
	assert.EqualValues(3, countUserMaster(orm), "users count after insert not in the transaction")

	err = s.Rollback()
	assert.Nil(err)

	s.Init()
	assert.EqualValues(3, countUserBySession(s), "users count after rollback")
}

func TestTransactionByKey(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)
	xsm := orm.XormSessionManager
	sl := xsm.getOrCreateSessionList(testID)

	assert.Len(sl.getTransactions(), 0)
	s, err := orm.TransactionByKey(testID, testUser{}, 1)
	assert.Nil(err)
	assert.NotNil(s)
	assert.Len(sl.getTransactions(), 1, "transaction is added")

	assert.EqualValues(3, countUserBySession(s), "initial users count")

	s.Insert(&testUser{ID: 4})
	assert.EqualValues(4, countUserBySession(s), "users count after insert in the transaction")
	assert.EqualValues(3, countUserMaster(orm), "users count after insert not in the transaction")

	err = s.Rollback()
	assert.Nil(err)

	s.Init()
	assert.EqualValues(3, countUserBySession(s), "users count after rollback")
}

func TestAutoTransaction(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)
	xsm := orm.XormSessionManager
	sl := xsm.getOrCreateSessionList(testID)
	assert.Len(sl.getTransactions(), 0)

	user1 := testUser{ID: 1}

	s, _ := orm.NewMasterSession(user1)

	err := orm.AutoTransaction(testID, user1, s)
	assert.Nil(err)
	assert.Len(sl.getTransactions(), 0, "transaction is not added")

	orm.SetAutoTransaction(testID, true)
	err = orm.AutoTransaction(testID, user1, s)
	assert.Nil(err)
	assert.Len(sl.getTransactions(), 1, "transaction is added")

	assert.EqualValues(3, countUserBySession(s), "initial users count")
	s.Insert(&testUser{ID: 4})
	assert.EqualValues(4, countUserBySession(s), "users count after insert in the transaction")
	assert.EqualValues(4, countUserMaster(orm), "users count after insert in the transaction")

	s2, _ := newSession(orm.Master(user1), user1)
	assert.EqualValues(3, countUserBySession(s2), "users count after insert in another session")

	err = s.Rollback()
	assert.Nil(err)
	s.Init()
	assert.EqualValues(3, countUserBySession(s), "users count after rollback")

}

func TestAutoTransactionDuplicateTx(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)
	xsm := orm.XormSessionManager
	sl := xsm.getOrCreateSessionList(testID)
	assert.Len(sl.getTransactions(), 0)

	var err error

	orm.SetAutoTransaction(testID, true)
	s1, _ := orm.NewMasterSession(testUser{ID: 1})
	s2, _ := orm.NewMasterSession(testUser{ID: 500})
	sl.addTransaction(orm.Master(testUser{ID: 1}), s1)

	err = orm.AutoTransaction(testID, testUser{ID: 1}, s1)
	assert.Nil(err, "error does not occur if same session exists")

	err = orm.AutoTransaction(testID, testUser{ID: 1}, s2)
	assert.NotNil(err, "error occurs if another session exists")
}

func TestCommitAll(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)
	xsm := orm.XormSessionManager
	sl := xsm.getOrCreateSessionList(testID)
	assert.Len(sl.getTransactions(), 0)

	user1 := testUser{ID: 1}
	user500 := testUser{ID: 500}

	s1, _ := orm.NewMasterSession(user1)
	s2, _ := orm.NewMasterSession(user500)

	orm.SetAutoTransaction(testID, true)
	orm.AutoTransaction(testID, user1, s1)
	orm.AutoTransaction(testID, user500, s2)
	assert.Len(sl.getTransactions(), 2, "transaction is added")

	assert.EqualValues(3, countUserBySession(s1), "initial users count")
	assert.EqualValues(3, countUserBySession(s2), "initial users count")

	num, err := s1.Insert(&testUser{ID: 4})
	assert.Nil(err)
	assert.EqualValues(1, num)
	s2.Insert(&testUser{ID: 504})
	assert.EqualValues(4, countUserBySession(s1), "users count after insert in the transaction")
	assert.EqualValues(4, countUserBySession(s2), "users count after insert in the transaction")
	assert.EqualValues(4, countUserMaster(orm), "users count after insert in the transaction")
	assert.EqualValues(4, countUserMasterB(orm), "users count after insert in the transaction")

	orm.SetAutoTransaction(testID, false)
	s1b, _ := newSession(orm.Master(user1), user1)
	s2b, _ := newSession(orm.Master(user500), user500)
	assert.EqualValues(3, countUserBySession(s1b), "users count after insert  in another session")
	assert.EqualValues(3, countUserBySession(s2b), "users count after insert  in another session")

	orm.ReadOnly(testID, true)
	err = orm.CommitAll(testID)
	assert.Nil(err)
	assert.Len(sl.getTransactions(), 2, "transaction is not removed when readonly")

	orm.ReadOnly(testID, false)
	err = orm.CommitAll(testID)
	assert.Nil(err)
	assert.Len(sl.getTransactions(), 0, "transaction is removed")

	assert.EqualValues(4, countUserMaster(orm), "users count after commit")
	assert.EqualValues(4, countUserMasterB(orm), "users count after commit")

	initTestDB()
}

func TestRollbackAll(t *testing.T) {
	assert := assert.New(t)
	wiz := testCreateWizard()
	orm := New(wiz)
	xsm := orm.XormSessionManager
	sl := xsm.getOrCreateSessionList(testID)
	assert.Len(sl.getTransactions(), 0)

	user1 := testUser{ID: 1}
	user500 := testUser{ID: 500}

	s1, _ := orm.NewMasterSession(user1)
	s2, _ := orm.NewMasterSession(user500)

	orm.SetAutoTransaction(testID, true)
	orm.AutoTransaction(testID, user1, s1)
	orm.AutoTransaction(testID, user500, s2)
	assert.Len(sl.getTransactions(), 2, "transaction is added")

	assert.EqualValues(3, countUserBySession(s1), "initial users count")
	assert.EqualValues(3, countUserBySession(s2), "initial users count")

	s1.Insert(&testUser{ID: 4})
	s2.Insert(&testUser{ID: 504})
	assert.EqualValues(4, countUserBySession(s1), "users count after insert in the transaction")
	assert.EqualValues(4, countUserBySession(s2), "users count after insert in the transaction")
	assert.EqualValues(4, countUserMaster(orm), "users count after insert in the transaction")
	assert.EqualValues(4, countUserMasterB(orm), "users count after insert in the transaction")

	orm.SetAutoTransaction(testID, false)
	s1b, _ := newSession(orm.Master(user1), user1)
	s2b, _ := newSession(orm.Master(user500), user500)
	assert.EqualValues(3, countUserBySession(s1b), "users count after insert  in another session")
	assert.EqualValues(3, countUserBySession(s2b), "users count after insert  in another session")

	orm.ReadOnly(testID, true)
	err := orm.RollbackAll(testID)
	assert.Nil(err)
	assert.Len(sl.getTransactions(), 2, "transaction is not removed when readonly")
	assert.EqualValues(4, countUserBySession(s1), "rollback does not occur when read only")
	assert.EqualValues(4, countUserBySession(s2), "rollback does not occur when read only")

	orm.ReadOnly(testID, false)
	err = orm.RollbackAll(testID)
	assert.Nil(err)
	assert.Len(sl.getTransactions(), 0, "transaction is removed")

	assert.EqualValues(3, countUserMaster(orm), "users count after rollback")
	assert.EqualValues(3, countUserMasterB(orm), "users count after rollback")
}

func TestTransactionNilDB(t *testing.T) {
	assert := assert.New(t)
	orm := New(emptyWiz)

	var s Session
	var err error

	s, err = orm.ForceNewTransaction(testUser{ID: 1})
	assert.NotNil(err)
	assert.Nil(s)

	s, err = orm.Transaction(testID, testUser{ID: 1})
	assert.NotNil(err)
	assert.Nil(s)

	s, err = orm.TransactionByKey(testID, testUser{}, 1)
	assert.NotNil(err)
	assert.Nil(s)
}
