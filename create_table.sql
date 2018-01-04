CREATE TABLE todolist (
    id serial PRIMARY KEY,
    created_on DATE,
    due_on DATE,
    status VARCHAR(50),
    description VARCHAR(255)
);