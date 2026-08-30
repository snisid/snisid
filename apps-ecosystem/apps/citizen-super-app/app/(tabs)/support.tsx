import { useState, useCallback, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  ActivityIndicator,
  Linking,
  TextInput,
  Alert,
  RefreshControl,
} from 'react-native';
import { useRouter } from 'expo-router';

import { useAuth } from '@/lib/auth';
import type { FAQItem, SupportTicket } from '@/lib/api';

const TICKET_STATUS_COLORS: Record<string, { bg: string; text: string; label: string }> = {
  open: { bg: '#DBEAFE', text: '#1E40AF', label: 'Open' },
  in_progress: { bg: '#FEF3C7', text: '#92400E', label: 'In Progress' },
  resolved: { bg: '#DCFCE7', text: '#166534', label: 'Resolved' },
  closed: { bg: '#F3F4F6', text: '#6B7280', label: 'Closed' },
};

function FAQAccordion({ item }: { item: FAQItem }) {
  const [open, setOpen] = useState(false);

  return (
    <View style={styles.faqItem}>
      <TouchableOpacity
        style={styles.faqHeader}
        onPress={() => setOpen((p) => !p)}
        activeOpacity={0.7}>
        <Text style={styles.faqQuestion}>{item.question}</Text>
        <Text style={styles.faqChevron}>{open ? '▲' : '▼'}</Text>
      </TouchableOpacity>
      {open && (
        <View style={styles.faqBody}>
          <Text style={styles.faqAnswer}>{item.answer}</Text>
        </View>
      )}
    </View>
  );
}

