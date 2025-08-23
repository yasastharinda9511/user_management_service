CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       first_name VARCHAR(50) NOT NULL,
                       last_name VARCHAR(50) NOT NULL,
                       phone VARCHAR(20),
                       is_active BOOLEAN DEFAULT TRUE,
                       is_email_verified BOOLEAN DEFAULT FALSE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       last_login TIMESTAMP NULL
);

-- Create indexes for users table
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_active ON users(is_active);


-- User sessions table for JWT token management
-- Create the user_sessions table
CREATE TABLE user_sessions (
                               id SERIAL PRIMARY KEY,
                               user_id INT NOT NULL,
                               token_hash VARCHAR(255) NOT NULL,
                               expires_at TIMESTAMP NOT NULL,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               is_revoked BOOLEAN DEFAULT FALSE,

                               FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes separately
CREATE INDEX idx_user_id ON user_sessions(user_id);
CREATE INDEX idx_token_hash ON user_sessions(token_hash);
CREATE INDEX idx_expires_at ON user_sessions(expires_at);


CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) UNIQUE NOT NULL,
                       description TEXT,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_roles (
                            user_id INT NOT NULL,
                            role_id INT NOT NULL,
                            assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                            PRIMARY KEY (user_id, role_id),
                            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                            FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE TABLE password_reset_tokens (
                                       id SERIAL PRIMARY KEY,
                                       user_id INT NOT NULL,
                                       token VARCHAR(255) NOT NULL,
                                       expires_at TIMESTAMP NOT NULL,
                                       used BOOLEAN DEFAULT FALSE,
                                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_token ON password_reset_tokens(token);
CREATE INDEX idx_token_user_id ON password_reset_tokens(user_id);




CREATE TABLE permissions (
                             id SERIAL PRIMARY KEY,
                             name VARCHAR(100) UNIQUE NOT NULL,
                             resource VARCHAR(50) NOT NULL,        -- e.g., 'users', 'orders', 'products'
                             action VARCHAR(50) NOT NULL,          -- e.g., 'create', 'read', 'update', 'delete'
                             description TEXT,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE role_permissions (
                                  id SERIAL PRIMARY KEY,
                                  role_id INT NOT NULL,
                                  permission_id INT NOT NULL,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
                                  FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
                                  UNIQUE(role_id, permission_id)
);


INSERT INTO user_roles (user_id, role_id) VALUES
    (13, 1);

INSERT INTO roles (name, description) VALUES
                                          ('admin', 'Administrator with full access'),
                                          ('user', 'Regular user with basic access'),
                                          ('moderator', 'Moderator with limited admin access');

INSERT INTO permissions (name, resource, action, description) VALUES
                                                                  ('users.create', 'users', 'create', 'Create new user accounts'),
                                                                  ('users.read', 'users', 'read', 'View user profiles and information'),
                                                                  ('users.update', 'users', 'update', 'Update user profiles and information'),
                                                                  ('users.delete', 'users', 'delete', 'Delete user accounts'),
                                                                  ('users.activate', 'users', 'activate', 'Activate/deactivate user accounts'),
                                                                  ('users.reset_password', 'users', 'reset_password', 'Reset user passwords'),
                                                                  ('users.impersonate', 'users', 'impersonate', 'Login as another user');

INSERT INTO role_permissions (role_id, permission_id)
SELECT
    r.id as role_id,
    p.id as permission_id
FROM roles r
         CROSS JOIN permissions p
WHERE r.name = 'admin';