
/*
--------------
    schema 
--------------
*/
create schema if not exists emails_hillary;

/*
--------------
     TABLES 
--------------
*/

CREATE TABLE IF NOT EXISTS emails_hillary.emails (
    id INT PRIMARY KEY,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    subject TEXT,
    "from" TEXT,
    "to" TEXT ,
    content TEXT
);

CREATE TABLE IF NOT EXISTS emails_hillary.emails_search (
    id INT PRIMARY KEY,
    search_vector TSVECTOR,
    FOREIGN KEY (id) REFERENCES emails_hillary.emails(id)
);

CREATE INDEX IF NOT EXISTS idx_emails_search_vector
ON emails_hillary.emails_search USING GIN (search_vector);

-- CREATE INVERTED INDEX IF NOT EXISTS idx_emails_search_vector
-- ON emails_hillary.emails_search (search_vector);

create index IF NOT EXISTS idx_emails_date 
on emails_hillary.emails (date);