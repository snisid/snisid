import { useState, useCallback } from 'react';
import { FlatList, Pressable, RefreshControl, StyleSheet, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { SymbolView } from 'expo-symbols';

import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';

import { BottomTabInset, MaxContentWidth, Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

type LogEntry = {
  id: string;
  event: string;
  category: string;
  severity: 'critical' | 'warning' | 'info';
  timestamp: string;
  user: string;
};

const ALL_LOGS: LogEntry[] = [
  { id: '1', event: 'Citizen #SN-88472 verified', category: 'Verification', severity: 'info', timestamp: '2026-06-03 09:42', user: 'admin@nia.gov.ng' },
  { id: '2', event: 'Suspicious login from IP 203.0.113.42', category: 'Auth', severity: 'critical', timestamp: '2026-06-03 09:15', user: 'system' },
  { id: '3', event: 'Biometric enrollment completed', category: 'Biometric', severity: 'info', timestamp: '2026-06-03 08:30', user: 'oprator@lagos.gov.ng' },
  { id: '4', event: 'Document #DOC-4454 expired', category: 'Document', severity: 'warning', timestamp: '2026-06-03 07:00', user: 'system' },
  { id: '5', event: 'Identity data sync with NIA completed', category: 'System', severity: 'info', timestamp: '2026-06-03 06:00', user: 'system' },
  { id: '6', event: 'Batch verification of 156 records', category: 'Verification', severity: 'info', timestamp: '2026-06-02 23:00', user: 'batch@system' },
  { id: '7', event: 'Citizen #SN-77231 suspended (court order)', category: 'Status', severity: 'critical', timestamp: '2026-06-02 18:22', user: 'admin@nia.gov.ng' },
  { id: '8', event: 'Failed login attempt (user unknown)', category: 'Auth', severity: 'warning', timestamp: '2026-06-02 16:45', user: 'anonymous' },
  { id: '9', event: 'Admin session terminated (timeout)', category: 'Auth', severity: 'info', timestamp: '2026-06-02 15:30', user: 'auditor@fct.gov.ng' },
  { id: '10', event: 'Verification request #VR-4456 auto-rejected', category: 'Verification', severity: 'warning', timestamp: '2026-06-02 14:10', user: 'system' },
];

const SEVERITY_CONFIG = { critical: { color: '#c62828', bg: '#fce4ec' }, warning: { color: '#e65100', bg: '#fff3e0' }, info: { color: '#1565c0', bg: '#e3f2fd' } };
const CATEGORIES = ['All', 'Verification', 'Auth', 'Biometric', 'Document', 'System', 'Status'];
const SEVERITIES = ['All', 'critical', 'warning', 'info'];

export default function AuditLogScreen() {
  const theme = useTheme();
  const router = useRouter();
  const [refreshing, setRefreshing] = useState(false);
  const [category, setCategory] = useState('All');
  const [severity, setSeverity] = useState('All');
  const [dateRange, setDateRange] = useState('All');

  const filtered = ALL_LOGS.filter(l => {
    const matchCat = category === 'All' || l.category === category;
    const matchSev = severity === 'All' || l.severity === severity;
    return matchCat && matchSev;
  });

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    setTimeout(() => setRefreshing(false), 1500);
  }, []);

  const ChipRow = ({ items, selected, onSelect }: { items: string[]; selected: string; onSelect: (v: string) => void }) => (
    <View style={styles.chipRow}>
      {items.map(item => (
        <Pressable key={item} onPress={() => onSelect(item)} style={[styles.chip, { backgroundColor: selected === item ? theme.text : theme.backgroundSelected }]}>
          <ThemedText style={{ color: selected === item ? theme.background : theme.text, fontSize: 12, textTransform: 'capitalize' }}>{item}</ThemedText>
        </Pressable>
      ))}
    </View>
  );

  return (
    <SafeAreaView style={[styles.safeArea, { backgroundColor: theme.background }]}>
      <FlatList
        data={filtered}
        keyExtractor={(item) => item.id}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
        contentContainerStyle={styles.listContent}
        ListHeaderComponent={
          <ThemedView>
            <View style={styles.headerRow}>
              <Pressable onPress={() => router.back()} style={styles.backBtn}>
                <SymbolView name={{ ios: 'chevron.left', android: 'arrow_back', web: 'arrow_back' }} size={20} tintColor={theme.text} />
              </Pressable>
              <ThemedText type="subtitle" style={{ fontSize: 22, lineHeight: 28 }}>Audit Log</ThemedText>
            </View>

            <View style={styles.filterSection}>
              <ThemedText themeColor="textSecondary" style={{ fontSize: 12, marginBottom: Spacing.one }}>Date Range</ThemedText>
              <ChipRow items={['All', 'Today', 'This Week', 'This Month']} selected={dateRange} onSelect={setDateRange} />
            </View>
            <View style={styles.filterSection}>
              <ThemedText themeColor="textSecondary" style={{ fontSize: 12, marginBottom: Spacing.one }}>Event Type</ThemedText>
              <ChipRow items={CATEGORIES} selected={category} onSelect={setCategory} />
            </View>
            <View style={styles.filterSection}>
              <ThemedText themeColor="textSecondary" style={{ fontSize: 12, marginBottom: Spacing.one }}>Severity</ThemedText>
              <ChipRow items={SEVERITIES} selected={severity} onSelect={setSeverity} />
            </View>

            <ThemedText themeColor="textSecondary" style={{ fontSize: 12, paddingHorizontal: Spacing.four, paddingTop: Spacing.two, paddingBottom: Spacing.two }}>
              {filtered.length} entries
            </ThemedText>
          </ThemedView>
        }
        renderItem={({ item }) => {
          const cfg = SEVERITY_CONFIG[item.severity];
          return (
            <View style={styles.entryWrapper}>
              <ThemedView type="backgroundElement" style={styles.entry}>
                <View style={[styles.sevDot, { backgroundColor: cfg.bg }]}>
                  <View style={[styles.sevInner, { backgroundColor: cfg.color }]} />
                </View>
                <View style={styles.entryContent}>
                  <ThemedText style={{ fontSize: 13 }}>{item.event}</ThemedText>
                  <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>{item.category} · {item.user}</ThemedText>
                </View>
                <View style={[styles.severityBadge, { backgroundColor: cfg.bg }]}>
                  <ThemedText style={{ color: cfg.color, fontSize: 10, fontWeight: '700' }}>{item.severity}</ThemedText>
                </View>
                <ThemedText themeColor="textSecondary" style={{ fontSize: 11, minWidth: 60, textAlign: 'right' }}>{item.timestamp}</ThemedText>
              </ThemedView>
            </View>
          );
        }}
        ListEmptyComponent={
          <ThemedView style={styles.empty}>
            <ThemedText themeColor="textSecondary" style={{ textAlign: 'center' }}>No log entries match your filters.</ThemedText>
          </ThemedView>
        }
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1 },
  listContent: { paddingBottom: BottomTabInset + Spacing.three, maxWidth: MaxContentWidth, width: '100%', alignSelf: 'center' },
  headerRow: { flexDirection: 'row', alignItems: 'center', gap: Spacing.two, paddingHorizontal: Spacing.four, paddingTop: Spacing.three, paddingBottom: Spacing.two },
  backBtn: { padding: Spacing.one },
  filterSection: { paddingHorizontal: Spacing.four, marginBottom: Spacing.two },
  chipRow: { flexDirection: 'row', flexWrap: 'wrap', gap: Spacing.two },
  chip: { paddingHorizontal: Spacing.three, paddingVertical: Spacing.one, borderRadius: Spacing.five },
  entryWrapper: { paddingHorizontal: Spacing.four, marginBottom: Spacing.two },
  entry: { flexDirection: 'row', alignItems: 'center', padding: Spacing.three, borderRadius: Spacing.two, gap: Spacing.two },
  sevDot: { width: 24, height: 24, borderRadius: 12, justifyContent: 'center', alignItems: 'center' },
  sevInner: { width: 10, height: 10, borderRadius: 5 },
  entryContent: { flex: 1, gap: 1 },
  severityBadge: { paddingHorizontal: Spacing.one, paddingVertical: 2, borderRadius: Spacing.one },
  empty: { padding: Spacing.six, alignItems: 'center' },
});
