package repository

import "bastion/model"

type BastionRepository struct{}

func (repo *BastionRepository) List(userId int) ([]*model.Machine, error) {
	query := "SELECT " +
		"	v.instance_name, " +
		"	v.private_ip, " +
		"	v.public_ip, " +
		"	v.password " +
		"FROM " +
		"	bastion b " +
		"JOIN " +
		"	vm v ON b.vm_id = v.id " +
		"WHERE " +
		"	b.user_id = ? " +
		"	AND v.is_deleted = 0;"
	rows, err := MysqlClient.Query(query, userId)
	if err != nil {
		return nil, err
	}
	data := make([]*model.Machine, 0)
	for rows.Next() {
		obj := model.Machine{}
		if err := rows.Scan(&obj.InstanceName, &obj.PrivateIP, &obj.PublicIP, &obj.Password); err != nil {
			return nil, err
		}
		data = append(data, &obj)
	}
	return data, nil
}

func (repo *BastionRepository) GetPassword(privateIp string) (string, error) {
	query := "SELECT " +
		"	password " +
		"FROM " +
		"	vm " +
		"WHERE " +
		"	private_ip = ? "
	var encryptText string
	err := MysqlClient.QueryRow(query, privateIp).Scan(&encryptText)
	if err != nil {
		return "", err
	}
	return encryptText, nil
}
