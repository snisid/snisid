import { Pressable, StyleSheet, View } from 'react-native';
import { SymbolView } from 'expo-symbols';

import { ThemedText } from './themed-text';
import { ThemedView } from './themed-view';

import { Spacing } from '@/constants/theme';

type Severity = 'critical' | 'warning' | 'info';

type SecurityAlertProps = {
  icon: { ios: string; android: string; web: string };
  description: string;
  timestamp: string;
  severity: Severity;
  actionLabel?: string;
  onAction?: () => void;
};

const SEVERITY_STYLES: Record<Severity, { color: string; bg: string; border: string }> = {
  critical: { color: '#c62828', bg: '#fce4ec', border: '#ef9a9a' },
  warning: { color: '#e65100', bg: '#fff3e0', border: '#ffcc80' },
  info: { color: '#1565c0', bg: '#e3f2fd', border: '#90caf9' },
};

export function SecurityAlert({ icon, description, timestamp, severity, actionLabel, onAction }: SecurityAlertProps) {
  const s = SEVERITY_STYLES[severity];

  return (
    <ThemedView type="backgroundElement" style={[styles.card, { borderLeftColor: s.border, borderLeftWidth: 3 }]}>
      <View style={[styles.iconWrapper, { backgroundColor: s.bg }]}>
        <SymbolView name={icon} size={20} tintColor={s.color} />
      </View>
      <View style={styles.content}>
        <ThemedText style={{ fontSize: 13 }}>{description}</ThemedText>
        <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>{timestamp}</ThemedText>
      </View>
      {actionLabel && onAction && (
        <Pressable onPress={onAction} style={[styles.actionBtn, { backgroundColor: s.bg }]}>
          <ThemedText style={{ color: s.color, fontSize: 12, fontWeight: '600' }}>{actionLabel}</ThemedText>
        </Pressable>
      )}
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  card: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: Spacing.three,
    borderRadius: Spacing.two,
    gap: Spacing.two,
  },
  iconWrapper: {
    width: 36,
    height: 36,
    borderRadius: 10,
    justifyContent: 'center',
    alignItems: 'center',
  },
  content: {
    flex: 1,
    gap: 2,
  },
  actionBtn: {
    paddingHorizontal: Spacing.two,
    paddingVertical: Spacing.one,
    borderRadius: Spacing.two,
  },
});
