import { useState, useCallback } from 'react';
import { ScrollView, Pressable, RefreshControl, StyleSheet, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { SymbolView } from 'expo-symbols';

import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';

import { BottomTabInset, MaxContentWidth, Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

type AuthAttempt = { id: string; type: 'success' | 'failure'; method: string; timestamp: string; ip: string };
type AdminSession = { id: string; admin: string; role: string; lastActive: string; ip: string };

const RECENT_ATTEMPTS: AuthAttempt[] = [
  { id: '1', type: 'failure', method: 'password', timestamp: '2 min ago', ip: '203.0.113.42' },
  { id: '2', type: 'success', method: 'biometric', timestamp: '15 min ago', ip: '198.51.100.10' },
  { id: '3', type: 'failure', method: 'password', timestamp: '18 min ago', ip: '203.0.113.42' },
  { id: '4', type: 'failure', method: 'password', timestamp: '22 min ago', ip: '203.0.113.42' },
  { id: '5', type: 'success', method: 'smartcard', timestamp: '1 hour ago', ip: '192.0.2.55' },
];

const ACTIVE_SESSIONS: AdminSession[] = [
  { id: '1', admin: 'admin@nia.gov.ng', role: 'Super Admin', lastActive: 'Now', ip: '192.168.1.100' },
  { id: '2', admin: 'oprator@lagos.gov.ng', role: 'Operator', lastActive: '2 min ago', ip: '10.0.0.45' },
  { id: '3', admin: 'auditor@fct.gov.ng', role: 'Auditor', lastActive: '5 min ago', ip: '10.0.0.88' },
];

export default function SecurityScreen() {
  const theme = useTheme();
  const [refreshing, setRefreshing] = useState(false);

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    setTimeout(() => setRefreshing(false), 1500);
  }, []);

  const failedCount = RECENT_ATTEMPTS.filter(a => a.type === 'failure').length;

  return (
    <SafeAreaView style={[styles.safeArea, { backgroundColor: theme.background }]}>
      <ScrollView
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
        contentContainerStyle={styles.scrollContent}
      >
        <View style={styles.header}>
          <ThemedText type="subtitle" style={{ fontSize: 26, lineHeight: 32 }}>Security</ThemedText>
        </View>

        <ThemedView type="backgroundElement" style={styles.section}>
          <View style={styles.sectionHeader}>
            <SymbolView name={{ ios: 'exclamationmark.shield', android: 'security', web: 'security' }} size={18} tintColor="#c62828" />
            <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>Failed Login Alerts</ThemedText>
            <View style={[styles.badge, { backgroundColor: '#fce4ec' }]}>
              <ThemedText style={{ color: '#c62828', fontSize: 12, fontWeight: '700' }}>{failedCount} alerts</ThemedText>
            </View>
          </View>
          <ThemedText themeColor="textSecondary" style={{ fontSize: 13 }}>
            {failedCount >= 3 ? 'Multiple failed attempts detected from IP 203.0.113.42. Consider blocking.' : 'No unusual activity detected.'}
          </ThemedText>
        </ThemedView>

        <ThemedView type="backgroundElement" style={styles.section}>
          <View style={styles.sectionHeader}>
            <SymbolView name={{ ios: 'clock.arrow.circlepath', android: 'history', web: 'history' }} size={18} tintColor={theme.text} />
            <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>Recent Auth Attempts</ThemedText>
          </View>
          {RECENT_ATTEMPTS.map(a => (
            <View key={a.id} style={styles.attemptRow}>
              <View style={[styles.dot, { backgroundColor: a.type === 'success' ? '#2e7d32' : '#c62828' }]} />
              <View style={styles.attemptInfo}>
                <ThemedText style={{ fontSize: 13, textTransform: 'capitalize' }}>{a.type} · {a.method}</ThemedText>
                <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>{a.ip}</ThemedText>
              </View>
              <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>{a.timestamp}</ThemedText>
            </View>
          ))}
        </ThemedView>

        <ThemedView type="backgroundElement" style={styles.section}>
          <View style={styles.sectionHeader}>
            <SymbolView name={{ ios: 'person.2', android: 'people', web: 'people' }} size={18} tintColor={theme.text} />
            <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>Active Admin Sessions</ThemedText>
          </View>
          {ACTIVE_SESSIONS.map(s => (
            <View key={s.id} style={styles.sessionRow}>
              <View style={[styles.onlineDot, { backgroundColor: '#2e7d32' }]} />
              <View style={styles.sessionInfo}>
                <ThemedText style={{ fontSize: 13, fontWeight: '500' }}>{s.admin}</ThemedText>
                <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>{s.role} · {s.ip}</ThemedText>
              </View>
              <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>{s.lastActive}</ThemedText>
            </View>
          ))}
        </ThemedView>

        <ThemedView type="backgroundElement" style={styles.section}>
          <View style={styles.sectionHeader}>
            <SymbolView name={{ ios: 'doc.text', android: 'description', web: 'description' }} size={18} tintColor={theme.text} />
            <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>Security Policy</ThemedText>
          </View>
          <View style={styles.policyRow}>
            <ThemedText style={{ fontSize: 13 }}>Password Policy</ThemedText>
            <ThemedText themeColor="textSecondary" style={{ fontSize: 12 }}>min 12 chars, 2FA required</ThemedText>
          </View>
          <View style={styles.policyRow}>
            <ThemedText style={{ fontSize: 13 }}>Session Timeout</ThemedText>
            <ThemedText themeColor="textSecondary" style={{ fontSize: 12 }}>15 min inactivity</ThemedText>
          </View>
          <View style={styles.policyRow}>
            <ThemedText style={{ fontSize: 13 }}>Max Login Attempts</ThemedText>
            <ThemedText themeColor="textSecondary" style={{ fontSize: 12 }}>5 before lockout</ThemedText>
          </View>
          <View style={styles.policyRow}>
            <ThemedText style={{ fontSize: 13 }}>IP Whitelist</ThemedText>
            <ThemedText themeColor="textSecondary" style={{ fontSize: 12 }}>Enabled (3 ranges)</ThemedText>
          </View>
        </ThemedView>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1 },
  scrollContent: { paddingBottom: BottomTabInset + Spacing.three, maxWidth: MaxContentWidth, width: '100%', alignSelf: 'center', paddingHorizontal: Spacing.four },
  header: { paddingTop: Spacing.three, paddingBottom: Spacing.three },
  section: { padding: Spacing.three, borderRadius: Spacing.three, marginBottom: Spacing.three, gap: Spacing.two },
  sectionHeader: { flexDirection: 'row', alignItems: 'center', gap: Spacing.two },
  badge: { paddingHorizontal: Spacing.two, paddingVertical: 2, borderRadius: Spacing.one, marginLeft: 'auto' },
  attemptRow: { flexDirection: 'row', alignItems: 'center', gap: Spacing.two, paddingVertical: Spacing.one },
  dot: { width: 8, height: 8, borderRadius: 4 },
  attemptInfo: { flex: 1, gap: 1 },
  sessionRow: { flexDirection: 'row', alignItems: 'center', gap: Spacing.two, paddingVertical: Spacing.one },
  onlineDot: { width: 8, height: 8, borderRadius: 4 },
  sessionInfo: { flex: 1, gap: 1 },
  policyRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', paddingVertical: Spacing.one },
});
