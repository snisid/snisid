import { useState } from 'react';
import { ScrollView, Pressable, StyleSheet, View, Switch } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { SymbolView } from 'expo-symbols';

import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';

import { BottomTabInset, MaxContentWidth, Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

type ToggleKey = 'emailAlerts' | 'smsAlerts' | 'pushAlerts' | 'criticalAlerts' | 'dailyDigest';

export default function SettingsScreen() {
  const theme = useTheme();
  const [toggles, setToggles] = useState<Record<ToggleKey, boolean>>({
    emailAlerts: true, smsAlerts: false, pushAlerts: true, criticalAlerts: true, dailyDigest: false,
  });

  const toggle = (key: ToggleKey) => setToggles(prev => ({ ...prev, [key]: !prev[key] }));

  const SettingRow = ({ label, value }: { label: string; value: string }) => (
    <View style={styles.settingRow}>
      <ThemedText style={{ fontSize: 13 }}>{label}</ThemedText>
      <ThemedText themeColor="textSecondary" style={{ fontSize: 12 }}>{value}</ThemedText>
    </View>
  );

  return (
    <SafeAreaView style={[styles.safeArea, { backgroundColor: theme.background }]}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <View style={styles.header}>
          <ThemedText type="subtitle" style={{ fontSize: 26, lineHeight: 32 }}>Settings</ThemedText>
        </View>

        <ThemedView type="backgroundElement" style={styles.section}>
          <View style={styles.sectionHeader}>
            <SymbolView name={{ ios: 'gearshape.2', android: 'settings', web: 'settings' }} size={18} tintColor={theme.text} />
            <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>System Configuration</ThemedText>
          </View>
          <SettingRow label="Region" value="Nigeria (NG)" />
          <SettingRow label="Session Timeout" value="15 minutes" />
          <SettingRow label="Backup Interval" value="Every 6 hours" />
          <SettingRow label="Data Sync" value="Real-time" />
        </ThemedView>

        <ThemedView type="backgroundElement" style={styles.section}>
          <View style={styles.sectionHeader}>
            <SymbolView name={{ ios: 'bell', android: 'notifications', web: 'notifications' }} size={18} tintColor={theme.text} />
            <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>Notification Preferences</ThemedText>
          </View>
          {([
            { key: 'emailAlerts' as ToggleKey, label: 'Email Alerts', desc: 'Security and system notifications via email' },
            { key: 'smsAlerts' as ToggleKey, label: 'SMS Alerts', desc: 'Critical alerts via SMS' },
            { key: 'pushAlerts' as ToggleKey, label: 'Push Notifications', desc: 'Real-time push notifications' },
            { key: 'criticalAlerts' as ToggleKey, label: 'Critical Alerts Only', desc: 'Only receive critical severity alerts' },
            { key: 'dailyDigest' as ToggleKey, label: 'Daily Digest', desc: 'End-of-day summary report' },
          ]).map(({ key, label, desc }) => (
            <View key={key} style={styles.toggleRow}>
              <View style={styles.toggleInfo}>
                <ThemedText style={{ fontSize: 13 }}>{label}</ThemedText>
                <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>{desc}</ThemedText>
              </View>
              <Switch
                value={toggles[key]}
                onValueChange={() => toggle(key)}
                trackColor={{ false: theme.backgroundSelected, true: '#2e7d32' }}
                thumbColor="white"
              />
            </View>
          ))}
        </ThemedView>

        <ThemedView type="backgroundElement" style={styles.section}>
          <View style={styles.sectionHeader}>
            <SymbolView name={{ ios: 'clock.arrow.circlepath', android: 'history', web: 'history' }} size={18} tintColor={theme.text} />
            <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>Data Retention</ThemedText>
          </View>
          <SettingRow label="Audit Logs" value="90 days" />
          <SettingRow label="Auth Attempts" value="30 days" />
          <SettingRow label="Verification Records" value="365 days" />
          <SettingRow label="Activity History" value="180 days" />
        </ThemedView>

        <ThemedView type="backgroundElement" style={styles.section}>
          <View style={styles.sectionHeader}>
            <SymbolView name={{ ios: 'info.circle', android: 'info', web: 'info' }} size={18} tintColor={theme.text} />
            <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>About</ThemedText>
          </View>
          <SettingRow label="Version" value="2.1.0 (build 847)" />
          <SettingRow label="Build Date" value="2026-05-15" />
          <SettingRow label="Environment" value="Production" />
          <Pressable style={styles.linkRow}>
            <ThemedText style={{ color: '#1565c0', fontSize: 13 }}>Licenses &amp; Attribution</ThemedText>
            <SymbolView name={{ ios: 'chevron.right', android: 'chevron_right', web: 'chevron_right' }} size={14} tintColor={theme.textSecondary} />
          </Pressable>
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
  settingRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', paddingVertical: Spacing.one },
  toggleRow: { flexDirection: 'row', alignItems: 'center', paddingVertical: Spacing.one, gap: Spacing.two },
  toggleInfo: { flex: 1, gap: 1 },
  linkRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', paddingVertical: Spacing.one },
});
