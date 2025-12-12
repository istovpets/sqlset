--META
{
    "name": "User Queries",
    "description": "A set of queries for user management."
}
--end

--SQL:GetUserByID
SELECT id, name, email FROM users WHERE id = ?;
--end

--SQL:CreateUser
INSERT INTO users (name, email) VALUES (?, ?);
--end
