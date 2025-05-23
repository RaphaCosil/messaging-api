-- Customer
CREATE TABLE customer (
    customer_id    SERIAL PRIMARY KEY,
    username   VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- Chats
CREATE TABLE chat (
    chat_id    SERIAL PRIMARY KEY,
    chat_name  VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- Many-to-many association customer<->chats
CREATE TABLE customer_chat (
    customer_id   INT NOT NULL,
    chat_id   INT NOT NULL,
    role      VARCHAR(50) NOT NULL DEFAULT 'participant',
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (customer_id, chat_id),
    CONSTRAINT fk_uc_customer FOREIGN KEY(customer_id) REFERENCES customer(customer_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_uc_chat FOREIGN KEY(chat_id) REFERENCES chat(chat_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Messages
CREATE TABLE message (
    message_id SERIAL PRIMARY KEY,
    chat_id    INT NOT NULL,
    customer_id    INT NOT NULL,
    content    TEXT NOT NULL,
    sent_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    CONSTRAINT fk_msg_chat FOREIGN KEY(chat_id) REFERENCES chat(chat_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_msg_customer FOREIGN KEY(customer_id) REFERENCES customer(customer_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Insert Customers
INSERT INTO customer (username, password) VALUES 
('carl', '123abc'),
('sigmund', '456def'),
('gustav', '789ghi');

-- Insert chats
INSERT INTO chat (chat_name) VALUES 
('general'), 
('random');

-- Associate Customers & chats
INSERT INTO customer_chat (customer_id, chat_id, role) VALUES 
(1, 1, 'admin'),
(2, 1, 'participant'),
(3, 2, 'participant');

-- Insert messages
INSERT INTO message (chat_id, customer_id, content) VALUES
(1, 1, 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'),
(1, 2, 'Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.'),
(1, 1, 'Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.'),
(1, 2, 'Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore.'),
(2, 3, 'Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia.'),
(2, 3, 'Mollit anim id est laborum.'),
(2, 3, 'Praesent commodo cursus magna, vel scelerisque nisl consectetur.');