import { useState, useCallback } from 'react';
import { FlatList, Pressable, RefreshControl, StyleSheet, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { SymbolView } from 'expo-symbols';

import { StatCard } from '@/components/StatCard';
import { AuditLogEntry } from '@/components/AuditLogEntry';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';

import { BottomTabInset, MaxContentWidth, Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

type DashboardStat = {
  title: string;
  value: string;
  icon: { ios: string; android: string; web: string };
  trend?: { value: string; positive: boolean };
  color: string;
};

const STATS: DashboardStat[] = [
  { title: 'Total Identities', value: '1,284,732', icon: { ios: 'person.crop.rectangle', android: 'badge', web: 'badge' }, trend: { value: '2.4%', positive: true }, color: '#1565c0' },
  { title: 'Pending Verifications', value: '3,847', icon: { ios: 'clock', android: 'pending', web: 'pending' }, trend: { value: '12%', positive: false }, color: '#e65100' },
  { title: 'Active Alerts', value: '23', icon: { ios: 'exclamationmark.triangle', android: 'warning', web: 'warning' }, trend: { value: '5', positive: false }, color: '#c62828' },
  { title: "Today's Registrations", value: '1,042', icon: { ios: 'person.badge.plus', android: 'person_add', web: 'person_add' }, trend: { value: '8.1%', positive: true }, color: '#2e7d32' },
];

type RecentActivity = {
  id: string;
  icon: { ios: string; android: string; web: string };
  description: string;
  timestamp: string;
  severity: 'critical' | 'warning' | 'info';
};

const RECENT_ACTIVITIES: RecentActivity[] = [
  { id: '1', icon: { ios: 'person.fill', android: 'person', web: 'person' }, description: 'Citizen #SN-88472 verified successfully', timestamp: '2 min ago', severity: 'info' },
  { id: '2', icon: { ios: 'exclamationmark.shield', android: 'security', web: 'security' }, description: 'Suspicious login attempt detected from IP 203.0.113.42', timestamp: '15 min ago', severity: 'critical' },
  { id: '3', icon: { ios: 'doc.text.fill', android: 'description', web: 'description' }, description: 'Biometric enrollment for citizen #SN-99103 completed', timestamp: '32 min ago', severity: 'info' },
  { id: '4', icon: { ios: 'arrow.triangle.branch', android: 'sync', web: 'sync' }, description: 'Identity data sync with agency NIA completed', timestamp: '1 hour ago', severity: 'info' },
  { id: '5', icon: { ios: 'exclamationmark.triangle', android: 'warning', web: 'warning' }, description: 'Verification request #VR-4456 auto-rejected (document expired)', timestamp: '2 hours ago', severity: 'warning' },
  { id: '6', icon: { ios: 'xmark.shield', android: 'block', web: 'block' }, description: 'Citizen #SN-77231 suspended by admin per court order', timestamp: '3 hours ago', severity: 'critical' },
  { id: '7', icon: { ios: 'checkmark.circle', android: 'check_circle', web: 'check_circle' }, description: 'Batch verification of 156 records completed', timestamp: '5 hours ago', severity: 'info' },
];

export default function DashboardScreen() {
  const theme = useTheme();
  const router = useRouter();
  const [refreshing, setRefreshing] = useState(false);

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    setTimeout(() => setRefreshing(false), 1500);
  }, []);

  return (
    <SafeAreaView style={[styles.safeArea, { backgroundColor: theme.background }]}>
      <FlatList
        data={RECENT_ACTIVITIES}
        keyExtractor={(item) => item.id}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
        contentContainerStyle={styles.listContent}
        ListHeaderComponent={
          <ThemedView>
            <View style={styles.headerRow}>
              <ThemedText type="subtitle" style={{ fontSize: 26, lineHeight: 32 }}>Dashboard</ThemedText>
              <Pressable onPress={() => router.push('/alerts')} style={styles.alertBell}>
                <SymbolView name={{ ios: 'bell.badge', android: 'notifications', web: 'notifications' }} size={22} tintColor={theme.text} />
              </Pressable>
            </View>

            <View style={styles.statsGrid}>
              {STATS.map((stat, i) => (
                <StatCard key={i} {...stat} />
              ))}
            </View>

            <View style={styles.sectionHeader}>
              <ThemedText style={{ fontWeight: '600', fontSize: 17 }}>Recent Activity</ThemedText>
              <Pressable onPress={() => router.push('/audit-log')}>
                <ThemedText themeColor="textSecondary" style={{ fontSize: 13 }}>View All →</ThemedText>
              </Pressable>
            </View>
          </ThemedView>
        }
        renderItem={({ item }) => (
          <View style={styles.activityItem}>
            <AuditLogEntry {...item} />
          </View>
        )}
        ListFooterComponent={
          <View style={styles.footer}>
            <ThemedView type="backgroundElement" style={styles.quickActions}>
              <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>Quick Actions</ThemedText>
              <View style={styles.actionRow}>
                <Pressable onPress={() => router.push('/reports')} style={[styles.quickBtn, { backgroundColor: theme.backgroundSelected }]}>
                  <SymbolView name={{ ios: 'chart.bar', android: 'bar_chart', web: 'bar_chart' }} size={20} tintColor={theme.text} />
                  <ThemedText style={{ fontSize: 12 }}>Reports</ThemedText>
                </Pressable>
                <Pressable onPress={() => router.push('/alerts')} style={[styles.quickBtn, { backgroundColor: theme.backgroundSelected }]}>
                  <SymbolView name={{ ios: 'bell', android: 'notifications', web: 'notifications' }} size={20} tintColor={theme.text} />
                  <ThemedText style={{ fontSize: 12 }}>Alerts</ThemedText>
                </Pressable>
                <Pressable onPress={() => router.push('/(tabs)/citizens')} style={[styles.quickBtn, { backgroundColor: theme.backgroundSelected }]}>
                  <SymbolView name={{ ios: 'person.2', android: 'people', web: 'people' }} size={20} tintColor={theme.text} />
                  <ThemedText style={{ fontSize: 12 }}>Citizens</ThemedText>
                </Pressable>
                <Pressable onPress={() => router.push('/(tabs)/verification')} style={[styles.quickBtn, { backgroundColor: theme.backgroundSelected }]}>
                  <SymbolView name={{ ios: 'checkmark.shield', android: 'verified_user', web: 'verified_user' }} size={20} tintColor={theme.text} />
                  <ThemedText style={{ fontSize: 12 }}>Verify</ThemedText>
                </Pressable>
              </View>
            </ThemedView>
          </View>
        }
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
  },
  listContent: {
    paddingBottom: BottomTabInset + Spacing.three,
    maxWidth: MaxContentWidth,
    width: '100%',
    alignSelf: 'center',
  },
  headerRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: Spacing.four,
    paddingTop: Spacing.three,
    paddingBottom: Spacing.two,
  },
  alertBell: {
    padding: Spacing.two,
  },
  statsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    paddingHorizontal: Spacing.four,
    gap: Spacing.two,
    marginBottom: Spacing.four,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: Spacing.four,
    marginBottom: Spacing.two,
  },
  activityItem: {
    paddingHorizontal: Spacing.four,
    marginBottom: Spacing.two,
  },
  footer: {
    paddingHorizontal: Spacing.four,
    marginTop: Spacing.two,
  },
  quickActions: {
    padding: Spacing.three,
    borderRadius: Spacing.three,
    gap: Spacing.three,
  },
  actionRow: {
    flexDirection: 'row',
    gap: Spacing.two,
  },
  quickBtn: {
    flex: 1,
    paddingVertical: Spacing.three,
    borderRadius: Spacing.two,
    alignItems: 'center',
    gap: Spacing.one,
  },
});
