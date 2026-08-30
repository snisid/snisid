// SNISID: Identity Collision & Biometric Fraud Detection
// This query detects identities sharing highly similar biometric embeddings
// suggesting potential duplicate enrollment or fraud rings.

MATCH (a:Identity)-[:HAS_BIOMETRIC]->(b1:Biometric),
      (c:Identity)-[:HAS_BIOMETRIC]->(b2:Biometric)
WHERE a.id <> c.id
  AND gds.similarity.cosine(b1.embedding, b2.embedding) > 0.98
WITH a, c, gds.similarity.cosine(b1.embedding, b2.embedding) AS similarity
MERGE (a)-[r:POTENTIAL_DUPLICATE {score: similarity}]->(c)
RETURN a.national_id, c.national_id, similarity
ORDER BY similarity DESC;

// Detect Fraud Rings via Shared Device/IP
MATCH (i1:Identity)-[:USES_DEVICE]->(d:Device)<-[:USES_DEVICE]-(i2:Identity)
WHERE i1 <> i2
WITH d, collect(i1.national_id) AS victims, count(i1) AS intensity
WHERE intensity > 3
RETURN d.id AS device_fingerprint, victims, intensity
ORDER BY intensity DESC;
