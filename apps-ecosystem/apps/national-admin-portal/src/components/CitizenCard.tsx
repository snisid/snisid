import { Pressable, StyleSheet, View } from 'react-native';

import { ThemedText } from './themed-text';
import { ThemedView } from './themed-view';

import { Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

export type CitizenStatus = 'active' | 'pending' | 'suspended' | 'revoked';

type CitizenCardProps = {
  name: string;
  nnu: string;
  status: CitizenStatus;
  agency: string;
  onPress?: () => void;
};

const STATUS_COLORS: Record<CitizenStatus, { bg: string; text: string }> = {
  active: { bg: '#e8f5e9', text: '#2e7d32' },
  pending: { bg: '#fff3e0', text: '#e65100' },
  suspended: { bg: '#fce4ec', text: '#c62828' },
  revoked: { bg: '#f3e5f5', text: '#6a1b9a' },
};

export function CitizenCard({ name, nnu, status, agency, onPress }: CitizenCardProps) {
  const theme = useTheme();
  const initials = name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
  const statusColor = STATUS_COLORS[status];

  return (
    <Pressable onPress={onPress}>
      <ThemedView type="backgroundElement" style={styles.card}>
        <View style={[styles.avatar, { backgroundColor: theme.backgroundSelected }]}>
          <ThemedText style={{ fontWeight: '700', fontSize: 16 }}>{initials}</ThemedText>
        </View>
        <View style={styles.info}>
          <ThemedText style={{ fontWeight: '600', fontSize: 15 }} numberOfLines={1}>{name}</ThemedText>
          <ThemedText themeColor="textSecondary" style={{ fontSize: 13 }}>NNU: {nnu}</ThemedText>
          <ThemedText themeColor="textSecondary" style={{ fontSize: 13 }} numberOfLines={1}>{agency}</ThemedText>
        </View>
        <View style={[styles.statusBadge, { backgroundColor: statusColor.bg }]}>
          <ThemedText style={{ color: statusColor.text, fontSize: 11, fontWeight: '700', textTransform: 'uppercase' }}>
            {status}
          </ThemedText>
        </View>
      </ThemedView>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: Spacing.three,
    borderRadius: Spacing.three,
    gap: Spacing.three,
  },
  avatar: {
    width: 44,
    height: 44,
    borderRadius: 22,
    justifyContent: 'center',
    alignItems: 'center',
  },
  info: {
    flex: 1,
    gap: 2,
  },
  statusBadge: {
    paddingHorizontal: Spacing.two,
    paddingVertical: Spacing.half,
    borderRadius: Spacing.one,
  },
});
