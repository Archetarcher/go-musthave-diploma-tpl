package repositories

const (
	userCreateQuery     = "insert into users ( login, hash) values (:login, :hash) returning id"
	userUpdateQuery     = "update users set  balance = :balance where id = :id"
	userGetByLoginQuery = "SELECT id, login, hash from users where login = $1 "
	userGetByIDQuery    = "SELECT id, login, balance from users where id = $1 "

	orderAccrualGetByIDQuery           = "SELECT * from order_accrual where order_id = $1 "
	orderAccrualGetByUserIDQuery       = "SELECT * from order_accrual where user_id = $1 and order_id = $2 "
	orderAccrualGetAllByUserIDQuery    = "SELECT * from order_accrual where user_id = $1  order by id desc"
	orderAccrualGetOrdersByStatusQuery = "SELECT * from order_accrual where status in (?)"
	orderAccrualCreateQuery            = "insert into order_accrual ( user_id, order_id, status, amount) values (:user_id, :order_id, :status, :amount)  returning id"
	orderAccrualUpdateQuery            = "update  order_accrual set amount = :amount, status = :status where order_id = :order_id"

	orderWithdrawalCreateQuery            = "insert into order_withdrawal ( user_id, order_id, amount) values (:user_id, :order_id, :amount)  returning id"
	orderWithdrawalGetAllByUserIDQuery    = "SELECT * from order_withdrawal where user_id = $1 order by id desc "
	orderWithdrawalGetAllByUserSumIDQuery = "SELECT sum(amount) from order_withdrawal where user_id = $1 group by user_id"
	orderWithdrawalGetByUserIDQuery       = "SELECT * from order_withdrawal where user_id = $1 and order_id = $2 "
)
