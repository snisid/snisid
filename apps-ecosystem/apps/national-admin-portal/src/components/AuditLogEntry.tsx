import { StyleSheet, View } from 'react-native';
import { SymbolView } from 'expo-symbols';

import { ThemedText } from './themed-text';
import { ThemedView } from './themed-view';

import { Spacing } from '@/constants/theme';

type Severity = 'critical' | 'warning' | 'info';

type AuditLogEntryProps = {
  icon: { ios: string; android: string; web: string };
  description: string;
  timestamp: string;
  severity: Severity;
};

const SEVERITY_CONFIG: Record<Severity, { color: string; bg: string }> = {
  critical: { color: '#c62828', bg: '#fce4ec' },
  warning: { color: '#e65100', bg: '#fff3e0' },
  info: { color: '#1565c0', bg: '#e3f2fd' },
};

export function AuditLogEntry({ icon, description, timestamp, severity }: AuditLogEntryProps) {
  const config = SEVERITY_CONFIG[severity];

  return (
    <ThemedView type="backgroundElement" style={styles.entry}>
      <View style={[styles.iconWrapper, { backgroundColor: config.bg }]}>
        <SymbolView name={icon} size={16} tintColor={config.color} />
      </View>
      <View style={styles.content}>
        <ThemedText style={{ fontSize: 13 }} numberOfLines={2}>{description}</ThemedText>
        <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>{timestamp}</ThemedText>
      </View>
      <View style={[styles.severityBadge, { backgroundColor: config.bg }]}>
        <ThemedText style={{ color: config.color, fontSize: 10, fontWeight: '700', textTransform: 'uppercase' }}>
          {severity}
        </ThemedText>
      </View>
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  entry: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: Spacing.three,
    borderRadius: Spacing.two,
    gap: Spacing.two,
  },
  iconWrapper: {
    width: 32,
    height: 32,
    borderRadius: 8,
    justifyContent: 'center',
    alignItems: 'center',
  },
  content: {
    flex: 1,
    gap: 2,
  },
  severityBadge: {
    paddingHorizontal: Spacing.one,
    paddingVertical: 2,
    borderRadius: Spacing.one,
  },
});
