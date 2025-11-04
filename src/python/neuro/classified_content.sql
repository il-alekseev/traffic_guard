CREATE TABLE classified_content (
    request_id VARCHAR(50) PRIMARY KEY,
    content_id VARCHAR(50),
    url VARCHAR(255),
    source_content TEXT,
    identified_class VARCHAR(50)
);