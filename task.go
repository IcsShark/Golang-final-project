package main

import "context"

type Item struct {
	ID        int
	Title     string
	Completed bool
}

type Tasks struct {
	Items          []Item
	Count          int
	CompletedCount int
}

func fetchTasks() ([]Item, error) {
	var items []Item
	rows, err := DB.Query("select id, title, completed from tasks order by position;")
	if err != nil {
		return []Item{}, err
	}
	defer rows.Close()
	for rows.Next() {
		item := Item{}
		err := rows.Scan(&item.ID, &item.Title, &item.Completed)
		if err != nil {
			return []Item{}, err
		}
		items = append(items, item)
	}
	return items, nil
}

func fetchTask(ID int) (Item, error) {
	var item Item
	err := DB.QueryRow("select id, title, completed from tasks where id = (?)", ID).Scan(&item.ID, &item.Title, &item.Completed)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func updateTask(ID int, title string) (Item, error) {
	_, err := DB.Exec("update tasks set title = (?) where id = (?)", title, ID)
	if err != nil {
		return Item{}, err
	}
	var item Item
	err = DB.QueryRow("select id, title, completed from tasks where id = (?)", ID).Scan(&item.ID, &item.Title, &item.Completed)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func fetchCount() (int, error) {
	var count int
	err := DB.QueryRow("select count(*) from tasks;").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func fetchCompletedCount() (int, error) {
	var count int
	err := DB.QueryRow("select count(*) from tasks where completed = 1;").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func insertTask(title string) (Item, error) {
	tx, err := DB.Begin()
	if err != nil {
		return Item{}, err
	}
	count, err := fetchCount()
	if err != nil {
		tx.Rollback()
		return Item{}, err
	}

	var id int
	err = tx.QueryRow("insert into tasks (title, position) values (?, ?) returning id", title, count).Scan(&id)
	if err != nil {
		tx.Rollback()
		return Item{}, err
	}
	item := Item{ID: id, Title: title, Completed: false}

	err = tx.Commit()
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func deleteTask(ctx context.Context, ID int) error {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.Exec("delete from tasks where id = (?)", ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows, err := tx.Query("select id, title, completed from tasks order by position")
	if err != nil {
		tx.Rollback()
		return err
	}
	var ids []int
	for rows.Next() {
		var id int
		var title string
		var completed bool
		err := rows.Scan(&id, &title, &completed)
		if err != nil {
			tx.Rollback()
			return err
		}
		ids = append(ids, id)
	}
	for idx, id := range ids {
		_, err := tx.Exec("update tasks set position = (?) where id = (?)", idx, id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func orderTasks(ctx context.Context, values []int) error {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for i, v := range values {
		_, err := tx.Exec("update tasks set position = (?) where id = (?)", i, v)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func toggleTask(ID int) (Item, error) {
	var item Item
	err := DB.QueryRow("update tasks set completed = case when completed = 1 then 0 else 1 end where id = (?) returning id, title, completed", ID).Scan(&item.ID, &item.Title, &item.Completed)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}
