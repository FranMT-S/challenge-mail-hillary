CREATE TABLE IF NOT EXISTS emails_hillary (
    id INT PRIMARY KEY,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    subject TEXT,
    "from" TEXT,
    "to" TEXT,
    content TEXT
);