export default function SupportScreen() {
  const router = useRouter();
  const { api } = useAuth();
  const [faq, setFaq] = useState<FAQItem[]>([]);
  const [tickets, setTickets] = useState<SupportTicket[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState('');
  const [showNewTicket, setShowNewTicket] = useState(false);
  const [ticketSubject, setTicketSubject] = useState('');
  const [ticketMessage, setTicketMessage] = useState('');
  const [ticketCategory, setTicketCategory] = useState('general');
  const [submitting, setSubmitting] = useState(false);

  const fetchData = useCallback(async () => {
    try {
      setError('');
      const [faqRes, ticketRes] = await Promise.all([
        api.getSupportFAQ(),
        api.getSupportTickets({ limit: 10 }),
      ]);
      setFaq(faqRes.faq);
      setTickets(ticketRes.tickets);
    } catch (err: any) {
      setError(
        err?.response?.data?.message ?? err?.message ?? 'Failed to load support data'
      );
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, [api]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    fetchData();
  }, [fetchData]);

  const handleCall = () => {
    Linking.openURL('tel:+27123456789').catch(() =>
      Alert.alert('Error', 'Unable to make a phone call on this device')
    );
  };

  const handleEmail = () => {
    Linking.openURL('mailto:support@snisid.gov.za').catch(() =>
      Alert.alert('Error', 'Unable to open email client')
    );
  };

  const handleSubmitTicket = useCallback(async () => {
    if (!ticketSubject.trim() || !ticketMessage.trim()) {
      Alert.alert('Validation Error', 'Please fill in all fields');
      return;
    }
    setSubmitting(true);
    try {
      const res = await api.createSupportTicket({
        subject: ticketSubject.trim(),
        message: ticketMessage.trim(),
        category: ticketCategory,
      });
      setTickets((prev) => [res.ticket, ...prev]);
      setShowNewTicket(false);
      setTicketSubject('');
      setTicketMessage('');
      setTicketCategory('general');
      Alert.alert('Ticket Created', `Your ticket #${res.ticket.id} has been submitted.`);
    } catch (err: any) {
      Alert.alert(
        'Error',
        err?.response?.data?.message ?? err?.message ?? 'Failed to create ticket'
      );
    } finally {
      setSubmitting(false);
    }
  }, [ticketSubject, ticketMessage, ticketCategory, api]);

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#0033a0" />
        <Text style={styles.loadingText}>Loading support information...</Text>
      </View>
    );
  }

  if (error && faq.length === 0) {
    return (
      <View style={styles.center}>
        <Text style={styles.errorIcon}>⚠️</Text>
        <Text style={styles.errorText}>{error}</Text>
        <TouchableOpacity style={styles.retryButton} onPress={fetchData}>
          <Text style={styles.retryText}>Retry</Text>
        </TouchableOpacity>
      </View>
    );
  }

  const categories = [...new Set(faq.map((f) => f.category))];

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.content}
      refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}>
      <Text style={styles.pageTitle}>Support</Text>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Contact Us</Text>
        <View style={styles.contactRow}>
          <TouchableOpacity style={styles.contactCard} onPress={handleCall}>
            <Text style={styles.contactIcon}>📞</Text>
            <Text style={styles.contactLabel}>Call</Text>
            <Text style={styles.contactValue}>012 345 6789</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.contactCard} onPress={handleEmail}>
            <Text style={styles.contactIcon}>📧</Text>
            <Text style={styles.contactLabel}>Email</Text>
            <Text style={styles.contactValue}>support@snisid.gov.za</Text>
          </TouchableOpacity>
        </View>
      </View>

      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Frequently Asked Questions</Text>
        </View>
        {faq.length === 0 ? (
          <Text style={styles.emptyText}>No FAQs available at this time.</Text>
        ) : (
          categories.map((cat) => (
            <View key={cat} style={styles.faqCategory}>
              <Text style={styles.categoryLabel}>{cat}</Text>
              {faq
                .filter((f) => f.category === cat)
                .map((item) => (
                  <FAQAccordion key={item.id} item={item} />
                ))}
            </View>
          ))
        )}
      </View>

      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>
            My Tickets ({tickets.length})
          </Text>
          <TouchableOpacity
            style={styles.newTicketButton}
            onPress={() => setShowNewTicket(!showNewTicket)}>
            <Text style={styles.newTicketText}>
              {showNewTicket ? 'Cancel' : 'New Ticket'}
            </Text>
          </TouchableOpacity>
        </View>

        {showNewTicket && (
          <View style={styles.newTicketForm}>
            <Text style={styles.formLabel}>Category</Text>
            <View style={styles.categoryRow}>
              {['general', 'account', 'technical', 'identity'].map((cat) => (
                <TouchableOpacity
                  key={cat}
                  style={[
                    styles.categoryChip,
                    ticketCategory === cat && styles.categoryChipActive,
                  ]}
                  onPress={() => setTicketCategory(cat)}>
                  <Text
                    style={[
                      styles.categoryChipText,
                      ticketCategory === cat && styles.categoryChipTextActive,
                    ]}>
                    {cat.charAt(0).toUpperCase() + cat.slice(1)}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>

            <Text style={styles.formLabel}>Subject</Text>
            <TextInput
              style={styles.input}
              value={ticketSubject}
              onChangeText={setTicketSubject}
              placeholder="Brief summary of your issue"
              placeholderTextColor="#999"
            />

            <Text style={styles.formLabel}>Message</Text>
            <TextInput
              style={[styles.input, styles.textArea]}
              value={ticketMessage}
              onChangeText={setTicketMessage}
              placeholder="Describe your issue in detail..."
              placeholderTextColor="#999"
              multiline
              numberOfLines={4}
              textAlignVertical="top"
            />

            <TouchableOpacity
              style={[styles.submitButton, submitting && styles.submitButtonDisabled]}
              onPress={handleSubmitTicket}
              disabled={submitting}>
              {submitting ? (
                <ActivityIndicator size="small" color="#fff" />
              ) : (
                <Text style={styles.submitButtonText}>Submit Ticket</Text>
              )}
            </TouchableOpacity>
          </View>
        )}

        {tickets.length === 0 ? (
          <Text style={styles.emptyText}>No support tickets yet.</Text>
        ) : (
          tickets.map((ticket) => {
            const st = TICKET_STATUS_COLORS[ticket.status] ?? TICKET_STATUS_COLORS.open;
            return (
              <View key={ticket.id} style={styles.ticketCard}>
                <View style={styles.ticketHeader}>
                  <Text style={styles.ticketSubject} numberOfLines={1}>
                    {ticket.subject}
                  </Text>
                  <View style={[styles.ticketStatusBadge, { backgroundColor: st.bg }]}>
                    <Text style={[styles.ticketStatusText, { color: st.text }]}>
                      {st.label}
                    </Text>
                  </View>
                </View>
                <Text style={styles.ticketDate}>
                  Created {new Date(ticket.createdAt).toLocaleDateString()}
                  {' · '}Priority: {ticket.priority}
                </Text>
                {ticket.lastMessage && (
                  <Text style={styles.ticketMessage} numberOfLines={2}>
                    {ticket.lastMessage}
                  </Text>
                )}
              </View>
            );
          })
        )}
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#F5F5F7',
  },
  content: {
    paddingTop: 24,
    paddingBottom: 40,
    gap: 24,
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
    backgroundColor: '#F5F5F7',
  },
  pageTitle: {
    fontSize: 28,
    fontWeight: '700',
    color: '#000',
    paddingHorizontal: 20,
  },
  section: {
    paddingHorizontal: 16,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  sectionTitle: {
    fontSize: 13,
    fontWeight: '600',
    color: '#888',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
    paddingHorizontal: 4,
    marginBottom: 8,
  },
  loadingText: {
    marginTop: 12,
    fontSize: 14,
    color: '#666',
  },
  errorIcon: {
    fontSize: 48,
    marginBottom: 12,
  },
  errorText: {
    fontSize: 14,
    color: '#991B1B',
    textAlign: 'center',
    marginBottom: 16,
  },
  retryButton: {
    backgroundColor: '#0033a0',
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 8,
  },
  retryText: {
    color: '#fff',
    fontWeight: '600',
  },
  contactRow: {
    flexDirection: 'row',
    gap: 12,
  },
  contactCard: {
    flex: 1,
    backgroundColor: '#fff',
    borderRadius: 14,
    padding: 20,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 4,
    elevation: 2,
  },
  contactIcon: {
    fontSize: 32,
    marginBottom: 8,
  },
  contactLabel: {
    fontSize: 15,
    fontWeight: '600',
    color: '#000',
  },
  contactValue: {
    fontSize: 11,
    color: '#666',
    marginTop: 4,
  },
  faqCategory: {
    marginBottom: 12,
  },
  categoryLabel: {
    fontSize: 14,
    fontWeight: '700',
    color: '#0033a0',
    marginBottom: 6,
    paddingHorizontal: 4,
    textTransform: 'capitalize',
  },
  faqItem: {
    backgroundColor: '#fff',
    borderRadius: 12,
    marginBottom: 4,
    overflow: 'hidden',
  },
  faqHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
  },
  faqQuestion: {
    flex: 1,
    fontSize: 15,
    fontWeight: '500',
    color: '#000',
    marginRight: 12,
  },
  faqChevron: {
    fontSize: 10,
    color: '#999',
  },
  faqBody: {
    paddingHorizontal: 16,
    paddingBottom: 16,
  },
  faqAnswer: {
    fontSize: 14,
    color: '#555',
    lineHeight: 20,
  },
  newTicketButton: {
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 8,
    backgroundColor: '#0033a0',
  },
  newTicketText: {
    color: '#fff',
    fontSize: 13,
    fontWeight: '600',
  },
  newTicketForm: {
    backgroundColor: '#fff',
    borderRadius: 14,
    padding: 16,
    marginBottom: 12,
    gap: 12,
  },
  formLabel: {
    fontSize: 13,
    fontWeight: '600',
    color: '#555',
  },
  categoryRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  categoryChip: {
    paddingHorizontal: 14,
    paddingVertical: 8,
    borderRadius: 20,
    backgroundColor: '#F0F0F5',
    borderWidth: 1,
    borderColor: '#E0E0E0',
  },
  categoryChipActive: {
    backgroundColor: '#0033a0',
    borderColor: '#0033a0',
  },
  categoryChipText: {
    fontSize: 13,
    color: '#555',
    fontWeight: '500',
  },
  categoryChipTextActive: {
    color: '#fff',
  },
  input: {
    borderWidth: 1,
    borderColor: '#E5E5EA',
    borderRadius: 10,
    padding: 12,
    fontSize: 15,
    color: '#000',
    backgroundColor: '#FAFAFA',
  },
  textArea: {
    minHeight: 100,
    paddingTop: 12,
  },
  submitButton: {
    backgroundColor: '#0033a0',
    borderRadius: 10,
    paddingVertical: 14,
    alignItems: 'center',
  },
  submitButtonDisabled: {
    opacity: 0.6,
  },
  submitButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  ticketCard: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 14,
    marginBottom: 6,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.03,
    shadowRadius: 3,
    elevation: 1,
  },
  ticketHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 6,
  },
  ticketSubject: {
    flex: 1,
    fontSize: 15,
    fontWeight: '600',
    color: '#000',
    marginRight: 8,
  },
  ticketStatusBadge: {
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: 8,
  },
  ticketStatusText: {
    fontSize: 11,
    fontWeight: '700',
  },
  ticketDate: {
    fontSize: 11,
    color: '#999',
    marginBottom: 4,
  },
  ticketMessage: {
    fontSize: 13,
    color: '#666',
    lineHeight: 18,
  },
  emptyText: {
    fontSize: 14,
    color: '#999',
    textAlign: 'center',
    paddingVertical: 16,
  },
});
