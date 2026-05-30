import axios from 'axios';

interface AuditEvent {
  id: string;
  action: string;
  details: Record<string, any>;
  timestamp: string;
  url: string;
  userAgent: string;
}

// Offline queue stored in IndexedDB (simulated here via localStorage since logs are non-PII operator metadata)
// In a true implementation, idb-keyval would be better, but localStorage works for small queues.
const AUDIT_QUEUE_KEY = 'snisid_audit_queue';

const getQueue = (): AuditEvent[] => {
  try {
    const q = localStorage.getItem(AUDIT_QUEUE_KEY);
    return q ? JSON.parse(q) : [];
  } catch {
    return [];
  }
};

const saveQueue = (queue: AuditEvent[]) => {
  localStorage.setItem(AUDIT_QUEUE_KEY, JSON.stringify(queue));
};

export const logAuditAction = async (action: string, details: Record<string, any> = {}) => {
  const event: AuditEvent = {
    id: crypto.randomUUID(),
    action,
    details,
    timestamp: new Date().toISOString(),
    url: window.location.href,
    userAgent: navigator.userAgent
  };

  // Try to send immediately if online
  if (navigator.onLine) {
    try {
      await axios.post('/api/v1/audit/logs', event, {
        headers: { 'Content-Type': 'application/json' },
        // Token interceptor handles Authorization
      });
      return;
    } catch (err) {
      console.warn('Failed to send audit log, queuing for offline sync', err);
    }
  }

  // Queue if offline or request failed
  const queue = getQueue();
  queue.push(event);
  saveQueue(queue);
};

// Background sync function
export const syncAuditLogs = async () => {
  if (!navigator.onLine) return;

  const queue = getQueue();
  if (queue.length === 0) return;

  try {
    // Send bulk request
    await axios.post('/api/v1/audit/logs/bulk', { events: queue });
    // Clear queue on success
    saveQueue([]);
    console.log(`Successfully synced ${queue.length} audit logs`);
  } catch (err) {
    console.error('Failed to sync audit logs, will retry later', err);
  }
};

// Auto-sync when coming back online
if (typeof window !== 'undefined') {
  window.addEventListener('online', syncAuditLogs);
  
  // Also try periodically
  setInterval(() => {
    if (navigator.onLine && getQueue().length > 0) {
      syncAuditLogs();
    }
  }, 60000); // Check every minute
}
