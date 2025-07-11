package repository

type UserRepository struct{}

func (repo *UserRepository) GetUserID(account string) (int, error) {
	query := "SELECT" +
		"	id " +
		"FROM " +
		"	user " +
		"WHERE " +
		"	account = ?"
	var userId int
	err := MysqlClient.QueryRow(query, account).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
