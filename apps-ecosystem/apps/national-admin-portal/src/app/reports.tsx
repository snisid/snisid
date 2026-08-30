import { ScrollView, Pressable, StyleSheet, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { SymbolView } from 'expo-symbols';

import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';

import { BottomTabInset, MaxContentWidth, Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

type Report = {
  id: string;
  title: string;
  icon: { ios: string; android: string; web: string };
  stats: { label: string; value: string }[];
  trend: string;
  trendPositive: boolean;
};

const REPORTS: Report[] = [
  {
    id: '1', title: 'Registration Stats', icon: { ios: 'person.badge.plus', android: 'person_add', web: 'person_add' },
    stats: [{ label: 'Daily', value: '1,042' }, { label: 'Weekly', value: '7,891' }, { label: 'Monthly', value: '32,450' }],
    trend: '8.1% vs last month', trendPositive: true,
  },
  {
    id: '2', title: 'Verification Success Rate', icon: { ios: 'checkmark.shield', android: 'verified_user', web: 'verified_user' },
    stats: [{ label: 'Success Rate', value: '97.3%' }, { label: 'Total Requests', value: '4,215' }, { label: 'Failed', value: '114' }],
    trend: '1.2% improvement', trendPositive: true,
  },
  {
    id: '3', title: 'Document Issuance', icon: { ios: 'doc.text.fill', android: 'description', web: 'description' },
    stats: [{ label: 'Issued Today', value: '856' }, { label: 'Pending', value: '234' }, { label: 'Rejected', value: '42' }],
    trend: '3.4% vs yesterday', trendPositive: true,
  },
  {
    id: '4', title: 'Biometric Enrollment', icon: { ios: 'touchid', android: 'fingerprint', web: 'fingerprint' },
    stats: [{ label: 'Enrolled', value: '1,241,890' }, { label: 'Pending', value: '42,842' }, { label: 'Completion', value: '96.7%' }],
    trend: '0.5% increase', trendPositive: true,
  },
];

export default function ReportsScreen() {
  const theme = useTheme();
  const router = useRouter();

  return (
    <SafeAreaView style={[styles.safeArea, { backgroundColor: theme.background }]}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <View style={styles.headerRow}>
          <Pressable onPress={() => router.back()} style={styles.backBtn}>
            <SymbolView name={{ ios: 'chevron.left', android: 'arrow_back', web: 'arrow_back' }} size={20} tintColor={theme.text} />
          </Pressable>
          <ThemedText type="subtitle" style={{ fontSize: 22, lineHeight: 28 }}>Reports &amp; Analytics</ThemedText>
        </View>

        {REPORTS.map(report => (
          <ThemedView key={report.id} type="backgroundElement" style={styles.section}>
            <View style={styles.sectionHeader}>
              <SymbolView name={report.icon} size={18} tintColor={theme.text} />
              <ThemedText style={{ fontWeight: '600', fontSize: 15, flex: 1 }}>{report.title}</ThemedText>
              <Pressable style={[styles.exportBtn, { backgroundColor: theme.backgroundSelected }]}>
                <SymbolView name={{ ios: 'square.and.arrow.up', android: 'share', web: 'share' }} size={14} tintColor={theme.text} />
                <ThemedText style={{ fontSize: 11 }}>Export</ThemedText>
              </Pressable>
            </View>

            <View style={styles.statsRow}>
              {report.stats.map((s, i) => (
                <View key={i} style={styles.statItem}>
                  <ThemedText style={{ fontWeight: '700', fontSize: 20 }}>{s.value}</ThemedText>
                  <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>{s.label}</ThemedText>
                </View>
              ))}
            </View>

            <View style={styles.chartPlaceholder}>
              <View style={styles.barSet}>
                {[60, 85, 45, 90, 70, 55, 80].map((h, i) => (
                  <View key={i} style={[styles.bar, { height: h * 0.6, backgroundColor: theme.text + '30' }]} />
                ))}
              </View>
            </View>

            <View style={[styles.trendBadge, { backgroundColor: report.trendPositive ? '#e8f5e9' : '#fce4ec' }]}>
              <ThemedText style={{ color: report.trendPositive ? '#2e7d32' : '#c62828', fontSize: 11, fontWeight: '600' }}>
                {report.trendPositive ? '↑' : '↓'} {report.trend}
              </ThemedText>
            </View>
          </ThemedView>
        ))}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1 },
  scrollContent: { paddingBottom: BottomTabInset + Spacing.three, maxWidth: MaxContentWidth, width: '100%', alignSelf: 'center', paddingHorizontal: Spacing.four },
  headerRow: { flexDirection: 'row', alignItems: 'center', gap: Spacing.two, paddingTop: Spacing.three, paddingBottom: Spacing.three },
  backBtn: { padding: Spacing.one },
  section: { padding: Spacing.three, borderRadius: Spacing.three, marginBottom: Spacing.three, gap: Spacing.three },
  sectionHeader: { flexDirection: 'row', alignItems: 'center', gap: Spacing.two },
  exportBtn: { flexDirection: 'row', alignItems: 'center', gap: Spacing.half, paddingHorizontal: Spacing.two, paddingVertical: Spacing.one, borderRadius: Spacing.one },
  statsRow: { flexDirection: 'row', gap: Spacing.three },
  statItem: { flex: 1, alignItems: 'center', gap: 2 },
  chartPlaceholder: { height: 80, justifyContent: 'flex-end' },
  barSet: { flexDirection: 'row', alignItems: 'flex-end', justifyContent: 'space-around', height: '100%' },
  bar: { width: 20, borderRadius: Spacing.one },
  trendBadge: { paddingHorizontal: Spacing.two, paddingVertical: Spacing.one, borderRadius: Spacing.one, alignSelf: 'flex-start' },
});
