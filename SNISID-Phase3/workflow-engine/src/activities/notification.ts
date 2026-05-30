/**
 * SNISID — Multi-channel notifications (SMS, Email, Push, In-App).
 * Routes via the National Notification Gateway.
 */
import axios from 'axios';
const GW = process.env.NOTIF_GATEWAY ?? 'https://notifications.snisid.ht';

export type Channel = 'SMS' | 'EMAIL' | 'PUSH' | 'INAPP';

export async function send(input: {
  channels: Channel[];
  recipient: { nin?: string; phone?: string; email?: string };
  template: string;
  vars: Record<string, string>;
}) {
  await axios.post(`${GW}/send`, input, { timeout: 10_000 });
}
