CREATE UNIQUE INDEX `uk_todos_business_duplicate` 
ON `todos` (`todo_list_id`,`title`,`assignee_id`,`due_date`);