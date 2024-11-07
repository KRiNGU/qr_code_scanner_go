ALTER TABLE transactions ADD processed BOOLEAN;
UPDATE transactions SET processed = FALSE;