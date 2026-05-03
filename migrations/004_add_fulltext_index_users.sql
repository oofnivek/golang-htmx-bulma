ALTER TABLE users ADD FULLTEXT INDEX idx_users_search_ngram (email, first_name, last_name, mobile) WITH PARSER ngram;
