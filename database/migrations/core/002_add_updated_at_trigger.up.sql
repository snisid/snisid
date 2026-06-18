CREATE OR REPLACE FUNCTION snisid_core.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_citizens_updated_at
    BEFORE UPDATE ON snisid_core.citizens
    FOR EACH ROW
    EXECUTE FUNCTION snisid_core.update_updated_at_column();
