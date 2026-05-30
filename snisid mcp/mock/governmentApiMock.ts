import 'dotenv/config';
import express from 'express';

const app = express();
const port = Number(process.env['MOCK_GOV_PORT'] ?? 4000);

app.disable('x-powered-by');
app.use(express.json({ limit: '1mb' }));

app.get('/healthz', (_req, res) => {
  res.json({
    ok: true,
    service: 'snisid-mock-government-apis',
    port
  });
});

app.use((req, res, next) => {
  const expectedApiKey = process.env['GOV_API_KEY'];

  if (!expectedApiKey) {
    res.status(500).json({
      error: 'GOV_API_KEY_NOT_CONFIGURED'
    });
    return;
  }

  const receivedApiKey = req.header('x-api-key');

  if (receivedApiKey !== expectedApiKey) {
    res.status(401).json({
      error: 'INVALID_GOV_API_KEY'
    });
    return;
  }

  next();
});

const endpoints = [
  '/oni/identity/verify',
  '/oni/identity/profile',
  '/oni/civil/birth-certificate',
  '/oni/identity/nationality',
  '/oni/identity/risk',

  '/dgi/tax/nif/verify',
  '/dgi/tax/compliance',
  '/dgi/tax/business-registry',
  '/dgi/tax/financial-risk',

  '/pnh/police/wanted-person',
  '/pnh/police/incidents',
  '/pnh/police/gang-affiliation',
  '/pnh/police/weapon-permit',
  '/pnh/police/threat-monitoring',

  '/mjsp/justice/criminal-record',
  '/mjsp/justice/warrants',
  '/mjsp/justice/court-cases',
  '/mjsp/justice/detention-status',
  '/mjsp/justice/history',

  '/immigration/immigration/border-alerts',
  '/immigration/immigration/travel-history',
  '/immigration/immigration/visa',
  '/immigration/immigration/entry-exit',
  '/immigration/immigration/watchlist-scan',

  '/anh/archives/lookup',
  '/anh/archives/attestation',

  '/passport/passport/lookup',
  '/passport/passport/status',

  '/biometric/biometric/match',
  '/biometric/biometric/face/verify',

  '/education/education/student/verify',
  '/education/education/diploma/verify',
  '/education/education/institution/lookup',
  '/education/education/academic-history',

  '/intelligence/intelligence/fusion-analysis',
  '/intelligence/intelligence/risk-score',
  '/intelligence/intelligence/network-analysis',
  '/intelligence/intelligence/threat-detection',
  '/intelligence/intelligence/behavior-analysis'
] as const;

function buildMockResponse(
  path: string,
  body: unknown,
  headers: Record<string, string | undefined>
) {
  return {
    ok: true,
    mock: true,
    path,
    timestamp: new Date().toISOString(),
    correlationId: headers['x-correlation-id'],
    purpose: headers['x-purpose'],
    actor: {
      subject: headers['x-actor-subject'],
      ministry: headers['x-actor-ministry']
    },
    request: body,
    data: {
      status: 'MOCK_SUCCESS',
      verified: true,
      confidence: 0.97,
      riskScore: 12,
      riskFlags: [],
      dataMinimized: true,
      note: 'Réponse simulée pour environnement local de développement SNISID.'
    }
  };
}

for (const endpoint of endpoints) {
  app.post(endpoint, (req, res) => {
    res.json(
      buildMockResponse(endpoint, req.body, {
        'x-correlation-id': req.header('x-correlation-id') ?? undefined,
        'x-purpose': req.header('x-purpose') ?? undefined,
        'x-actor-subject': req.header('x-actor-subject') ?? undefined,
        'x-actor-ministry': req.header('x-actor-ministry') ?? undefined
      })
    );
  });
}

app.use((req, res) => {
  res.status(404).json({
    error: 'MOCK_ENDPOINT_NOT_FOUND',
    method: req.method,
    path: req.path
  });
});

app.listen(port, () => {
  console.log(`SNISID mock government APIs listening on http://localhost:${port}`);
});
