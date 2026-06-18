import { StyleSheet, View } from 'react-native';
import { SymbolView } from 'expo-symbols';

import { ThemedText } from './themed-text';
import { ThemedView } from './themed-view';

import { Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

type StatCardProps = {
  title: string;
  value: string | number;
  icon: { ios: string; android: string; web: string };
  trend?: { value: string; positive: boolean };
  color?: string;
};

export function StatCard({ title, value, icon, trend, color }: StatCardProps) {
  const theme = useTheme();

  return (
    <ThemedView type="backgroundElement" style={styles.card}>
      <View style={styles.header}>
        <View style={[styles.iconContainer, { backgroundColor: color ? color + '20' : theme.backgroundSelected }]}>
          <SymbolView name={icon} size={18} tintColor={color || theme.text} />
        </View>
        {trend && (
          <View style={[styles.trendBadge, { backgroundColor: trend.positive ? '#e8f5e9' : '#fce4ec' }]}>
            <ThemedText style={{ color: trend.positive ? '#2e7d32' : '#c62828', fontSize: 11, fontWeight: '700' }}>
              {trend.positive ? '↑' : '↓'} {trend.value}
            </ThemedText>
          </View>
        )}
      </View>
      <ThemedText style={styles.value}>{value}</ThemedText>
      <ThemedText themeColor="textSecondary" style={styles.title}>{title}</ThemedText>
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  card: {
    flex: 1,
    padding: Spacing.three,
    borderRadius: Spacing.three,
    gap: Spacing.two,
    minWidth: 140,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  iconContainer: {
    width: 34,
    height: 34,
    borderRadius: 10,
    justifyContent: 'center',
    alignItems: 'center',
  },
  trendBadge: {
    paddingHorizontal: Spacing.two,
    paddingVertical: 2,
    borderRadius: Spacing.one,
  },
  value: {
    fontSize: 26,
    lineHeight: 32,
    fontWeight: '700',
  },
  title: {
    fontSize: 13,
    lineHeight: 18,
  },
});
