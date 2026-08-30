import { useState, useCallback } from 'react';
import { FlatList, Pressable, RefreshControl, StyleSheet, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { SymbolView } from 'expo-symbols';

import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';

import { BottomTabInset, MaxContentWidth, Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

type AlertItem = {
  id: string;
  severity: 'critical' | 'warning' | 'info';
  description: string;
  timestamp: string;
  acknowledged: boolean;
};

const INITIAL_ALERTS: AlertItem[] = [
  { id: '1', severity: 'critical', description: 'Multiple failed login attempts detected from IP 203.0.113.42', timestamp: '15 min ago', acknowledged: false },
  { id: '2', severity: 'critical', description: 'Biometric verification timeout — system degraded in Kano region', timestamp: '1 hour ago', acknowledged: false },
  { id: '3', severity: 'warning', description: 'Document batch #DOC-9981 has 12 expired records', timestamp: '2 hours ago', acknowledged: false },
  { id: '4', severity: 'warning', description: 'Data sync with agency NIA delayed (queue backlog)', timestamp: '3 hours ago', acknowledged: false },
  { id: '5', severity: 'warning', description: 'Certificate for api.nia.gov.ng expires in 7 days', timestamp: '5 hours ago', acknowledged: false },
  { id: '6', severity: 'info', description: 'Scheduled maintenance: 2025-06-05 02:00-04:00 UTC', timestamp: '1 day ago', acknowledged: false },
  { id: '7', severity: 'info', description: 'Admin audit trail review completed — 0 discrepancies', timestamp: '2 days ago', acknowledged: true },
  { id: '8', severity: 'info', description: 'Batch verification of 156 records completed successfully', timestamp: '3 days ago', acknowledged: true },
];

const SEVERITY_ORDER = { critical: 0, warning: 1, info: 2 };
const SEVERITY_CONFIG = { critical: { color: '#c62828', bg: '#fce4ec', label: 'Critical' }, warning: { color: '#e65100', bg: '#fff3e0', label: 'Warning' }, info: { color: '#1565c0', bg: '#e3f2fd', label: 'Info' } };

export default function AlertsScreen() {
  const theme = useTheme();
  const router = useRouter();
  const [refreshing, setRefreshing] = useState(false);
  const [alerts, setAlerts] = useState(INITIAL_ALERTS);

  const sorted = [...alerts].sort((a, b) => SEVERITY_ORDER[a.severity] - SEVERITY_ORDER[b.severity]);
  const activeCritical = sorted.filter(a => !a.acknowledged && a.severity === 'critical');
  const activeWarnings = sorted.filter(a => !a.acknowledged && a.severity === 'warning');
  const resolved = sorted.filter(a => a.acknowledged);

  const sections = [
    { title: 'Active Critical Alerts', data: activeCritical, icon: 'exclamationmark.triangle' as const },
    { title: 'Active Warnings', data: activeWarnings, icon: 'exclamationmark.circle' as const },
    { title: 'Resolved Alerts', data: resolved, icon: 'checkmark.circle' as const },
  ];

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    setTimeout(() => setRefreshing(false), 1500);
  }, []);

  const toggleAcknowledge = (id: string) => {
    setAlerts(prev => prev.map(a => a.id === id ? { ...a, acknowledged: !a.acknowledged } : a));
  };

  const flatData = sections.flatMap((s, si) => [{ _type: 'header' as const, title: s.title, icon: s.icon, count: s.data.length }, ...s.data.map(d => ({ _type: 'item' as const, ...d }))]);

  return (
    <SafeAreaView style={[styles.safeArea, { backgroundColor: theme.background }]}>
      <FlatList
        data={flatData}
        keyExtractor={(item) => item._type === 'header' ? `h-${item.title}` : `i-${item.id}`}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
        contentContainerStyle={styles.listContent}
        ListHeaderComponent={
          <View style={styles.header}>
            <Pressable onPress={() => router.back()} style={styles.backBtn}>
              <SymbolView name={{ ios: 'chevron.left', android: 'arrow_back', web: 'arrow_back' }} size={20} tintColor={theme.text} />
            </Pressable>
            <ThemedText type="subtitle" style={{ fontSize: 22, lineHeight: 28 }}>Alert Center</ThemedText>
          </View>
        }
        renderItem={({ item }) => {
          if (item._type === 'header') {
            return (
              <View style={styles.sectionHeader}>
                <SymbolView name={{ ios: item.icon, android: item.icon === 'exclamationmark.triangle' ? 'warning' : item.icon === 'exclamationmark.circle' ? 'error' : 'check_circle', web: item.icon }} size={16} tintColor={theme.text} />
                <ThemedText style={{ fontWeight: '600', fontSize: 14 }}>{item.title}</ThemedText>
                <View style={[styles.countBadge, { backgroundColor: theme.backgroundSelected }]}>
                  <ThemedText style={{ fontSize: 11, fontWeight: '600' }}>{item.count}</ThemedText>
                </View>
              </View>
            );
          }
          const cfg = SEVERITY_CONFIG[item.severity];
          return (
            <ThemedView type="backgroundElement" style={[styles.alertCard, { borderLeftColor: cfg.color, borderLeftWidth: 3 }]}>
              <View style={styles.alertContent}>
                <View style={styles.alertHeader}>
                  <View style={[styles.severityBadge, { backgroundColor: cfg.bg }]}>
                    <ThemedText style={{ color: cfg.color, fontSize: 10, fontWeight: '700' }}>{cfg.label}</ThemedText>
                  </View>
                  <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>{item.timestamp}</ThemedText>
                </View>
                <ThemedText style={{ fontSize: 13, marginTop: Spacing.one }}>{item.description}</ThemedText>
                <Pressable onPress={() => toggleAcknowledge(item.id)} style={[styles.ackBtn, { backgroundColor: theme.backgroundSelected }]}>
                  <SymbolView name={item.acknowledged ? { ios: 'arrow.uturn.left', android: 'undo', web: 'undo' } : { ios: 'checkmark', android: 'check', web: 'check' }} size={12} tintColor={theme.text} />
                  <ThemedText style={{ fontSize: 11 }}>{item.acknowledged ? 'Reopen' : 'Acknowledge'}</ThemedText>
                </Pressable>
              </View>
            </ThemedView>
          );
        }}
        ListEmptyComponent={
          <ThemedView style={styles.empty}>
            <SymbolView name={{ ios: 'bell.slash', android: 'notifications_off', web: 'notifications_off' }} size={40} tintColor={theme.textSecondary} />
            <ThemedText themeColor="textSecondary" style={{ textAlign: 'center', marginTop: Spacing.two }}>No alerts</ThemedText>
          </ThemedView>
        }
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1 },
  listContent: { paddingBottom: BottomTabInset + Spacing.three, maxWidth: MaxContentWidth, width: '100%', alignSelf: 'center' },
  header: { flexDirection: 'row', alignItems: 'center', gap: Spacing.two, paddingHorizontal: Spacing.four, paddingTop: Spacing.three, paddingBottom: Spacing.three },
  backBtn: { padding: Spacing.one },
  sectionHeader: { flexDirection: 'row', alignItems: 'center', gap: Spacing.two, paddingHorizontal: Spacing.four, paddingTop: Spacing.three, paddingBottom: Spacing.two },
  countBadge: { paddingHorizontal: Spacing.two, paddingVertical: 1, borderRadius: Spacing.one },
  alertCard: { marginHorizontal: Spacing.four, marginBottom: Spacing.two, borderRadius: Spacing.two },
  alertContent: { padding: Spacing.three, gap: Spacing.one },
  alertHeader: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' },
  severityBadge: { paddingHorizontal: Spacing.one, paddingVertical: 2, borderRadius: Spacing.one },
  ackBtn: { flexDirection: 'row', alignItems: 'center', alignSelf: 'flex-end', gap: Spacing.half, paddingHorizontal: Spacing.two, paddingVertical: Spacing.one, borderRadius: Spacing.one, marginTop: Spacing.one },
  empty: { padding: Spacing.six, alignItems: 'center' },
});
