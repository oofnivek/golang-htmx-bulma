-- Migration: Add FULLTEXT n-gram index on vehicle_color.name
-- Required for fast partial string matching via MATCH...AGAINST in boolean mode.
-- The n-gram parser tokenizes the name into small overlapping character sequences,
-- allowing searches like "low" to match "Yellow" without a full table scan.
--
-- Note: The minimum n-gram token length is controlled by innodb_ft_min_token_size
-- (default: 2) or ngram_token_size (default: 2). Ensure your search terms are
-- at least 2 characters long for best results.

ALTER TABLE vehicle_color
    ADD FULLTEXT INDEX idx_color_name_ngram (name) WITH PARSER ngram;
